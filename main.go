package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
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
	var responseModel AdventureResponse
	err := json.Unmarshal(b64, &responseModel)
	if err != nil {
		log.Fatal(err)
	}
	file, _ := json.MarshalIndent(responseModel, "", "  ")
	_ = ioutil.WriteFile("adventure_result.json", file, 0644)

	startPagination := ScoreboardPagination{
		MinRank: 1,
		Count:   50,
	}
	startPaginationStr, err := encodeToBase64(startPagination)
	if err != nil {
		log.Fatal(err)
	}
	var result []ScoreboardResponse
	for _, adventure := range responseModel.Adventures {
		var scoreboardResponse ScoreboardResponse
		scoreboardResponse.AdventureID = adventure.ID
		scoreboardResponse.AdventureName = adventure.Title

		for _, challenge := range adventure.Challenges {
			newUrl := fmt.Sprintf("https://codejam.googleapis.com/scoreboard/%v/poll?p=%v", challenge.ID, startPaginationStr)
			resp := makeResponse(newUrl)
			data := decodeFromBase64(resp)

			var scoreboard Scoreboard
			err = json.Unmarshal(data, &scoreboard)
			if err != nil {
				log.Fatal(err)
			}
			scoreboardResponse.ChallengeID = challenge.ID
			scoreboardResponse.ChallengeName = challenge.Title
			scoreboardResponse.UserScores = append(scoreboardResponse.UserScores, scoreboard.UserScores...)

			var sum int64
			sum = 51
			for sum < scoreboard.Size {
				pagination := ScoreboardPagination{
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

				var includedScoreboard Scoreboard
				err = json.Unmarshal(data, &includedScoreboard)
				if err != nil {
					log.Fatal(err)
				}
				scoreboardResponse.UserScores = append(scoreboardResponse.UserScores, includedScoreboard.UserScores...)

				sum += 50
			}

			fmt.Printf("Количество участников в %v, %v: %v\n", adventure.Title, challenge.Title, scoreboard.Size)
		}
		result = append(result, scoreboardResponse)
	}
	resultFile, _ := json.MarshalIndent(result, "", "  ")
	_ = ioutil.WriteFile("result.json", resultFile, 0644)
}
