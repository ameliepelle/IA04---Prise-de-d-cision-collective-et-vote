package restclientagent

import (
	rad "TD3/agt/restagentdemo"
	procedures "TD3/comsoc"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// votant
// operator a mettre dans le serveur au lieu du client

type RestClientAgent struct {
	id       string
	url      string
	operator string
	prefs    []procedures.Alternative
}

func NewRestClientAgent(id string, url string, op string, prefs []procedures.Alternative) *RestClientAgent {
	return &RestClientAgent{id, url, op, prefs}
}

func (rca *RestClientAgent) treatResponse(r *http.Response) int {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	var resp rad.Response
	json.Unmarshal(buf.Bytes(), &resp)

	return resp.Result
}

func (rca *RestClientAgent) doRequest() (res int, err error) {
	req := rad.Request{
		Operator: rca.operator,
		Prefs:    rca.prefs,
	}

	// sérialisation de la requête
	url := rca.url + "/calculator"
	data, _ := json.Marshal(req)

	// envoi de la requête
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	// traitement de la réponse
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		return
	}
	res = rca.treatResponse(resp)

	return
}

func (rca *RestClientAgent) Start() {
	log.Printf("démarrage de %s", rca.id)
	res, err := rca.doRequest()

	if err != nil {
		log.Fatal(rca.id, "error:", err.Error())
	} else {
		log.Printf("[%s] %v %s = %d\n", rca.id, rca.prefs, rca.operator, res)
	}
}
