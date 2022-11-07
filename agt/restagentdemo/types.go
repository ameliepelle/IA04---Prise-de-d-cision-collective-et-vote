package restagentdemo

import (
	"time"
)

type BallotRequest struct {
	Rule     string    `json:"rule"`
	Deadline time.Time `json:"deadline"`
	VoterIds []string  `json:"voter-ids"`
	NbrAlts  int       `json:"#alts"`
}

type BallotResponse struct {
	BallotId string `json:"ballot-id"`
}

type VoteRequest struct {
	AgentId string `json:"agent-id"`
	VoteId  string `json:"vote-id"`
	Prefs   []int  `json:"prefs"`
	Options []int  `json:"options"`
}

type ResultRequest struct {
	BallotId string `json:"ballot-id"`
}

type ResultResponse struct {
	Winner  int   `json:"winner"`
	Ranking []int `json:"ranking"`
}
