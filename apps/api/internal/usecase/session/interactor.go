package session

type interactor struct {
	repo Repository
}

func NewInteractor(repo Repository) *interactor {
	return &interactor{repo: repo}
}
