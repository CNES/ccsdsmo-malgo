/**
 * MIT License
 *
 * Copyright (c) 2017 - 2019 CNES
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
package tcp_test

import (
	"fmt"
	. "github.com/CNES/ccsdsmo-malgo/mal"
	"github.com/CNES/ccsdsmo-malgo/mal/transport/tcp" // Needed to initialize TCP transport factory
	"testing"
	"time"
)

var (
	maltcp string = "maltcp://127.0.0.1:16000?par1=value1"
)

// Test TCP transport using a unique context
func TestTCP1(t *testing.T) {
	ctx, err := NewContext(maltcp)
	if err != nil {
		t.Fatal("Error creating context, ", err)
	}

	consumer, err := NewEndPoint(ctx, "consumer", nil)
	if err != nil {
		t.Fatal("Error creating consumer, ", err)
	}
	fmt.Println("Registered consumer: ", consumer.Uri)

	provider, err := NewEndPoint(ctx, "provider", nil)
	if err != nil {
		t.Fatal("Error creating provider, ", err)
	}
	fmt.Println("Registered provider: ", provider.Uri)

	endpoint, err := ctx.GetEndPoint(ctx.NewURI("consumer"))
	fmt.Println("consumer: ", endpoint, err)

	nbmsg := 0
	go func() {
		var msg *Message = nil
		var err error = nil
		for err == nil {
			msg, err = provider.Recv()
			if msg != nil {
				par, err := msg.DecodeLastParameter(NullString, false)
				fmt.Println("receive: ", *par.(*String), ", ", err)
				nbmsg += 1
			}
		}
		t.Log("end: ", err)
	}()

	body := tcp.NewTCPBody(make([]byte, 0, 1024), true)
	body.EncodeLastParameter(NewString("message1"), false)
	msg1 := &Message{
		UriFrom:          consumer.Uri,
		UriTo:            provider.Uri,
		TransactionId:    consumer.TransactionId(),
		InteractionType:  MAL_INTERACTIONTYPE_SEND,
		InteractionStage: MAL_IP_STAGE_SEND,
		Body:             body,
	}
	consumer.Send(msg1)

	body = tcp.NewTCPBody(make([]byte, 0, 1024), true)
	body.EncodeLastParameter(NewString("message2"), false)
	msg2 := &Message{
		UriFrom:          consumer.Uri,
		UriTo:            provider.Uri,
		TransactionId:    consumer.TransactionId(),
		InteractionType:  MAL_INTERACTIONTYPE_SEND,
		InteractionStage: MAL_IP_STAGE_SEND,
		Body:             body,
	}
	consumer.Send(msg2)

	time.Sleep(1000 * time.Millisecond)
	ctx.Close()

	if nbmsg != 2 {
		t.Errorf("Receives %d messages, expect %d ", nbmsg, 2)
	}
	time.Sleep(1000 * time.Millisecond)
}

var (
	maltcp1 string = "maltcp://127.0.0.1:16001"
	maltcp2 string = "maltcp://127.0.0.1:16002"
)

// Test TCP transport using 2 different contexts
func TestTCP2(t *testing.T) {
	ctx1, err := NewContext(maltcp1)
	if err != nil {
		t.Fatal("Error creating context, ", err)
	}

	consumer, err := NewEndPoint(ctx1, "consumer", nil)
	if err != nil {
		t.Fatal("Error creating consumer, ", err)
	}
	fmt.Println("Registered consumer: ", *consumer.Uri)

	ctx2, err := NewContext(maltcp2)
	if err != nil {
		t.Fatal("Error creating context, ", err)
	}

	provider, err := NewEndPoint(ctx2, "provider", nil)
	if err != nil {
		t.Fatal("Error creating provider, ", err)
	}
	fmt.Println("Registered provider: ", *provider.Uri)

	endpoint, err := ctx1.GetEndPoint(ctx1.NewURI("consumer"))
	fmt.Println("consumer: ", endpoint, err)

	nbmsg := 0
	go func() {
		var msg *Message = nil
		var err error = nil
		for err == nil {
			msg, err = provider.Recv()
			if msg != nil {
				par, err := msg.DecodeLastParameter(NullString, false)
				fmt.Println("receive: ", *par.(*String), ", ", err)
				nbmsg += 1
			}
		}
		t.Log("end: ", err)
	}()

	body := tcp.NewTCPBody(make([]byte, 0, 1024), true)
	body.EncodeLastParameter(NewString("message1"), false)
	msg1 := &Message{
		UriFrom:          consumer.Uri,
		UriTo:            provider.Uri,
		TransactionId:    consumer.TransactionId(),
		InteractionType:  MAL_INTERACTIONTYPE_SEND,
		InteractionStage: MAL_IP_STAGE_SEND,
		Body:             body,
	}
	consumer.Send(msg1)

	body = tcp.NewTCPBody(make([]byte, 0, 1024), true)
	body.EncodeLastParameter(NewString("message2"), false)
	msg2 := &Message{
		UriFrom:          consumer.Uri,
		UriTo:            provider.Uri,
		TransactionId:    consumer.TransactionId(),
		InteractionType:  MAL_INTERACTIONTYPE_SEND,
		InteractionStage: MAL_IP_STAGE_SEND,
		Body:             body,
	}
	consumer.Send(msg2)

	time.Sleep(250 * time.Millisecond)
	ctx1.Close()
	ctx2.Close()

	if nbmsg != 2 {
		t.Errorf("Receives %d messages, expect %d ", nbmsg, 2)
	}
	time.Sleep(1000 * time.Millisecond)
}
