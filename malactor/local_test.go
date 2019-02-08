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
package malactor_test

import (
	"fmt"
	. "github.com/CNES/ccsdsmo-malgo/mal"
	_ "github.com/CNES/ccsdsmo-malgo/mal/transport/invm"
	_ "github.com/CNES/ccsdsmo-malgo/mal/transport/tcp"
	. "github.com/CNES/ccsdsmo-malgo/malactor"
	"testing"
	"time"
)

var (
	invm string = "invm://local"
)

func TestLocal(t *testing.T) {
	ctx, err := NewContext(invm)
	if err != nil {
		t.Fatal("Error creating context, ", err)
	}

	consumer, err := NewEndPoint(ctx, "consumer", nil)
	if err != nil {
		t.Fatal("Error creating consumer, ", err)
	}
	fmt.Println("Registered consumer: ", consumer.Uri)

	routing := new(Routing)
	//	provider := new(MyProvider)
	//	routing.registerProviderProgress(..., provider)
	actor, err := NewActor(ctx, "provider", routing)
	fmt.Println("Registered provider: ", actor)
	actor.Start()

	endpoint, err := ctx.GetEndPoint(ctx.NewURI("consumer"))
	fmt.Println("consumer: ", endpoint, err)

	time.Sleep(1 * time.Second)
	body := ctx.NewBody()
	body.EncodeLastParameter(NewString("message1"), false)
	msg1 := &Message{
		UriFrom: consumer.Uri,
		UriTo:   actor.Uri(),
		Body:    body,
	}
	consumer.Send(msg1)
	fmt.Println("Message sent: ", msg1.Body)

	time.Sleep(250 * time.Millisecond)
	body = ctx.NewBody()
	body.EncodeLastParameter(NewString("message2"), false)
	msg2 := &Message{
		UriFrom: consumer.Uri,
		UriTo:   actor.Uri(),
		Body:    body,
	}
	consumer.Send(msg2)
	fmt.Println("Message sent: ", msg2.Body)

	time.Sleep(250 * time.Millisecond)
	actor.Stop()

	time.Sleep(250 * time.Millisecond)
	ctx.Close()

	time.Sleep(250 * time.Millisecond)
}
