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
package api2_test

import (
	"fmt"
	. "mal"
	. "mal/api2"
	_ "mal/transport/tcp" // Needed to initialize TCP transport factory
	"testing"
	"time"
)

const (
	provider_url = "maltcp://127.0.0.1:16001"
	consumer_url = "maltcp://127.0.0.1:16002"
)

// ########## ########## ########## ########## ########## ########## ########## ##########
// Test Send interaction

// Define SendProvider

type SendProvider struct {
	pctx  *ProviderContext
	nbmsg int
}

func NewSendProvider(ctx *Context, service string) (*SendProvider, error) {
	pctx, err := NewProviderContext(ctx, service)
	if err != nil {
		return nil, err
	}
	sendProvider := &SendProvider{pctx: pctx}
	pctx.RegisterSendHandler(200, 1, 1, 1, sendProvider)

	return sendProvider, nil
}

func (provider *SendProvider) OnSend(msg *Message, transaction SendTransaction) error {
	if msg != nil {
		fmt.Println("\n\n $$$$$ sendHandler receive: ", string(msg.Body), "\n\n")
		provider.nbmsg += 1
	} else {
		fmt.Println("receive: nil")
	}
	return nil
}

// Test TCP transport Send Interaction using the high level API
func TestSend(t *testing.T) {
	provider_ctx, err := NewContext(provider_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	provider, err := NewSendProvider(provider_ctx, "provider")
	if err != nil {
		t.Fatal("Error creating provider, ", err)
		return
	}

	consumer_ctx, err := NewContext(consumer_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	consumer, err := NewOperationContext(consumer_ctx, "consumer")
	if err != nil {
		t.Fatal("Error creating consumer, ", err)
		return
	}

	op1, err := consumer.NewSendOperation(200, 1, 1, 1)
	msg1 := &Message{
		UriTo: provider.pctx.Uri,
		Body:  []byte("message1"),
	}
	op1.Send(msg1)

	op2, err := consumer.NewSendOperation(200, 1, 1, 1)
	msg2 := &Message{
		UriTo: provider.pctx.Uri,
		Body:  []byte("message2"),
	}
	op2.Send(msg2)

	time.Sleep(250 * time.Millisecond)
	provider_ctx.Close()
	consumer_ctx.Close()

	if provider.nbmsg != 2 {
		t.Errorf("Receives %d messages, expect %d ", provider.nbmsg, 2)
	}

	// Waits for socket close
	time.Sleep(250 * time.Millisecond)
}

// ########## ########## ########## ########## ########## ########## ########## ##########
// Test Submit interaction

// Define SubmitProvider

type SubmitProvider struct {
	pctx  *ProviderContext
	nbmsg int
}

func NewSubmitProvider(ctx *Context, service string) (*SubmitProvider, error) {
	pctx, err := NewProviderContext(ctx, service)
	if err != nil {
		return nil, err
	}
	submitProvider := &SubmitProvider{pctx: pctx}
	pctx.RegisterSubmitHandler(200, 1, 1, 1, submitProvider)

	return submitProvider, nil
}

func (provider *SubmitProvider) OnSubmit(msg *Message, transaction SubmitTransaction) error {
	if msg != nil {
		fmt.Println("\n\n $$$$$ submitHandler receive: ", string(msg.Body), "\n\n")
		provider.nbmsg += 1
		transaction.Ack(nil)
	} else {
		fmt.Println("receive: nil")
	}
	return nil
}

// Test TCP transport Submit Interaction using the high level API
func TestSubmit(t *testing.T) {
	provider_ctx, err := NewContext(provider_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	provider, err := NewSubmitProvider(provider_ctx, "provider")
	if err != nil {
		t.Fatal("Error creating provider, ", err)
		return
	}

	consumer_ctx, err := NewContext(consumer_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	consumer, err := NewOperationContext(consumer_ctx, "consumer")
	if err != nil {
		t.Fatal("Error creating consumer, ", err)
		return
	}

	op1, err := consumer.NewSubmitOperation(200, 1, 1, 1)
	msg1 := &Message{
		UriTo: provider.pctx.Uri,
		Body:  []byte("message1"),
	}
	op1.Submit(msg1)

	fmt.Println("\n\n &&&&& Submit1: OK\n\n")

	op2, err := consumer.NewSubmitOperation(200, 1, 1, 1)
	msg2 := &Message{
		UriTo: provider.pctx.Uri,
		Body:  []byte("message2"),
	}
	op2.Submit(msg2)

	fmt.Println("\n\n &&&&& Submit2: OK\n\n")

	time.Sleep(250 * time.Millisecond)
	provider_ctx.Close()
	consumer_ctx.Close()

	if provider.nbmsg != 2 {
		t.Errorf("Receives %d messages, expect %d ", provider.nbmsg, 2)
	}

	// Waits for socket close
	time.Sleep(250 * time.Millisecond)
}

// ########## ########## ########## ########## ########## ########## ########## ##########
// Test Request interaction

// Define RequestProvider

type RequestProvider struct {
	pctx  *ProviderContext
	nbmsg int
}

func NewRequestProvider(ctx *Context, service string) (*RequestProvider, error) {
	pctx, err := NewProviderContext(ctx, service)
	if err != nil {
		return nil, err
	}
	requestProvider := &RequestProvider{pctx: pctx}
	pctx.RegisterRequestHandler(200, 1, 1, 1, requestProvider)

	return requestProvider, nil
}

func (provider *RequestProvider) OnRequest(msg *Message, transaction RequestTransaction) error {
	if msg != nil {
		fmt.Println("\n\n $$$$$ requestHandler receive: ", string(msg.Body), "\n\n")
		provider.nbmsg += 1
		transaction.Reply([]byte("reply message"), nil)
	} else {
		fmt.Println("receive: nil")
	}
	return nil
}

// Test TCP transport Request Interaction using the high level API
func TestRequest(t *testing.T) {
	provider_ctx, err := NewContext(provider_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	provider, err := NewRequestProvider(provider_ctx, "provider")
	if err != nil {
		t.Fatal("Error creating provider, ", err)
		return
	}

	consumer_ctx, err := NewContext(consumer_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	consumer, err := NewOperationContext(consumer_ctx, "consumer")
	if err != nil {
		t.Fatal("Error creating consumer, ", err)
		return
	}

	op1, err := consumer.NewRequestOperation(200, 1, 1, 1)
	msg1 := &Message{
		UriTo: provider.pctx.Uri,
		Body:  []byte("message1"),
	}
	op1.Request(msg1)

	fmt.Println("\n\n &&&&& Request1: OK\n\n")

	op2, err := consumer.NewRequestOperation(200, 1, 1, 1)
	msg2 := &Message{
		UriTo: provider.pctx.Uri,
		Body:  []byte("message2"),
	}
	op2.Request(msg2)

	fmt.Println("\n\n &&&&& Request2: OK\n\n")

	time.Sleep(250 * time.Millisecond)
	provider_ctx.Close()
	consumer_ctx.Close()

	if provider.nbmsg != 2 {
		t.Errorf("Receives %d messages, expect %d ", provider.nbmsg, 2)
	}

	// Waits for socket close
	time.Sleep(250 * time.Millisecond)
}

// ########## ########## ########## ########## ########## ########## ########## ##########
// Test Invoke interaction

// Define InvokeProvider

type InvokeProvider struct {
	pctx  *ProviderContext
	nbmsg int
}

func NewInvokeProvider(ctx *Context, service string) (*InvokeProvider, error) {
	pctx, err := NewProviderContext(ctx, service)
	if err != nil {
		return nil, err
	}
	invokeProvider := &InvokeProvider{pctx: pctx}
	pctx.RegisterInvokeHandler(200, 1, 1, 1, invokeProvider)

	return invokeProvider, nil
}

func (provider *InvokeProvider) OnInvoke(msg *Message, transaction InvokeTransaction) error {
	if msg != nil {
		fmt.Println("\n\n $$$$$ invokeHandler receive: ", string(msg.Body), "\n\n")
		transaction.Ack(nil)
		provider.nbmsg += 1
		time.Sleep(250 * time.Millisecond)
		transaction.Reply([]byte("reply message"), nil)
	} else {
		fmt.Println("receive: nil")
	}
	return nil
}

// Test TCP transport Invoke Interaction using the high level API
func TestInvoke(t *testing.T) {
	provider_ctx, err := NewContext(provider_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	provider, err := NewInvokeProvider(provider_ctx, "provider")
	if err != nil {
		t.Fatal("Error creating provider, ", err)
		return
	}

	consumer_ctx, err := NewContext(consumer_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	consumer, err := NewOperationContext(consumer_ctx, "consumer")
	if err != nil {
		t.Fatal("Error creating consumer, ", err)
		return
	}

	op1, err := consumer.NewInvokeOperation(200, 1, 1, 1)
	msg1 := &Message{
		UriTo: provider.pctx.Uri,
		Body:  []byte("message1"),
	}
	op1.Invoke(msg1)

	r1, err := op1.GetResponse()
	fmt.Println("\n\n &&&&& Invoke1: OK\n", string(r1.Body), "\n\n")

	op2, err := consumer.NewInvokeOperation(200, 1, 1, 1)
	msg2 := &Message{
		UriTo: provider.pctx.Uri,
		Body:  []byte("message2"),
	}
	op2.Invoke(msg2)

	r2, err := op1.GetResponse()
	fmt.Println("\n\n &&&&& Invoke2: OK\n", string(r2.Body), "\n\n")

	time.Sleep(250 * time.Millisecond)
	provider_ctx.Close()
	consumer_ctx.Close()

	if provider.nbmsg != 2 {
		t.Errorf("Receives %d messages, expect %d ", provider.nbmsg, 2)
	}

	// Waits for socket close
	time.Sleep(250 * time.Millisecond)
}

// ########## ########## ########## ########## ########## ########## ########## ##########
// Test Progress interaction

// Define ProgressProvider

type ProgressProvider struct {
	pctx *ProviderContext
}

func NewProgressProvider(ctx *Context, service string) (*ProgressProvider, error) {
	pctx, err := NewProviderContext(ctx, service)
	if err != nil {
		return nil, err
	}
	progressProvider := &ProgressProvider{pctx: pctx}
	pctx.RegisterProgressHandler(200, 1, 1, 1, progressProvider)

	return progressProvider, nil
}

func (provider *ProgressProvider) OnProgress(msg *Message, transaction ProgressTransaction) error {
	if msg != nil {
		fmt.Println("\n\n $$$$$ progressHandler receive: ", string(msg.Body), "\n\n")
		transaction.Ack(nil)
		for i := 0; i < 10; i++ {
			transaction.Update([]byte(fmt.Sprintf("messsage#%d", i)), nil)
		}
		transaction.Reply([]byte("last message"), nil)
	} else {
		fmt.Println("receive: nil")
	}
	return nil
}

// Test TCP transport Progress Interaction using the high level API
func TestProgress(t *testing.T) {
	provider_ctx, err := NewContext(provider_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	provider, err := NewProgressProvider(provider_ctx, "provider")
	if err != nil {
		t.Fatal("Error creating provider, ", err)
		return
	}

	consumer_ctx, err := NewContext(consumer_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	consumer, err := NewOperationContext(consumer_ctx, "consumer")
	if err != nil {
		t.Fatal("Error creating consumer, ", err)
		return
	}

	nbmsg := 0

	op1, err := consumer.NewProgressOperation(200, 1, 1, 1)
	msg1 := &Message{
		UriTo: provider.pctx.Uri,
		Body:  []byte("message1"),
	}
	op1.Progress(msg1)
	fmt.Println("\n\n &&&&& Progress1: OK\n\n")

	updt, err := op1.GetUpdate()
	if err != nil {
		t.Error(err)
	}
	for updt != nil {
		nbmsg += 1
		fmt.Println("\n\n &&&&& Progress1: Update -> ", string(updt.Body), "\n\n")
		updt, err = op1.GetUpdate()
		if err != nil {
			t.Error(err)
		}
	}
	rep, err := op1.GetResponse()
	if err != nil {
		t.Error(err)
	}
	nbmsg += 1
	fmt.Println("\n\n &&&&& Progress1: Response -> ", string(rep.Body), "\n\n")

	time.Sleep(250 * time.Millisecond)
	provider_ctx.Close()
	consumer_ctx.Close()

	if nbmsg != 11 {
		t.Errorf("Receives %d messages, expect %d ", nbmsg, 2)
	}

	// Waits for socket close
	time.Sleep(250 * time.Millisecond)
}
