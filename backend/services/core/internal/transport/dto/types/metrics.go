package dto

import (
	"time"

	"github.com/google/uuid"
)

type SeriesDTO struct {
	// TODO: replace kind value with key value
	Key    int32      `json:"key"`             // "cpu.percent"
	Label  string     `json:"label,omitempty"` // "C:", "eth0", "gpu0"
	Unit   string     `json:"unit"`            // "%", "bytes", "bps"
	Ts     []int64    `json:"ts,omitempty"`
	Values []*float64 `json:"values"`
} // @Name Series

type ActivityUpdate struct {
	AgentID   uuid.UUID `json:"agentID"`
	Kind      int32     `json:"kind"`
	Title     *string   `json:"title,omitempty"`
	Timestamp time.Time `json:"ts"`
}
