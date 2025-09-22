package user

type interactor struct {
	repo       Repository
	hasher     PasswordHasher
	idProvider IDProvider
}

func NewInteractor(repo Repository, hasher PasswordHasher, idProvider IDProvider) UseCase {
	return &interactor{
		repo:       repo,
		hasher:     hasher,
		idProvider: idProvider,
	}
}