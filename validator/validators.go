package validator

type XTCEValidator interface {
	Validate(filename string) error
}
