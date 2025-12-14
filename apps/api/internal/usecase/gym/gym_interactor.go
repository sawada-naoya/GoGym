package gym

type interactor struct {
	repo Repository
}

func NewInteractor(repo Repository) GymUseCase {
	return &interactor{
		repo: repo,
	}
}
