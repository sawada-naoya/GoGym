package gym

type gymInteractor struct {
	repo Repository
}

func NewGymInteractor(repo Repository) GymUseCase {
	return &gymInteractor{
		repo: repo,
	}
}
