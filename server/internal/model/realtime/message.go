package realtime_model

import (
	"encoding/json"

	"github.com/google/uuid"
)

type MessageType string

const (
	MessageTypeSeries MessageType = "series"
)

type Message struct {
	Type    MessageType     `json:"type"`
	AgentID uuid.UUID       `json:"agentID"`
	Payload json.RawMessage `json:"payload"`
}

type ClientMessageType string

const (
	ClientMessageTypeAgentDetailed ClientMessageType = "agent.detailed"
)

type AgentDetailedMessage struct {
	Agents []uuid.UUID `json:"agents"`
}

type ClientMessage struct {
	Type    ClientMessageType `json:"type"`
	Payload json.RawMessage   `json:"payload"`
}
