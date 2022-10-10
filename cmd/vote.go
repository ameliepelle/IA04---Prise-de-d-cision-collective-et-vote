package main

import (
	procedures "TD3/comsoc"
	"fmt"
)

func main() {
	alt := []procedures.Alternative{4, 3, 5, 7, 1}
	m := make(procedures.Count, 4)
	m[1] = 6
	m[2] = 1
	m[3] = 5
	m[4] = 6

	fmt.Println(procedures.Rank(5, alt))
	fmt.Println(procedures.IsPref(3, 4, alt))
	fmt.Println(procedures.IsPref(3, 1, alt))
	fmt.Println(procedures.MaxCount(m))

	alt2 := []procedures.Alternative{2, 1, 7, 8}
	prof := procedures.Profile{
		{2, 1, 7, 8},
		{2, 1, 7, 8},
		{8, 7, 2, 1},
		{8, 7, 2, 1},
	}
	thre := []int{3, 3, 1, 1}

	fmt.Println(procedures.ApprovalSCF(prof, thre))
	alt3, err := procedures.ApprovalSCF(prof, thre)
	fmt.Println(err)
	fmt.Println(procedures.TieBreak(alt3))
	fmt.Println("Ordre", alt2)
	f := procedures.TieBreakFactory(alt2)
	fmt.Println(f(alt3))
}
