package gym

type DomainError struct {
	Key string
}

func (e *DomainError) Error() string { return e.Key }

func NewDomainError(key string) *DomainError {
	return &DomainError{Key: key}
}
