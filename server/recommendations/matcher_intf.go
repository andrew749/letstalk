package recommendations

type Matcher interface {
	Match(matches []UserMatch) error
}
