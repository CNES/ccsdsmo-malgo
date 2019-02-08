/**
 * MIT License
 *
 * Copyright (c) 2018 - 2019 ccsdsmo
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */
package malactor

import (
	"fmt"
	. "github.com/CNES/ccsdsmo-malgo/mal"
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
					fmt.Println("receive: ", msg.Body)
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
