package database

import (
	"time"

	"github.com/google/uuid"
)

type BaseDbModel struct {
	ID        uuid.UUID `bun:"id,pk"`
	CreatedAt time.Time `bun:",notnull,default:current_timestamp`
	UpdatedAt time.Time `bun:",notnull,default:current_timestamp`
}
