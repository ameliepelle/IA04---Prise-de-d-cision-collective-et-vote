package procedures

import "gitlab.utc.fr/lagruesy/ia04/utils"

func TauDeKendall(alt1, alt2 []Alternative) (tau int, err error) {
	err = CheckProfile([][]Alternative{alt1, alt2})
	if err != nil {
		return
	}
	couples := make(map[Alternative]Alternative)
	perm := utils.FirstPermutation(len(alt1))
	for i := 0; i < len(alt1); i++ {
		if IsPref(Alternative(perm[i]), Alternative(perm[i+1]), alt1) != IsPref(Alternative(perm[i]), Alternative(perm[i+1]), alt2) {
			couples[Alternative(perm[i])] = Alternative(perm[i+1])
		}

	}
	perm, ok := utils.NextPermutation(perm)
	for ok {
		for i := 0; i < len(alt1); i++ {
			if IsPref(Alternative(perm[i]), Alternative(perm[i+1]), alt1) != IsPref(Alternative(perm[i]), Alternative(perm[i+1]), alt2) {
				if couples[Alternative(perm[i])] != Alternative(perm[i+1]) && couples[Alternative(perm[i+1])] != Alternative(perm[i]) {
					couples[Alternative(perm[i])] = Alternative(perm[i+1])
				}
			}
		}
		perm, ok = utils.NextPermutation(perm)
	}
	tau = len(couples)
	return
} // a tester

// func DistanceRangementProfil()
