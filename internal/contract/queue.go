package contract

type QueueContract interface {
	Name() string
	InjectPayload(payload map[string]interface{}) error
	Handle() error
}
