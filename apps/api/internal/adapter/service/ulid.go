package ulid

import (
	crypto "crypto/rand"
	"time"

	uc "gogym-api/internal/usecase/user"

	"github.com/oklog/ulid/v2"
)

type ULIDProvider struct{ ent *ulid.MonotonicEntropy }

func NewULIDProvider() *ULIDProvider { return &ULIDProvider{ent: ulid.Monotonic(crypto.Reader, 0)} }

func (p *ULIDProvider) NewUserID() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), p.ent).String()
}

var _ uc.IDProvider = (*ULIDProvider)(nil)
