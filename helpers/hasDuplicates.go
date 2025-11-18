package helpers

func HasDuplicates(cats []string) bool {
	for i := 0; i < len(cats); i++ {
		for j := i + 1; j < len(cats); j++ {
			if cats[i] == cats[j] {
				return true
			}
		}
	}
	return false
}
