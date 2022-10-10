package main

import (
	ras "TD3/agt/restagentdemo/restserveragent"
	"fmt"
)

func main() {
	server := ras.NewRestServerAgent(":8080")
	server.Start()
	fmt.Scanln()
}
