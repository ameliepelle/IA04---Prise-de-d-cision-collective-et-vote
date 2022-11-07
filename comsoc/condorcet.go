package procedures

func CondorcetWinner(p Profile) (bestAlts []Alternative, err error) {

	err = CheckProfile(p)

	if err != nil {
		return
	}

	//Créer une map pour compter le nombre de fois où chaque alternative gagne
	count := make(Count, len(p[0]))

	//On parcours les alternatives 2 à 2
	for i := range p[0] {
		for j := i + 1; j < len(p[0]); j++ {
			scoreI := 0
			scoreJ := 0
			for k := range p {
				if IsPref(p[0][i], p[0][j], p[k]) {
					scoreI = scoreI + 1
				} else {
					scoreJ = scoreJ + 1
				}
			}
			if scoreI > scoreJ {
				count[p[0][i]] += 1
			} else {
				count[p[0][j]] += 1
			}
		}

	}
	for k := range count {
		//Si len(alt)==score(k) k est un gagnat de condorcet
		if count[k] == len(p[0])-1 {
			bestAlts = append(bestAlts, k)
		}
	}
	return bestAlts, CheckProfile(p)
}
