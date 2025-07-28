package helper

func RemoveDuplicate(arr []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, v := range arr {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}

	return result
}
