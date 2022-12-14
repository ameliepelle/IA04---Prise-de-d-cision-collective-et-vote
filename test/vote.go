package main

import (
	"fmt"
	"sort"

	procedures "github.com/ameliepelle/IA04---Prise-de-d-cision-collective-et-vote/comsoc"

	"gitlab.utc.fr/lagruesy/ia04/utils"
)

func main() {
	//Test ranking
	m2 := make(procedures.Count, 4)
	keys := make([]procedures.Alternative, 0, len(m2))
	m2[1] = 6
	m2[2] = 1
	m2[3] = 5
	m2[4] = 6
	for key := range m2 {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return m2[keys[i]] < m2[keys[j]]
	})

	fmt.Println("keys !", keys)

	//Test Kemeny
	prof2 := procedures.Profile{
		{1, 2, 3},
		{3, 2, 1},
		{1, 3, 2},
	}
	fmt.Println("vainqueur Kemeny ")
	fmt.Println(procedures.Kemeny(prof2))

	arr1 := []procedures.Alternative{1, 2, 3}
	arr2 := []procedures.Alternative{1, 2, 3}
	fmt.Println(procedures.TauKendall(arr1, arr2))

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
	thre := []int{4, 3, 1, 1}

	fmt.Println(procedures.ApprovalSCF(prof, thre))
	alt3, err := procedures.ApprovalSCF(prof, thre)
	fmt.Println(err)
	fmt.Println(procedures.TieBreak(alt3))
	fmt.Println("Ordre", alt2)
	f := procedures.TieBreakFactory(alt2)
	fmt.Println(f(alt3))

	fmt.Println("permutations")
	const n = 4

	// création et affichage de la première permutation
	perm := utils.FirstPermutation(n)
	fmt.Println(perm)

	// compteur des permutations
	count := 1

	// itération et affichage de la permutation suivante
	perm, ok := utils.NextPermutation(perm)
	for ok {
		count++
		fmt.Println(perm)
		perm, ok = utils.NextPermutation(perm)
	}

	// on affiche la valeur de la factorielle...
	fmt.Printf("\n%d!=%d\n", n, count)

}
