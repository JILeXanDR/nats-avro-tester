package sse

type client struct {
	id     string
	data   chan interface{}
	closed chan struct{}
}

func NewClient(id string) *client {
	return &client{
		id:     id,
		data:   make(chan interface{}),
		closed: make(chan struct{}, 1),
	}
}

func (c *client) Wait(next func(interface{})) error {
	for {
		select {
		case <-c.closed:
			return nil
		case v, ok := <-c.data:
			if ok {
				next(v)
			}
		}
	}
}

func (c *client) ID() string {
	return c.id
}
