package logger

import "fmt"

func collectArgs(args ...interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	if len(args) == 0 {
		return result
	}
	var badKeyCount int
	for i := 0; i < len(args); i++ {
		switch v := args[i].(type) {
		case string:
			if i+1 < len(args) {
				result[v] = args[i+1]
				i++
			} else {
				badKeyCount++
				result[fmt.Sprintf("bad_key_%d", badKeyCount)] = v
			}
		default:
			badKeyCount++
			result[fmt.Sprintf("bad_key_%d", badKeyCount)] = v
		}
	}
	return result
}
