package algo_helper

func DedupStringList(ss []string) []string {
	smap := make(map[string]interface{})
	for _, s := range ss {
		smap[s] = nil
	}
	ssNew := make([]string, 0, len(smap))
	for s := range smap {
		ssNew = append(ssNew, s)
	}
	return ssNew
}
