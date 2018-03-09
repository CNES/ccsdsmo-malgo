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
	. "github.com/ccsdsmo/malgo/src/mal"
	. "github.com/ccsdsmo/malgo/src/mal/api"
	_ "github.com/ccsdsmo/malgo/src/mal/transport/tcp" // Needed to initialize TCP transport factory
	"testing"
	"time"
)

const (
	nested2_consumer_url = "maltcp://127.0.0.1:16002"
	nested2_provider_url = "maltcp://127.0.0.1:16001"
)

// ########## ########## ########## ########## ########## ########## ########## ##########
// Test Invoke interaction in a nested scenario: Ctx1.C -> ctx2.P1 -> ctx2.P2
// Invoke a method in same service (same MAL context).

// Define Provider (Invoke interaction with nested invoke)

type TestNested2Provider1 struct {
	cctx  *ClientContext
	uri   *URI
	p2uri *URI
	nbmsg int
}

func newTestNested2Provider1(cctx *ClientContext, p2uri *URI) (*TestNested2Provider1, error) {
	provider := &TestNested2Provider1{cctx, cctx.Uri, p2uri, 0}

	// Register handler
	invokeHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			transaction := t.(InvokeTransaction)
			fmt.Println("\t$$$$$ Provider1 receive: ", string(msg.Body))
			transaction.Ack(nil, false)
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

			transaction.Reply(msg.Body, false)
		} else {
			fmt.Println("receive: nil")
		}
		return nil
	}
	cctx.RegisterInvokeHandler(200, 1, 1, 1, invokeHandler)

	return provider, nil
}

// Define Provider2 (Invoke interaction)

type TestNested2Provider2 struct {
	cctx  *ClientContext
	uri   *URI
	nbmsg int
}

func newTestNested2Provider2(cctx *ClientContext) (*TestNested2Provider2, error) {
	provider := &TestNested2Provider2{cctx, cctx.Uri, 0}

	// Register handler
	invokeHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			transaction := t.(InvokeTransaction)
			fmt.Println("\t$$$$$ Provider2 receive: ", string(msg.Body))
			transaction.Ack(nil, false)
			fmt.Println("\t$$$$$ Provider2 ack sent: OK")
			provider.nbmsg += 1
			time.Sleep(250 * time.Millisecond)
			//		transaction.Reply([]byte("reply message"), nil)
			transaction.Reply(msg.Body, false)
			fmt.Println("\t$$$$$ Provider2 reply sent: OK")
		} else {
			fmt.Println("receive: nil")
		}
		return nil
	}
	cctx.RegisterInvokeHandler(200, 1, 2, 2, invokeHandler)

	return provider, nil
}

// Test TCP transport Invoke Interaction using the high level API
func TestNested2(t *testing.T) {
	// Waits socket closing from previous test
	time.Sleep(250 * time.Millisecond)

	provider_ctx, err := NewContext(nested2_provider_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}
	defer provider_ctx.Close()

	cctx, err := NewClientContext(provider_ctx, "provider")
	if err != nil {
		t.Fatal("Error creating client context, ", err)
		return
	}

	provider2, err := newTestNested2Provider2(cctx)
	if err != nil {
		t.Fatal("Error creating provider, ", err)
		return
	}

	provider1, err := newTestNested2Provider1(cctx, provider2.cctx.Uri)
	if err != nil {
		t.Fatal("Error creating provider, ", err)
		return
	}

	consumer_ctx, err := NewContext(nested2_consumer_url)
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
