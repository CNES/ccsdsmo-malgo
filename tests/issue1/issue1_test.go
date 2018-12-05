/**
 * MIT License
 *
 * Copyright (c) 2017 - 2018 CNES
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
package issue1

import (
	"fmt"
	. "github.com/CNES/ccsdsmo-malgo/mal"
	. "github.com/CNES/ccsdsmo-malgo/mal/api"
	_ "github.com/CNES/ccsdsmo-malgo/mal/transport/tcp" // Needed to initialize TCP transport factory
	"testing"
	"time"
)

const (
	test1_provider_url = "maltcp://127.0.0.1:16001"
	test1_consumer_url = "maltcp://127.0.0.1:16002"
)

// ########## ########## ########## ########## ########## ########## ########## ##########
// Test Send interaction

// Define SendProvider

type TestSendProvider struct {
	ctx   *Context
	cctx  *ClientContext
	nbmsg int
}

func newTestSendProvider(ctx *Context) (*TestSendProvider, error) {
	cctx, err := NewClientContext(ctx, "sendProvider")
	if err != nil {
		return nil, err
	}
	provider := &TestSendProvider{ctx, cctx, 0}

	// Register handler
	sendHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			fmt.Println("\t$$$$$ sendHandler receive: ", string(msg.Body))
			provider.nbmsg += 1
		} else {
			fmt.Println("receive: nil")
		}
		return nil
	}
	cctx.RegisterSendHandler(200, 1, 1, 1, sendHandler)

	return provider, nil
}

func (provider *TestSendProvider) close() {
	provider.cctx.Close()
}

// ########## ########## ########## ########## ########## ########## ########## ##########
// Test Submit interaction

// Define SubmitProvider

type TestSubmitProvider struct {
	ctx   *Context
	cctx  *ClientContext
	nbmsg int
}

func newTestSubmitProvider(ctx *Context) (*TestSubmitProvider, error) {
	cctx, err := NewClientContext(ctx, "submitProvider")
	if err != nil {
		return nil, err
	}
	provider := &TestSubmitProvider{ctx, cctx, 0}

	// Register handler
	submitHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			transaction := t.(SubmitTransaction)
			fmt.Println("\t$$$$$ submitHandler receive: ", string(msg.Body))
			provider.nbmsg += 1
			transaction.Ack(nil, false)
		} else {
			fmt.Println("receive: nil")
		}
		return nil
	}
	cctx.RegisterSubmitHandler(200, 1, 1, 1, submitHandler)

	return provider, nil
}

func (provider *TestSubmitProvider) close() {
	provider.cctx.Close()
}

// Test TCP transport Send Interaction using the high level API
func Test1(t *testing.T) {
	// Waits socket closing from previous test
	time.Sleep(250 * time.Millisecond)

	provider_ctx, err := NewContext(test1_provider_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}
	defer provider_ctx.Close()

	sendProvider, err := newTestSendProvider(provider_ctx)
	if err != nil {
		t.Fatal("Error creating provider, ", err)
		return
	}
	defer sendProvider.close()

	submitProvider, err := newTestSubmitProvider(provider_ctx)
	if err != nil {
		t.Fatal("Error creating provider, ", err)
		return
	}
	defer submitProvider.close()

	consumer_ctx, err := NewContext(test1_consumer_url)
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

	op1 := consumer.NewSendOperation(sendProvider.cctx.Uri, 200, 1, 1, 1)
	op1.Send([]byte("message1"))

	op2 := consumer.NewSendOperation(sendProvider.cctx.Uri, 200, 1, 1, 1)
	op2.Send([]byte("message2"))

	// Waits for message reception
	time.Sleep(250 * time.Millisecond)

	if sendProvider.nbmsg != 2 {
		t.Errorf("Receives %d messages, expect %d ", sendProvider.nbmsg, 2)
	}

	// Test TCP transport Submit Interaction using the high level API
	op3 := consumer.NewSubmitOperation(submitProvider.cctx.Uri, 200, 1, 1, 1)
	_, err = op3.Submit([]byte("message1"))
	if err != nil {
		t.Fatal("Error during submit, ", err)
		return
	}
	fmt.Println("\t&&&&& Submit1: OK")

	op4 := consumer.NewSubmitOperation(submitProvider.cctx.Uri, 200, 1, 1, 1)
	_, err = op4.Submit([]byte("message2"))
	if err != nil {
		t.Fatal("Error during submit, ", err)
		return
	}
	fmt.Println("\t&&&&& Submit2: OK")

	if submitProvider.nbmsg != 2 {
		t.Errorf("Receives %d messages, expect %d ", submitProvider.nbmsg, 2)
	}
}
