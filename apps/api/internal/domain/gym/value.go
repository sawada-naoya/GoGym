// internal/domain/gym/value.go
// 役割: ジムドメインのバリューオブジェクト
// Equipment、PriceRange等のジム関連のバリューオブジェクトの定義と検証
package gym

// Equipment represents gym equipment as a value object
type Equipment struct {
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description" validate:"max=500"`
	Available   bool   `json:"available"`
}

// NewEquipment creates a new equipment value object
func NewEquipment(name, description string, available bool) Equipment {
	return Equipment{
		Name:        name,
		Description: description,
		Available:   available,
	}
}

// PriceRange represents gym pricing as a value object
type PriceRange struct {
	Min      int    `json:"min" validate:"min=0"`
	Max      int    `json:"max" validate:"min=0"`
	Currency string `json:"currency" validate:"required,len=3"`
}

// NewPriceRange creates a new price range value object
func NewPriceRange(min, max int, currency string) (*PriceRange, error) {
	if min < 0 || max < 0 {
		return nil, &ValueError{Field: "price", Message: "price cannot be negative"}
	}
	
	if min > max {
		return nil, &ValueError{Field: "price", Message: "min price cannot be greater than max price"}
	}
	
	if len(currency) != 3 {
		return nil, &ValueError{Field: "currency", Message: "currency must be 3 characters"}
	}

	return &PriceRange{
		Min:      min,
		Max:      max,
		Currency: currency,
	}, nil
}

// IsValid validates the price range
func (pr PriceRange) IsValid() bool {
	return pr.Min >= 0 && pr.Max >= 0 && pr.Min <= pr.Max && len(pr.Currency) == 3
}

// ValueError represents value object validation error
type ValueError struct {
	Field   string
	Message string
}

func (e *ValueError) Error() string {
	return e.Field + ": " + e.Message
}