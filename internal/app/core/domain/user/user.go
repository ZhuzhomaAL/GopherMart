package user

import (
	"github.com/gofrs/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	*bun.BaseModel `bun:"table:users,alias:u"`

	ID       uuid.UUID `bun:"id,type:uuid,pk"             json:"id"`
	Login    string    `bun:"login,notnull,unique"        json:"login"`
	Password string    `bun:"password,notnull"            json:"password"`
}
