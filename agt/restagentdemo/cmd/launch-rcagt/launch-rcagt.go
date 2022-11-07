package main

import (
	"fmt"
	"time"

	restclientagent "github.com/ameliepelle/IA04---Prise-de-d-cision-collective-et-vote/agt/restagentdemo/restclientagent"
	procedures "github.com/ameliepelle/IA04---Prise-de-d-cision-collective-et-vote/comsoc"
)

func main() {
	ag := restclientagent.NewRestClientAgent("id1", "vote0", "http://localhost:8000", []procedures.Alternative{3, 5, 4, 8, 1}, []int{3})
	ballot := restclientagent.NewBallotAgent("vote0", "borda", time.Now().Add(3*time.Second), []string{"id1"}, 5, "http://localhost:3000")
	ballot.Start()
	ag.Start()
	fmt.Scanln()
}
