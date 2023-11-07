package changeStrings

import "strings"

func RemoveDuplicates(str string) string {
	// Use map to record duplicates as we find them.
	mapa := map[rune]bool{}
	for _, c := range str {
		_, ok := mapa[c]
		if !ok {
			count := strings.Count(str, string(c))
			if count == 1 {
				mapa[c] = true
			}
		}
	}
	keys := make([]rune, len(mapa))
	i := 0
	for k := range mapa {
		if mapa[k] {
			keys[i] = k
			i++
		}
	}
	return string(keys)
}

func ReturnDuplicates(str string) string {
	// Use map to record duplicates as we find them.
	mapa := map[rune]bool{}
	for _, c := range str {
		_, ok := mapa[c]
		if !ok {
			count := strings.Count(str, string(c))
			if count > 1 {
				mapa[c] = true
			}
		}
	}
	keys := make([]rune, len(mapa))
	i := 0
	for k := range mapa {
		if mapa[k] {
			keys[i] = k
			i++
		}
	}
	return string(keys)
}
