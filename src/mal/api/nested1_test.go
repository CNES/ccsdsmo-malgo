/**
 * MIT License
 *
 * Copyright (c) 2018 CNES
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
package api_test

import (
	"errors"
	"fmt"
	. "mal"
	. "mal/api"
	_ "mal/transport/tcp" // Needed to initialize TCP transport factory
	"testing"
	"time"
)

const (
	nested1_consumer_url  = "maltcp://127.0.0.1:16002"
	nested1_provider1_url = "maltcp://127.0.0.1:16001"
	nested1_provider2_url = "maltcp://127.0.0.1:16003"
)

// ########## ########## ########## ########## ########## ########## ########## ##########
// Test Invoke interaction in a nested scenario: Ctx1.C -> ctx2.P -> ctx3.P2
// Invoke a method to another service (different MAL context).

// Define Provider1 (Invoke interaction with nested invoke)

type Nested1Provider1 struct {
	cctx  *ClientContext
	p2uri *URI
	nbmsg int
}

func NewNested1Provider1(ctx *Context, service string, p2uri *URI) (*Nested1Provider1, error) {
	cctx, err := NewClientContext(ctx, service)
	if err != nil {
		return nil, err
	}
	provider := &Nested1Provider1{cctx: cctx, p2uri: p2uri}
	cctx.RegisterInvokeProvider(200, 1, 1, 1, provider)

	return provider, nil
}

func (provider *Nested1Provider1) OnInvoke(msg *Message, transaction InvokeTransaction) error {
	if msg != nil {
		fmt.Println("\t$$$$$ Provider1 receive: ", string(msg.Body))
		transaction.Ack(nil)
		provider.nbmsg += 1

		op := provider.cctx.NewInvokeOperation(provider.p2uri, 200, 1, 2, 2)
		_, err := op.Invoke([]byte("message from provider1"))
		if err != nil {
			//			t.Fatal("Error during invoke, ", err)
			return errors.New("Error during invoke")
		}

		reply, err := op.GetResponse()
		if err != nil {
			//			t.Fatal("Error getting response, ", err)
			return errors.New("Error getting response")
		}
		fmt.Println("\t&&&&& Nested Invoke1: OK, ", string(reply.Body))

		transaction.Reply(msg.Body, nil)
	} else {
		fmt.Println("receive: nil")
	}
	return nil
}

// Define Provider2 (Invoke interaction)

type Nested1Provider2 struct {
	cctx  *ClientContext
	nbmsg int
}

func NewNested1Provider2(ctx *Context, service string) (*Nested1Provider2, error) {
	cctx, err := NewClientContext(ctx, service)
	if err != nil {
		return nil, err
	}
	provider := &Nested1Provider2{cctx: cctx}
	cctx.RegisterInvokeProvider(200, 1, 2, 2, provider)

	return provider, nil
}

func (provider *Nested1Provider2) OnInvoke(msg *Message, transaction InvokeTransaction) error {
	if msg != nil {
		fmt.Println("\t$$$$$ Provider2 receive: ", string(msg.Body))
		transaction.Ack(nil)
		fmt.Println("\t$$$$$ Provider2 ack sent: OK")
		provider.nbmsg += 1
		time.Sleep(250 * time.Millisecond)
		//		transaction.Reply([]byte("reply message"), nil)
		transaction.Reply(msg.Body, nil)
		fmt.Println("\t$$$$$ Provider2 reply sent: OK")
	} else {
		fmt.Println("receive: nil")
	}
	return nil
}

// Test TCP transport Invoke Interaction using the high level API
func TestNested1Provider(t *testing.T) {
	// Waits socket closing from previous test
	time.Sleep(250 * time.Millisecond)

	provider2_ctx, err := NewContext(nested1_provider2_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}
	defer provider2_ctx.Close()

	provider2, err := NewNested1Provider2(provider2_ctx, "provider2")
	if err != nil {
		t.Fatal("Error creating provider, ", err)
		return
	}

	provider1_ctx, err := NewContext(nested1_provider1_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}
	defer provider1_ctx.Close()

	provider1, err := NewNested1Provider1(provider1_ctx, "provider", provider2.cctx.Uri)
	if err != nil {
		t.Fatal("Error creating provider, ", err)
		return
	}

	consumer_ctx, err := NewContext(nested1_consumer_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}
	defer consumer_ctx.Close()

	consumer, err := NewClientContext(consumer_ctx, "consumer")
	if err != nil {
		t.Fatal("Error creating consumer, ", err)
		return
	}

	op1 := consumer.NewInvokeOperation(provider1.cctx.Uri, 200, 1, 1, 1)
	_, err = op1.Invoke([]byte("message1"))
	if err != nil {
		t.Fatal("Error during invoke, ", err)
		return
	}

	r1, err := op1.GetResponse()
	if err != nil {
		t.Fatal("Error getting response, ", err)
		return
	}
	fmt.Println("\t&&&&& Invoke1: OK, ", string(r1.Body))

	op2 := consumer.NewInvokeOperation(provider1.cctx.Uri, 200, 1, 1, 1)
	_, err = op2.Invoke([]byte("message2"))
	if err != nil {
		t.Fatal("Error during invoke, ", err)
		return
	}

	r2, err := op2.GetResponse()
	if err != nil {
		t.Fatal("Error getting response, ", err)
		return
	}
	fmt.Println("\t&&&&& Invoke2: OK, ", string(r2.Body))

	if provider1.nbmsg != 2 {
		t.Errorf("Receives %d messages, expect %d ", provider1.nbmsg, 2)
	}

	if provider2.nbmsg != 2 {
		t.Errorf("Receives %d messages, expect %d ", provider2.nbmsg, 2)
	}
}
