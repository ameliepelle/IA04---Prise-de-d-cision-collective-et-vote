package main

import (
	"fmt"

	ras "github.com/ameliepelle/IA04---Prise-de-d-cision-collective-et-vote/agt/restagentdemo/restserveragent"
)

func main() {
	server := ras.NewRestServerAgent(":8080")
	server.Start()
	fmt.Scanln()
}
