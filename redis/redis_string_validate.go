package redis

func ValidateSet(args [][]byte) bool {
	return true
}

// redis get 参数只能有一个
func ValidateGet(args [][]byte) bool {
	return len(args) == 1
}
