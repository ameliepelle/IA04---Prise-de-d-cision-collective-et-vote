package main

import (
	restclientagent "TD3/agt/restagentdemo/restclientagent"
	procedures "TD3/comsoc"
	"fmt"
)

func main() {
	ag := restclientagent.NewRestClientAgent("id1", "http://localhost:8000", "borda", []procedures.Alternative{3, 5, 4, 8, 1}) // mettre borda par ex a la place de "+"
	ag.Start()
	fmt.Scanln()
}
