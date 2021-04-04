package util

func SliceDiffStr(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	inter := SliceInterStr(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}
	for _, v := range slice1 {
		times, _ := m[v]
		if times == 0 {
			nn = append(nn, v)
		}
	}
	return nn
}

func SliceInterStr(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	for _, v := range slice1 {
		m[v]++
	}
	for _, v := range slice2 {
		times, _ := m[v]
		if times > 0 {
			nn = append(nn, v)
		}
	}
	return nn
}
