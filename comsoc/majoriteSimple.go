package procedures

// func MajoritySWF(p Profile) (count Count, err error)            // décompte à partir d'un profil
// func MajoritySCF(p Profile) (bestAlts []Alternative, err error) // alternatives préférées

func MajoritySWF(p Profile) (count Count, err error) {
	err = CheckProfile(p)

	if err != nil {
		return
	}

	count = make(Count, len(p[0]))
	for _, pref := range p {
		count[pref[0]] += 1
	}

	return
}

func MajoritySCF(p Profile) (bestAlts []Alternative, err error) {
	var count Count
	count, err = MajoritySWF(p)
	bestAlts = MaxCount(count)
	return
}
