package contract

type QueueContract interface {
	Name() string
	Handle(payload map[string]interface{}) error
}
