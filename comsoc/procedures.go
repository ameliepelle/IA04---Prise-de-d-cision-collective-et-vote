package procedures

import (
	"errors"
	"sort"
)

type Alternative int
type Profile [][]Alternative // [][]int
type Count map[Alternative]int

func Rank(alt Alternative, prefs []Alternative) int { // renvoie l'indice où se trouve alt dans prefs
	// retourne la valeur de l'indice où se trouve alt ou bien -1 si alt n'est pas dans prefs
	for i, v := range prefs {
		if v == alt {
			return i
		}
	}
	return -1
}

func IsPref(alt1, alt2 Alternative, prefs []Alternative) bool { // renvoie vrai ssi alt1 est préférée à alt2
	// en utilisant Rank, on récupère l'indice auquel chaque alternative se trouve
	indice1 := Rank(alt1, prefs)
	indice2 := Rank(alt2, prefs)

	// si alt1 n'est pas dans les préférences, on revoie false (elle ne peut pas être préférée à alt2)
	if indice1 == -1 {
		return false
	}

	return indice1 < indice2
}

// func IsPref(alt1, alt2 Alternative, prefs []Alternative) bool { // renvoie vrai ssi alt1 est préférée à alt2
// 	var indice1, indice2 int
// 	for i, v := range prefs {
// 		if v == alt1 {
// 			indice1 = i
// 		}
// 		if v == alt2 {
// 			indice2 = i
// 		}
// 	}
// 	return indice1 < indice2
// }

func MaxCount(count Count) (bestAlts []Alternative) { // renvoie les meilleures alternatives pour un décomtpe donné
	sortedAlts := make([]Alternative, len(count))
	for key := range count {
		sortedAlts = append(sortedAlts, key)
	}

	sort.SliceStable(sortedAlts, func(i, j int) bool {
		return count[sortedAlts[i]] < count[sortedAlts[j]]
	})

	bestAlts = append(bestAlts, sortedAlts[len(sortedAlts)-1])
	for i := len(sortedAlts) - 2; i > 0; i-- {
		if count[sortedAlts[i]] < count[sortedAlts[len(sortedAlts)-1]] {
			return
		}
		if count[sortedAlts[i]] == count[sortedAlts[len(sortedAlts)-1]] {
			bestAlts = append(bestAlts, sortedAlts[i])
		}
	}

	return
}

func CheckProfile(prefs Profile) error { // vérifie le profil donné, par ex. qu'ils sont tous complets et que chaque alternative n'apparaît qu'une seule fois par préférences
	len1 := len(prefs[0])
	alternatives := make(map[Alternative]bool, len1)
	newError := errors.New("profil incomplet")
	newErrorDoublon := errors.New(("doublons"))
	for _, pref := range prefs[0] {
		_, ok := alternatives[pref]
		if ok {
			return newErrorDoublon
		}
		alternatives[pref] = true
	}
	for _, profile := range prefs[1:] {
		if len(profile) != len1 {
			return newError
		}
		prof := make(map[Alternative]bool, len1)
		for _, alt := range profile {
			_, ok := alternatives[alt]
			if !ok {
				return newError
			}
			prof[alt] = true
		}

		for alt := range alternatives {
			_, ok := prof[alt]
			if !ok {
				return newError
			}
		}
	}
	return nil
}

func CheckProfileAlternative(prefs Profile, alts []Alternative) error { // vérifie le profil donné, par ex. qu'ils sont tous complets et que chaque alternative de alts apparaît exactement une fois par préférences
	newError := errors.New("alternative manquante")
	err := CheckProfile(prefs)
	if err != nil {
		return err
	}

	for _, profile := range prefs {
		prof := make(map[Alternative]bool, len(profile))
		for _, alt := range profile {
			prof[alt] = true
		}

		for _, alt := range alts {
			_, ok := prof[alt]
			if !ok {
				return newError
			}
		}
	}

	return nil
}
