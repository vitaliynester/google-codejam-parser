package models

type Response struct {
	FileName      string `json:"file_name"`
	FileUrl       string `json:"file_url"`
	UserID        string `json:"user_id"`
	UserName      string `json:"user_name"`
	ChallengeID   string `json:"challenge_id"`
	ChallengeName string `json:"challenge_name"`
	AdventureID   string `json:"adventure_id"`
	AdventureName string `json:"adventure_name"`
}

type AttemptRequest struct {
	CompetitorID   string `json:"competitor_id"`
	NonFinalResult bool   `json:"include_non_final_results"`
}

type AttemptsResult struct {
	Attempts []AttemptResult `json:"attempts"`
}

type AttemptResult struct {
	SourceFile struct {
		Filename string `json:"filename"`
		Url      string `json:"url"`
	} `json:"source_file"`
}

type Challenge struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Adventure struct {
	ID         string      `json:"id"`
	Title      string      `json:"title"`
	Challenges []Challenge `json:"challenges"`
}

type AdventureResponse struct {
	Adventures []Adventure `json:"adventures"`
}

type ScoreboardPagination struct {
	MinRank int64 `json:"min_rank"`
	Count   int64 `json:"num_consecutive_users"`
}

type UserScore struct {
	Competitor struct {
		Name string `json:"displayname"`
		ID   string `json:"id"`
	} `json:"competitor"`
}

type Scoreboard struct {
	Size       int64       `json:"full_scoreboard_size"`
	UserScores []UserScore `json:"user_scores"`
}

type ChallengeResponse struct {
	Challenge  Challenge   `json:"challenge"`
	UserScores []UserScore `json:"user_scores"`
}

type ScoreboardResponse struct {
	AdventureName string              `json:"adventure_name"`
	AdventureID   string              `json:"adventure_id"`
	Challenges    []ChallengeResponse `json:"challenges"`
}
