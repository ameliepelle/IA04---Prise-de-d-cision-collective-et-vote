package main

import (
	restclientagent "TD3/agt/restagentdemo/restclientagent"
	restserveragent "TD3/agt/restagentdemo/restserveragent"
	procedures "TD3/comsoc"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func main() {
	const n = 100
	const url1 = ":8080"
	const url2 = "http://localhost:8080"
	ops := [...]string{"borda", "majority", "approval", "condorcet", "kemeny"} // Il faut gérer le cas nul avec Condorcet

	clAgts := make([]restclientagent.RestClientAgent, 0, n)
	servAgt := restserveragent.NewRestServerAgent(url1)

	log.Println("démarrage du serveur...")
	go servAgt.Start()

	log.Println("démarrage des clients...")
	for i := 0; i < n; i++ {
		id := fmt.Sprintf("id%02d", i)
		prefsInt := rand.Perm(16)
		prefs := make([]procedures.Alternative, 16)
		for j, pref := range prefsInt {
			prefs[j] = procedures.Alternative(pref)
		}
		options := rand.Perm(16)
		agt := restclientagent.NewRestClientAgent(id, "vote0", url2, prefs, options) // mettre prefs a la place de op1 op2
		clAgts = append(clAgts, *agt)
	}
	op := ops[0]
	deadline := time.Now().Add(30 * time.Second)
	ballot := restclientagent.NewBallotAgent("vote0", op, deadline, []string{"id1", "id2", "id3", "id5"}, 3, url2)
	ballot.Start()

	for _, agt := range clAgts {
		agt.SetVoteId(ballot.GetId())
		// attention, obligation de passer par cette lambda pour faire capturer la valeur de l'itération par la goroutine
		func(agt restclientagent.RestClientAgent) {
			go agt.Start()
		}(agt)
	}

	time.Sleep((time.Until(ballot.Deadline)))
	fmt.Println("slept")
	ballot.Result()

	fmt.Scanln()
}
