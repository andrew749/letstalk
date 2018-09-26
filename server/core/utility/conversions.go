package utility

func GenderIdByName(gender string) int {
	switch gender {
	case "FEMALE":
		fallthrough
	case "female":
		return 1
	case "MALE":
		fallthrough
	case "male":
		return 2
	default:
		return 3
	}
}
