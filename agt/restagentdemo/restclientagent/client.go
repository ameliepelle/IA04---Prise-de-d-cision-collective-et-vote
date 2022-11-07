package restclientagent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	procedures "github.com/ameliepelle/IA04---Prise-de-d-cision-collective-et-vote/comsoc"

	rad "github.com/ameliepelle/IA04---Prise-de-d-cision-collective-et-vote/agt/restagentdemo"
)

// votant
// operator a mettre dans le serveur au lieu du client
type BallotAgent struct {
	id       string
	Rule     string
	Deadline time.Time
	VoterIds []string
	NbrAlts  int
	url      string
}

func NewBallotAgent(id string, rule string, deadline time.Time, voterIds []string, nbr int, url string) *BallotAgent {
	return &BallotAgent{id, rule, deadline, voterIds, nbr, url}
}

type RestClientAgent struct {
	id      string
	voteId  string
	url     string
	prefs   []procedures.Alternative
	options []int
}

func NewRestClientAgent(id string, vote string, url string, prefs []procedures.Alternative, opt []int) *RestClientAgent {
	return &RestClientAgent{id, vote, url, prefs, opt}
}

func (rca *BallotAgent) treatBallotResponse(r *http.Response) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	var resp rad.BallotResponse
	json.Unmarshal(buf.Bytes(), &resp)

	return resp.BallotId
}
func (b *BallotAgent) treatResultResponse(r *http.Response) int {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	var resp rad.ResultResponse
	json.Unmarshal(buf.Bytes(), &resp)
	return resp.Winner
}

func (b *BallotAgent) doBallotRequest() (res string, err error) {
	req := rad.BallotRequest{
		Rule:     b.Rule,
		Deadline: b.Deadline,
		VoterIds: b.VoterIds,
		NbrAlts:  b.NbrAlts,
	}

	// sérialisation de la requête
	url := b.url + "/new_ballot"
	data, _ := json.Marshal(req)

	// envoi de la requête
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	// traitement de la réponse
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusCreated {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		return
	}
	res = b.treatBallotResponse(resp)

	return
}

func (rca *RestClientAgent) doVoteRequest() (err error) {
	alts := make([]int, len(rca.prefs))
	for i, pref := range rca.prefs {
		alts[i] = int(pref)
	}
	req := rad.VoteRequest{
		AgentId: rca.id,
		VoteId:  rca.voteId,
		Prefs:   alts,
		Options: rca.options,
	}

	// sérialisation de la requête
	url := rca.url + "/vote"
	data, _ := json.Marshal(req)

	// envoi de la requête
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	// traitement de la réponse
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusCreated {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		return
	}

	return
}

func (b *BallotAgent) doResultRequest() (res int, err error) {
	req := rad.ResultRequest{
		BallotId: b.id,
	}

	// sérialisation de la requête
	url := b.url + "/result"
	data, _ := json.Marshal(req)
	// envoi de la requête
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	// traitement de la réponse
	if err != nil {
		log.Fatal("error:", err.Error())
		return
	}
	log.Println("status code", resp.StatusCode)
	log.Println("status created", http.StatusOK)
	err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
	if resp.StatusCode != http.StatusOK {
		return
	}
	res = b.treatResultResponse(resp)
	return
}

func (rca *RestClientAgent) Start() {
	log.Printf("démarrage de %s", rca.id)
	rca.doVoteRequest()

	log.Printf("[%s] %v \n", rca.id, rca.prefs)
}

func (rca *RestClientAgent) SetVoteId(voteId string) {
	rca.voteId = voteId
}

func (b *BallotAgent) GetId() string {
	return b.id
}

func (b *BallotAgent) Start() {
	log.Printf("démarrage du ballot")

	res, err := b.doBallotRequest()

	if err != nil {
		log.Fatal("error:", err.Error())
	} else {
		b.id = res
	}
}

func (b *BallotAgent) Result() {
	res, err := b.doResultRequest()
	if err.Error() != "[200] 200 OK" {
		log.Fatal("error:", err.Error())
		return
	}

	log.Println("gagnant :", res)
}
