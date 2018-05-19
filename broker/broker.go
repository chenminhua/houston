package broker

type Broker interface {
	ReceiveMessage()

	Subscribe()
}