package transaction

import (
	"github.com/gofrs/uuid"
	"github.com/uptrace/bun"
	"time"
)

const (
	TypeIncome   = "INCOME"
	TypeWithdraw = "WITHDRAW"
)

type Transaction struct {
	bun.BaseModel `bun:"table:transactions,alias:tr"`

	ID          uuid.UUID `bun:"id,type:uuid,pk"             json:"-"`
	UserId      uuid.UUID `bun:"user_id,type:uuid"           json:"-"`
	OderNumber  string    `bun:"order,notnull"               json:"order"`
	Sum         int       `bun:"sum,notnull"                 json:"sum"`
	ProcessedAt time.Time `bun:"processed_at,notnull"        json:"processed_at"`
	Type        string    `bun:"type"                        json:"-"`
}
