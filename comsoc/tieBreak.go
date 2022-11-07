package procedures

import (
	"errors"
	"fmt"
	"math/rand"
)

func TieBreak(alternatives []Alternative) (alt Alternative, err error) {
	err = errors.New("slice vide")
	if alternatives == nil {
		return
	}

	err = nil
	alt = alternatives[rand.Intn(len(alternatives))]
	return
}

// la fonction TieBreakFactory prend en argument un tri un ordre ex : [5, 3, 2, 1] le n°5 gagne en cas de tiebreak puis 3 etc
// puis la fonction renvoie une fonction qui fait le tiebreak en fonction de l'ordre indiqué
func TieBreakFactory(ordre []Alternative) func(alts []Alternative) (alt Alternative, err error) {
	return func(alts []Alternative) (alt Alternative, err error) {
		err = errors.New("slice vide")
		if alts == nil {
			return 0, err
		}
		err = nil
		var pref []Alternative
		var is int
		is = 0
		for _, e := range ordre {
			for _, a := range alts {
				if e == a {
					pref = append(pref, e)
					is = 1
					break
				}
			}
			if is == 1 {
				break
			}
		}
		fmt.Println("apres tiebreak", pref)
		return pref[0], err

	}

}

func SWFFactory(sfw func(p Profile) (Count, error), f func([]Alternative) (alt Alternative, err error)) func(Profile) ([]Alternative, error) {
	return func(p2 Profile) (alternatives []Alternative, err error) {
		count, err := sfw(p2)
		//on va chercher les meilleurs alternatives
		for {
			best := MaxCount(count)
			alt, err2 := f(best)
			fmt.Println(err2.Error())
			alternatives = append(alternatives, alt)
			delete(count, alt)
			if count == nil {
				break
			}
		}
		return alternatives, err

	}

}
