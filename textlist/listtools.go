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
