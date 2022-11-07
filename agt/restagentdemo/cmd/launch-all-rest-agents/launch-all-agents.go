package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	restclientagent "github.com/ameliepelle/IA04---Prise-de-d-cision-collective-et-vote/agt/restagentdemo/restclientagent"
	restserveragent "github.com/ameliepelle/IA04---Prise-de-d-cision-collective-et-vote/agt/restagentdemo/restserveragent"
	procedures "github.com/ameliepelle/IA04---Prise-de-d-cision-collective-et-vote/comsoc"
)

func main() {
	const n = 100
	const url1 = ":8080"
	const url2 = "http://localhost:8080"
	ops := [...]string{"borda", "majority", "approval", "condorcet", "kemeny"} // Il faut gérer le cas nul avec Condorcet

	clAgts := make([]restclientagent.RestClientAgent, 0, n)
	servAgt := restserveragent.NewRestServerAgent(url1)

	var voters []string
	log.Println("démarrage du serveur...")
	go servAgt.Start()

	log.Println("démarrage des clients...")
	for i := 0; i < n; i++ {
		id := fmt.Sprintf("id%02d", i)
		voters = append(voters, id)
		prefsInt := rand.Perm(5)
		prefs := make([]procedures.Alternative, 5)
		for j, pref := range prefsInt {
			prefs[j] = procedures.Alternative(pref)
		}
		agt := restclientagent.NewRestClientAgent(id, "vote0", url2, prefs, []int{rand.Intn(len(prefs))}) // mettre prefs a la place de op1 op2
		clAgts = append(clAgts, *agt)
	}
	op := ops[2]
	deadline := time.Now().Add(3 * time.Second)
	ballot := restclientagent.NewBallotAgent("vote0", op, deadline, voters, 5, url2)
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
