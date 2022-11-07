package main

import (
	"fmt"
	"time"

	restclientagent "github.com/ameliepelle/IA04---Prise-de-d-cision-collective-et-vote/agt/restagentdemo/restclientagent"
	procedures "github.com/ameliepelle/IA04---Prise-de-d-cision-collective-et-vote/comsoc"
)

func main() {
	ag := restclientagent.NewRestClientAgent("id1", "vote0", "http://localhost:8000", []procedures.Alternative{3, 5, 4, 8, 1}, []int{3})
	ballot := restclientagent.NewBallotAgent("vote0", "borda", time.Date(2022, 11, 7, 11, 30, 0, 0, time.Now().Location()), []string{"id1"}, 3, "http://localhost:3000")
	ballot.Start()
	ag.Start()
	fmt.Scanln()
}
