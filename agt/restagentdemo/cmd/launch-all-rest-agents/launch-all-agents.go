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
	const n = 100 // nombre de votants
	const m = 5   // nombre d'alternatives
	const url1 = ":8080"
	const url2 = "http://localhost:8080"
	ops := [...]string{"borda", "majority", "approval", "condorcet", "kemeny"} // méthodes de vote implémentées

	clAgts := make([]restclientagent.RestClientAgent, 0, n)
	servAgt := restserveragent.NewRestServerAgent(url1)

	var voters []string
	log.Println("démarrage du serveur...")
	go servAgt.Start()

	log.Println("démarrage des clients...")
	for i := 0; i < n; i++ {
		id := fmt.Sprintf("id%02d", i)
		voters = append(voters, id)
		prefsInt := rand.Perm(m) // preférences aléatoires
		prefs := make([]procedures.Alternative, m)
		for j, pref := range prefsInt { // conversion vers le type Alternative
			prefs[j] = procedures.Alternative(pref)
		}
		agt := restclientagent.NewRestClientAgent(id, "vote0", url2, prefs, []int{rand.Intn(len(prefs))})
		clAgts = append(clAgts, *agt)
	}
	op := ops[2]                                // changer l'indice pour sélectionner une méthode de vote différente
	deadline := time.Now().Add(3 * time.Second) // deadline après 3 secondes (changer l'argument de Add() pour changer la deadline)
	ballot := restclientagent.NewBallotAgent("vote0", op, deadline, voters, m, url2)
	ballot.Start()

	for _, agt := range clAgts {
		agt.SetVoteId(ballot.GetId())
		// attention, obligation de passer par cette lambda pour faire capturer la valeur de l'itération par la goroutine
		func(agt restclientagent.RestClientAgent) {
			go agt.Start()
		}(agt)
	}

	time.Sleep((time.Until(ballot.Deadline))) // attente de la deadline pour lancer le calcul du résultat de vote
	fmt.Println("calcul du résultat ...")
	ballot.Result()
	fmt.Scanln()
}
