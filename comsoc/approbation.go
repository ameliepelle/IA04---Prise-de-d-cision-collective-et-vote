package procedures

import (
	"errors"
)

func ApprovalSWF(p Profile, thresholds []int) (count Count, err error) {
	nbAlt := len(p[0])
	count = make(Count, nbAlt)
	if len(p) > len(thresholds) {
		return count, errors.New("pas le bon nombre de seuils")
	}
	err = CheckProfile(p)

	if err != nil {
		return
	}

	for i, pref := range p {
		for j, alt := range pref {
			if j > thresholds[i]-1 {
				continue
			}
			count[alt] += 1
		}
	}
	return
}

func ApprovalSCF(p Profile, thresholds []int) (bestAlts []Alternative, err error) {
	var count Count
	count, err = ApprovalSWF(p, thresholds)
	bestAlts = MaxCount(count)
	return
}
