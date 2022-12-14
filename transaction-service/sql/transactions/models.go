// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package transactions

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID          uuid.UUID
	Amount      int64
	ToAccount   uuid.UUID
	FromAccount uuid.UUID
	CreatedAt   time.Time
}
