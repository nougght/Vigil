package model

import "time"

const (
	ActivityKindIdle   = 0
	ActivityKindActive = 1
	ActivityKindFocus  = 2
)

type ActivityUpdate struct {
	Kind      int32
	Title     *string
	Timestamp time.Time
}
