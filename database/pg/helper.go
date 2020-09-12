package pg

func ToGoBool(v int8) bool {
	if v == 0 {
		return false
	}
	return true
}
func ToPgBool(v bool) int8 {
	if v {
		return 1
	}
	return 0
}
