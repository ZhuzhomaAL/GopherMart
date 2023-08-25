package order

import (
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/gofrs/uuid"
	"github.com/uptrace/bun"
	"time"
)

const (
	StatusNew        string = "NEW"
	StatusProcessing string = "PROCESSING"
	StatusInvalid    string = "INVALID"
	StatusProcessed  string = "PROCESSED"
)

type Order struct {
	bun.BaseModel `bun:"table:orders,alias:o"`

	ID         uuid.UUID `bun:"id,type:uuid,pk"             json:"-"`
	UserId     uuid.UUID `bun:"user_id,type:uuid"           json:"-"`
	Number     string    `bun:"number,notnull,unique"       json:"number"`
	Status     string    `bun:"status,notnull"             json:"status"`
	UploadedAt time.Time `bun:"uploaded_at,notnull"         json:"uploaded_at"`
}

func ValidateOrderFormat(orderNumber string) bool {
	if err := goluhn.Validate(orderNumber); err != nil {
		return false
	}

	return true
}
