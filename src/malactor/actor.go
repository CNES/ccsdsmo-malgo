package malactor

import (
	"fmt"
	. "mal"
)

type Actor interface {
	Uri() *URI
	Start() error
	Stop() error
	Join() error
}

type actor struct {
	cmd      chan int
	routing  *Routing
	endpoint *EndPoint
	ch       chan Message
}

func NewActor(ctx *Context, service string, routing *Routing) (Actor, error) {
	actor := new(actor)
	actor.cmd = make(chan int)
	actor.routing = routing

	ch := make(chan *Message, 10)
	endpoint, err := NewEndPoint(ctx, service, ch)
	if err != nil {
		fmt.Println("Error creating handler, ", err)
		return nil, err
	}
	actor.endpoint = endpoint

	return actor, nil
}

func (actor *actor) Uri() *URI {
	return actor.endpoint.Uri
}

//func (actor *actor) Start() error {
//	go func() {
//		var msg *Message = nil
//		var err error = nil
//		for err == nil {
//			msg, err = actor.endpoint.Recv()
//			if msg != nil {
//				fmt.Println("receive: ", string(msg.Body), ", ", err)
//				actor.routing.Handle(msg)
//			}
//		}
//		fmt.Println("end: ", err)
//	}()
//	return nil
//}

func (actor *actor) Start() error {
	go func() {
		running := true
		for running {
			select {
			case msg, ok := <-actor.ch:
				if ok {
					fmt.Println("receive: ", string(msg.Body))
					actor.routing.Handle(&msg)
				} else {
					// TODO (AF):
					fmt.Println("eror on channel")
					running = false
				}
			case i := <-actor.cmd:
				fmt.Println("quit:", i)
				running = false
			}
		}
		fmt.Println("end: ")
	}()
	return nil
}

func (actor *actor) Stop() error {
	actor.cmd <- 0
	return nil
}

func (actor *actor) Join() error {
	return nil
}
