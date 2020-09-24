/**
 * MIT License
 *
 * Copyright (c) 2020 CNES
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
	"fmt"
	. "github.com/CNES/ccsdsmo-malgo/mal"
	. "github.com/CNES/ccsdsmo-malgo/mal/api"
	_ "github.com/CNES/ccsdsmo-malgo/mal/transport/tcp" // Needed to initialize TCP transport factory
	"testing"
	"time"
)

const (
	reentrant_url = "maltcp://127.0.0.1:16001"
)

// ########## ########## ########## ########## ########## ########## ########## ##########
// Test Submit interaction

// Define SubmitProvider

type TestReentrantSubmitProvider struct {
	ctx   *Context
	cctx  *ClientContext
	nbmsg int
}

func newTestReentrantSubmitProvider(ctx *Context) (*TestSubmitProvider, error) {
	cctx, err := NewClientContext(ctx, "provider")
	if err != nil {
		return nil, err
	}
	provider := &TestSubmitProvider{ctx, cctx, 0}

	// Register handler
	submitHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			transaction := t.(SubmitTransaction)
			par, err := msg.DecodeLastParameter(NullString, false)
			fmt.Println("\t$$$$$ submitHandler receive: ", *par.(*String), err)
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

// Test TCP transport Submit Interaction using the high level API
func TestReentrantSubmit(t *testing.T) {
	// Waits socket closing from previous test
	time.Sleep(250 * time.Millisecond)

	ctx, err := NewContext(reentrant_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}
	defer ctx.Close()

	provider, err := newTestReentrantSubmitProvider(ctx)
	if err != nil {
		t.Fatal("Error creating provider, ", err)
		return
	}

	consumer, err := NewClientContext(ctx, "consumer")
	if err != nil {
		t.Fatal("Error creating consumer, ", err)
		return
	}

	body := ctx.NewBody()
	body.EncodeLastParameter(NewString("message1"), false)
	op1 := consumer.NewSubmitOperation(provider.cctx.Uri, 200, 1, 1, 1)
	_, err = op1.Submit(body)
	if err != nil {
		t.Fatal("Error during submit, ", err)
		return
	}
	fmt.Println("\t&&&&& Submit1: OK")

	body.Reset(true)
	body.EncodeLastParameter(NewString("message2"), false)
	op2 := consumer.NewSubmitOperation(provider.cctx.Uri, 200, 1, 1, 1)
	_, err = op2.Submit(body)
	if err != nil {
		t.Fatal("Error during submit, ", err)
		return
	}
	fmt.Println("\t&&&&& Submit2: OK")

	if provider.nbmsg != 2 {
		t.Errorf("Receives %d messages, expect %d ", provider.nbmsg, 2)
	}
}
