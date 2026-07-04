package realtime

type Event struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}
