package query

type CheckString func(v string) bool

func CheckStringNotEmpty(v string) bool {
	if v == "" {
		return false
	}
	return true
}

type stringValue struct {
	*baseValue

	p  *string
	cf CheckString
}

func NewStringValue(p *string, required bool, errno int, msg string, cf CheckString) *stringValue {
	s := &stringValue{
		baseValue: newBaseValue(required, errno, msg),

		p:  p,
		cf: cf,
	}

	return s
}

func (s *stringValue) Set(str string) error {
	*(s.p) = str

	return nil
}

func (s *stringValue) Check() bool {
	if s.cf == nil {
		return true
	}

	return s.cf(*(s.p))
}
