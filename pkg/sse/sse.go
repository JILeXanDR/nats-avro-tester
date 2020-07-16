package sse

type Hub interface {
	Run()
	Register(*client)
	Unregister(*client)
	NotifyAll(interface{})
}
