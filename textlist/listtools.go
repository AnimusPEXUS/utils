package textlist

func RemoveDuplicatedStrings(in []string) []string {

	ret := make([]string, 0)

search:
	for _, i := range in {
		for _, j := range ret {
			if j == i {
				continue search
			}
		}
		ret = append(ret, i)
	}

	return ret
}

func RemoveZeroLengthItems(in []string) []string {

	ret := make([]string, 0)

	for _, i := range in {
		if len(i) != 0 {
			ret = append(ret, i)
		}
	}

	return ret
}
