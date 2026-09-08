package dto

type SeriesDTO struct {
	Key    int32      `json:"key"`             // "cpu.percent"
	Label  string     `json:"label,omitempty"` // "C:", "eth0", "gpu0"
	Unit   string     `json:"unit"`            // "%", "bytes", "bps"
	Ts     []int64    `json:"ts,omitempty"`
	Values []*float64 `json:"values"`
} // @Name Series
