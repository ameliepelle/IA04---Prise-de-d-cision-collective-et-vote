package procedures

import (
	"fmt"

	"gitlab.utc.fr/lagruesy/ia04/utils"
)

func TauKendall(arr1 []Alternative, arr2 []Alternative) (tau int, err error) {
	err = CheckProfile([][]Alternative{arr1, arr2})
	if err != nil {
		return
	}
	//Utiliser la fonction IsPref
	var pref1 [][]Alternative
	var pref2 [][]Alternative

	for i := range arr1 {
		for j := i + 1; j < len(arr1); j++ {

			//On construit les premières préférences
			var inter1 []Alternative
			inter1 = append(inter1, arr1[i])
			inter1 = append(inter1, arr1[j])
			pref1 = append(pref1, inter1)
		}
	}

	//On construit les 2èmes préférences
	for i := range arr2 {
		for j := i + 1; j < len(arr2); j++ {

			//On construit les premières préférences
			var inter1 []Alternative
			inter1 = append(inter1, arr2[i])
			inter1 = append(inter1, arr2[j])
			pref2 = append(pref2, inter1)
		}
	}

	//fmt.Println(pref1)
	//fmt.Println(pref2)
	for _, e := range pref1 {
		for _, f := range pref2 {
			if e[0] == f[1] && e[1] == f[0] { // si les paires sont inversées
				tau = tau + 1
			}
		}

	}
	return tau, err
	//return (count - (len(pref1)-count))/len(pref1
} //ok

func DistRP(arr1 []Alternative, p Profile) (dist int, err error) {
	err = CheckProfile(append(p, arr1))
	if err != nil {
		return
	}

	for _, e := range p {
		temp, err2 := TauKendall(arr1, e)
		dist += temp
		if err2 != nil {
			return
		}

	}

	return dist, err
}

func Kemeny(p Profile) (bestAlt []Alternative, err error) {
	var distMin int                           //on garde simplement la permutation avec la distance minimale
	perm := utils.FirstPermutation(len(p[0])) //On liste toutes les permutations et on calcule la distance dessus
	fmt.Println(perm)
	//On crée l'équivalent en alternative de la permutation
	//On initialise les variables
	bestAlt = make([]Alternative, len(p[0]))
	for j, i := range perm {
		bestAlt[j] = p[0][i]
	}
	fmt.Println("bestAlt ", bestAlt)
	distMin, err = DistRP(bestAlt, p)
	fmt.Println("distMin", distMin)
	if err != nil {
		return
	}

	//On fait le premier test
	perm, ok := utils.NextPermutation(perm)

	for ok {
		altTemp := make([]Alternative, len(p[0]))
		fmt.Println("perm ", perm)
		for j, i := range perm {
			fmt.Println("p[0][i] ", p[0][i])
			altTemp[j] = p[0][i]
		}
		fmt.Println("altTemp", altTemp)
		distemp, err := DistRP(altTemp, p)
		if err != nil {
			return bestAlt, err
		}
		if distemp < distMin {
			distMin = distemp
			copy(bestAlt, altTemp)
			fmt.Println("bestAlt ", bestAlt)

		}
		perm, ok = utils.NextPermutation(perm)
	}
	fmt.Println("min ", distMin)
	return bestAlt, err
}
