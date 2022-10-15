package models

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

type ScoreboardResponse struct {
	AdventureName string      `json:"adventure_name"`
	AdventureID   string      `json:"adventure_id"`
	ChallengeName string      `json:"challenge_name"`
	ChallengeID   string      `json:"challenge_id"`
	UserScores    []UserScore `json:"user_scores"`
}
