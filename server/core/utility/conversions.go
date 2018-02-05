package utility

func GenderIdByName(gender string) int {
	switch gender {
	case "FEMALE":
		return 1
	case "MALE":
		return 2
	default:
		return 0
	}
}
