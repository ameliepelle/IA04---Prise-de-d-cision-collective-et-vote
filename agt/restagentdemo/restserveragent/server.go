package restserveragent

import (
	rad "TD3/agt/restagentdemo"
	procedures "TD3/comsoc"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// bureau de vote

type RestServerAgent struct {
	sync.Mutex // primitive de synchronisation utilisée en programmation informatique pour éviter que des ressources partagées d'un système ne soient utilisées en même temps.
	id         string
	reqCount   int
	addr       string
	profile    procedures.Profile
}

func NewRestServerAgent(addr string) *RestServerAgent {
	return &RestServerAgent{id: addr, addr: addr}
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

func (*RestServerAgent) decodeRequest(r *http.Request) (req rad.Request, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

func (rsa *RestServerAgent) doCalc(w http.ResponseWriter, r *http.Request) {
	// mise à jour du nombre de requêtes
	rsa.Lock()
	defer rsa.Unlock()
	rsa.reqCount++

	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := rsa.decodeRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	rsa.profile = append(rsa.profile, req.Prefs)

	if rsa.reqCount < 100 { // à changer pour mettre une deadline time
		fmt.Println("not enough votes yet")
		return
	}

	// traitement de la requête
	var resp rad.Response

	switch req.Operator {
	case "borda": // case borda etc
		var result []procedures.Alternative
		result, err = procedures.BordaSCF(rsa.profile)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		resp.Result = int(result[0]) // comment faire pour récup un profil alors que chaque votant donne une liste d'alternatives?
	case "majority":
		var result []procedures.Alternative
		result, err = procedures.MajoritySCF(rsa.profile)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		resp.Result = int(result[0])
	case "approval":
		var result []procedures.Alternative
		result, err = procedures.ApprovalSCF(rsa.profile, rand.Perm(100))
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		resp.Result = int(result[0])
	case "condorcet":
		var result []procedures.Alternative
		result, err = procedures.CondorcetWinner(rsa.profile)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		if result == nil {
			fmt.Fprint(w, "Pas de gagnant de Condorcet")
			return
		}
		resp.Result = int(result[0])
	default:
		w.WriteHeader(http.StatusNotImplemented)
		msg := fmt.Sprintf("Unkonwn command '%s'", req.Operator)
		w.Write([]byte(msg))
		return
	}

	w.WriteHeader(http.StatusOK)
	serial, _ := json.Marshal(resp)
	w.Write(serial)
}

func (rsa *RestServerAgent) doReqcount(w http.ResponseWriter, r *http.Request) {
	if !rsa.checkMethod("GET", w, r) {
		return
	}

	w.WriteHeader(http.StatusOK)
	rsa.Lock()
	defer rsa.Unlock()
	serial, _ := json.Marshal(rsa.reqCount)
	w.Write(serial)
}

func (rsa *RestServerAgent) Start() {
	// création du multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/calculator", rsa.doCalc)
	mux.HandleFunc("/reqcount", rsa.doReqcount)

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
