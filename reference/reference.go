package reference

type Reference interface {
	Reference() string
}

func New(text string) Reference {
	return &referenceString{text}
}

type referenceString struct {
	s string
}

func (r *referenceString) Reference() string {
	return r.s
}
