package dpds

type DotConsumer interface {
	Init()                                    // Initialize the consumer.
	Prepare() error                           // Prepare for consumption.
	Consume(params ...interface{}) error      // Consumes dot data.
	Commit() error                            // Commit all changes.
	Abort() bool                              // Cleanup and shut down.
	Finalize() bool                           // Cleanup and shut down.
}

type ConsumableDot struct {
	id       uint64
	parentId uint64
	name     string
	value    string
}

type DotConsumerFactory struct {
	dc     DotConsumer // Dot Consumer interface
}

func (dcf DotConsumerFactory) GetInstance() DotConsumer {
	return nil
}

var dcf DotConsumerFactory

func GetConsumerInstance() DotConsumer {
	return dcf.GetInstance()
}
