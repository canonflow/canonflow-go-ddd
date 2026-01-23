package contract

type ProducerContract interface {
	GetTopic() *string
	Send(event Event) error
}
