package procedures

func BordaSWF(p Profile) (count Count, err error) {
	err = CheckProfile(p)

	if err != nil {
		return
	}
	nbAlt := len(p[0])
	count = make(Count, nbAlt)
	for _, pref := range p {
		for i, alt := range pref {
			count[alt] += len(pref) - i - 1
		}
	}
	return
}

func BordaSCF(p Profile) (bestAlts []Alternative, err error) {
	var count Count
	count, err = BordaSWF(p)
	bestAlts = MaxCount(count)
	return
}
