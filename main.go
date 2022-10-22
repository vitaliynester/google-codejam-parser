package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"example.com/m/models"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

func encodeToBase64(v interface{}) (string, error) {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	err := json.NewEncoder(encoder).Encode(v)
	if err != nil {
		return "", err
	}
	encoder.Close()
	return buf.String(), nil
}

func decodeFromBase64(data []byte) []byte {
	b64 := make([]byte, base64.RawURLEncoding.DecodedLen(len(data)))
	_, _ = base64.RawURLEncoding.Decode(b64, data)
	return b64
}

func makeResponse(targetUrl string) []byte {
	options := cookiejar.Options{}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{Jar: jar}

	baseUrl, err := url.Parse(targetUrl)
	if err != nil {
		log.Fatal(err)
	}
	req, _ := http.NewRequest("GET", baseUrl.String(), nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	data, err := io.ReadAll(resp.Body)
	resp.Body.Close()

	return data
}

func main() {
	data := makeResponse("https://codejam.googleapis.com/poll?p=e30")
	b64 := decodeFromBase64(data)
	var responseModel models.AdventureResponse
	err := json.Unmarshal(b64, &responseModel)
	if err != nil {
		log.Fatal(err)
	}

	startPagination := models.ScoreboardPagination{
		MinRank: 1,
		Count:   50,
	}
	startPaginationStr, err := encodeToBase64(startPagination)
	if err != nil {
		log.Fatal(err)
	}
	var result []models.ScoreboardResponse
	for _, adventure := range responseModel.Adventures {
		var scoreboardResponse models.ScoreboardResponse
		scoreboardResponse.AdventureID = adventure.ID
		scoreboardResponse.AdventureName = adventure.Title

		for _, challenge := range adventure.Challenges {
			var challengeResponse models.ChallengeResponse
			challengeResponse.Challenge = challenge
			newUrl := fmt.Sprintf("https://codejam.googleapis.com/scoreboard/%v/poll?p=%v", challenge.ID, startPaginationStr)
			resp := makeResponse(newUrl)
			data := decodeFromBase64(resp)

			var scoreboard models.Scoreboard
			err = json.Unmarshal(data, &scoreboard)
			if err != nil {
				log.Fatal(err)
			}
			challengeResponse.UserScores = append(challengeResponse.UserScores, scoreboard.UserScores...)

			var sum int64
			sum = 51
			for sum < scoreboard.Size {
				pagination := models.ScoreboardPagination{
					MinRank: sum,
					Count:   50,
				}
				paginationStr, err := encodeToBase64(pagination)
				if err != nil {
					log.Fatal(err)
				}

				newUrl = fmt.Sprintf("https://codejam.googleapis.com/scoreboard/%v/poll?p=%v", challenge.ID, paginationStr)
				resp = makeResponse(newUrl)
				data = decodeFromBase64(resp)

				var includedScoreboard models.Scoreboard
				err = json.Unmarshal(data, &includedScoreboard)
				if err != nil {
					log.Fatal(err)
				}
				challengeResponse.UserScores = append(challengeResponse.UserScores, includedScoreboard.UserScores...)

				sum += 50
			}
			scoreboardResponse.Challenges = append(scoreboardResponse.Challenges, challengeResponse)

			fmt.Printf("Количество участников в %v, %v: %v\n", adventure.Title, challenge.Title, scoreboard.Size)
		}
		result = append(result, scoreboardResponse)
	}
	resultFile, _ := json.MarshalIndent(result, "", "  ")
	_ = ioutil.WriteFile("result.json", resultFile, 0644)

	var results []models.ScoreboardResponse
	info, _ := ioutil.ReadFile("result.json")
	err = json.Unmarshal(info, &results)
	if err != nil {
		log.Fatal(err)
	}
	var resultToFile []models.Response
	totalFiles := 0
	for _, adventure := range result {
		for _, challenge := range adventure.Challenges {
			for _, userScore := range challenge.UserScores {
				attemptRequest := models.AttemptRequest{
					CompetitorID:   userScore.Competitor.ID,
					NonFinalResult: false,
				}
				attemptRequestStr, _ := encodeToBase64(attemptRequest)
				newUrl := fmt.Sprintf("https://codejam.googleapis.com/attempts/%v/poll?p=%v", challenge.Challenge.ID, attemptRequestStr)
				resp := makeResponse(newUrl)
				data := decodeFromBase64(resp)

				var attemptsResponse models.AttemptsResult
				err = json.Unmarshal(data, &attemptsResponse)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("Количество файлов у %v в соревновании %v: %v шт.\n", userScore.Competitor.Name, challenge.Challenge.Title, len(attemptsResponse.Attempts))
				totalFiles += len(attemptsResponse.Attempts)
				for _, attempt := range attemptsResponse.Attempts {
					toFile := models.Response{
						FileName:      attempt.SourceFile.Filename,
						FileUrl:       attempt.SourceFile.Url,
						UserID:        userScore.Competitor.ID,
						UserName:      userScore.Competitor.Name,
						ChallengeID:   challenge.Challenge.ID,
						ChallengeName: challenge.Challenge.Title,
						AdventureID:   adventure.AdventureID,
						AdventureName: adventure.AdventureName,
					}
					resultToFile = append(resultToFile, toFile)
				}
			}
		}
	}
	fmt.Printf("Суммарное количество файлов для загрузки: %v\n", totalFiles)
	resultFile, _ = json.MarshalIndent(resultToFile, "", "  ")
	_ = ioutil.WriteFile("final_result.json", resultFile, 0644)
}
