package restserveragent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	rad "github.com/ameliepelle/IA04---Prise-de-d-cision-collective-et-vote/agt/restagentdemo"
	procedures "github.com/ameliepelle/IA04---Prise-de-d-cision-collective-et-vote/comsoc"
)

// bureau de vote

type RestServerAgent struct {
	sync.Mutex   // primitive de synchronisation utilisée en programmation informatique pour éviter que des ressources partagées d'un système ne soient utilisées en même temps.
	id           string
	ballotCount  int
	addr         string
	deadline     time.Time
	operator     string
	profile      procedures.Profile
	voters       []string
	alreadyVoted []string
	nbalts       int
}

func NewRestServerAgent(addr string) *RestServerAgent {
	ballot := make(procedures.Profile, 0)
	return &RestServerAgent{addr: addr, profile: ballot}
}

// Test de la méthode
func (rsa *RestServerAgent) checkMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "method %q not allowed", r.Method)
		return false
	}
	return true
}

func (*RestServerAgent) decodeBallotRequest(r *http.Request) (req rad.BallotRequest, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}
func (*RestServerAgent) decodeVoteRequest(r *http.Request) (req rad.VoteRequest, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}
func (*RestServerAgent) decodeResultRequest(r *http.Request) (req rad.ResultRequest, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

func (rsa *RestServerAgent) doNewBallot(w http.ResponseWriter, r *http.Request) {
	// mise à jour du nombre de requêtes
	rsa.Lock()
	defer rsa.Unlock()
	rsa.ballotCount++

	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := rsa.decodeBallotRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	rsa.deadline = req.Deadline
	rsa.operator = req.Rule
	rsa.voters = req.VoterIds
	rsa.id = "vote" + string(rune(rsa.ballotCount))
	rsa.nbalts = req.NbrAlts

	if rsa.deadline != time.Now() && rsa.operator != "" && rsa.voters != nil {
		resp := rad.BallotResponse{BallotId: "vote" + string(rune(rsa.ballotCount))}
		w.WriteHeader(http.StatusCreated)
		serial, _ := json.Marshal(resp)
		w.Write(serial)
		return
	}
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprint(w, err.Error())
}

func (rsa *RestServerAgent) doAddVote(w http.ResponseWriter, r *http.Request) {
	rsa.Lock()
	defer rsa.Unlock()

	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := rsa.decodeVoteRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	if time.Now().After(rsa.deadline) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	voted := false
	voterValid := false
	for _, voter := range rsa.voters {
		if voter == req.AgentId {
			voterValid = true
		}
	}
	if !voterValid {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	for _, voter := range rsa.alreadyVoted {
		if voter == req.AgentId {
			voted = true
		}
	}
	if voted {
		w.WriteHeader(http.StatusAlreadyReported)
		voted = false
		return
	}

	prefs := make([]procedures.Alternative, rsa.nbalts)
	for i, pref := range req.Prefs {
		if i < rsa.nbalts {
			prefs[i] = procedures.Alternative(pref)
		}
	}
	rsa.profile = append(rsa.profile, prefs)

	w.WriteHeader(http.StatusOK)

}

func (rsa *RestServerAgent) doVote(w http.ResponseWriter, r *http.Request) {
	rsa.Lock()
	defer rsa.Unlock()

	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}
	// décodage de la requête
	req, err := rsa.decodeResultRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	if req.BallotId != rsa.id {
		if req.BallotId > rsa.id {
			w.WriteHeader(http.StatusTooEarly)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if time.Now().Before(rsa.deadline) {
		w.WriteHeader(http.StatusTooEarly)
		return
	}
	// traitement de la requête
	var resp rad.ResultResponse
	switch rsa.operator {
	case "borda": // case borda etc
		var result []procedures.Alternative
		log.Println("Profile :", rsa.profile)
		result, err = procedures.BordaSCF(rsa.profile)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		var result2 procedures.Count
		result2, err = procedures.BordaSWF(rsa.profile)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		result3 := procedures.Ranking(result2)
		intResult := make([]int, len(result))
		for i, res := range result3 {
			intResult[i] = int(res)
		}
		resp.Ranking = intResult
		resp.Winner = int(result[0]) //faire un tiebreak ici
	case "majority":
		var result []procedures.Alternative
		log.Println("Profile :", rsa.profile)
		result, err = procedures.MajoritySCF(rsa.profile)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		var result2 procedures.Count
		result2, err = procedures.MajoritySWF(rsa.profile)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		result3 := procedures.Ranking(result2)
		intResult := make([]int, len(result))
		for i, res := range result3 {
			intResult[i] = int(res)
		}
		resp.Ranking = intResult
		resp.Winner = int(result[0])
	case "approval":
		var result []procedures.Alternative
		log.Println("Profile :", rsa.profile)
		result, err = procedures.ApprovalSCF(rsa.profile, rand.Perm(len(rsa.profile[0])))
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		/*intResult := make([]int, len(result))
		for i, res := range result {
			intResult[i] = int(res)
		}
		resp.Ranking = intResult
		var result2 procedures.Count
		result2, err = procedures.ApprovalSWF(rsa.profile, rand.Perm(len(rsa.profile[0])))
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		resp.Winner = result2[0]*/
		resp.Winner = int(result[0])
	case "condorcet":
		var result []procedures.Alternative
		log.Println("Profile :", rsa.profile)
		result, err = procedures.CondorcetWinner(rsa.profile)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		if result == nil {
			log.Println(w, "Pas de gagnant de Condorcet")
			resp.Winner = -1
			return
		}
		resp.Winner = int(result[0])
	case "kemeny":
		var result []procedures.Alternative
		log.Println("Profile :", rsa.profile)
		result, err = procedures.Kemeny(rsa.profile)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		intResult := make([]int, len(result))
		for i, res := range result {
			intResult[i] = int(res)
		}
		resp.Ranking = intResult
		resp.Winner = int(result[0])
	default:
		w.WriteHeader(http.StatusNotImplemented)
		msg := fmt.Sprintf("Unkonwn command '%s'", rsa.operator)
		w.Write([]byte(msg))
		return
	}

	w.WriteHeader(http.StatusOK)
	serial, _ := json.Marshal(resp)
	w.Write(serial)
}

func (rsa *RestServerAgent) doBallotcount(w http.ResponseWriter, r *http.Request) {
	if !rsa.checkMethod("GET", w, r) {
		return
	}

	w.WriteHeader(http.StatusOK)
	rsa.Lock()
	defer rsa.Unlock()
	serial, _ := json.Marshal(rsa.ballotCount)
	w.Write(serial)
}

func (rsa *RestServerAgent) Start() {
	// création du multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/new_ballot", rsa.doNewBallot)
	mux.HandleFunc("/vote", rsa.doAddVote)
	mux.HandleFunc("/result", rsa.doVote)
	mux.HandleFunc("/reqcount", rsa.doBallotcount)

	// création du serveur http
	s := &http.Server{
		Addr:           rsa.addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	// lancement du serveur
	log.Println("Listening on", rsa.addr)
	go log.Fatal(s.ListenAndServe())
}
