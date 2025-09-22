package review

type interactor struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &interactor{
		repo: repo,
	}
}
