package entities

type Messager interface {
	Send(msg string)
}

type ProgressMessager interface {
	Messager
	Finish(msg string)
	Stop()
}
