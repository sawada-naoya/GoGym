package session

import (
	"context"

	"gogym-api/internal/adapter/http/dto"
	dom "gogym-api/internal/domain/user"
)

func (i *interactor) CreateSession(ctx context.Context, email string) (dto.TokenPairResponse, error) {

	emailObj, err := dom.NewEmail(email)
	if err != nil {
		return dto.TokenPairResponse{}, err
	}

	user, err := i.ur.FindByEmail(ctx, emailObj)
	if err != nil {
		return dto.TokenPairResponse{}, dom.NewDomainError("email_not_found")
	}
	if user == nil {
		return dto.TokenPairResponse{}, dom.NewDomainError( "user_not_found")
	}

	access, accessTTL, err := i.jwt.IssueAccess(user.ID)
	if err != nil {
		return dto.TokenPairResponse{}, err
	}

	refresh, _, _, exp, err := i.jwt.IssueRefresh(user.ID)
	if err != nil {
		return dto.TokenPairResponse{}, err
	}

	return dto.TokenPairResponse{
		AccessToken:  access,
		ExpiresIn:    int64(accessTTL.Seconds()),
		RefreshToken: refresh,
		RefreshExp:   exp.Unix(),
	}, nil
}
