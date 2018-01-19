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
package api_test

import (
	"fmt"
	. "mal"
	. "mal/api"
	_ "mal/transport/tcp" // Needed to initialize TCP transport factory
	"testing"
	"time"
)

//const (
//	provider_url = "maltcp://127.0.0.1:16001"
//	consumer_url = "maltcp://127.0.0.1:16002"
//	subscriber_url = "maltcp://127.0.0.1:16001"
//	publisher_url  = "maltcp://127.0.0.1:16002"
//	broker_url     = "maltcp://127.0.0.1:16003"
//)

// ########## ########## ########## ########## ########## ########## ########## ##########
// Test Send interaction

// Define SendProvider

type MySendProvider struct {
	pctx  *ProviderContext
	nbmsg int
}

func NewSendProvider(ctx *Context, service string) (*MySendProvider, error) {
	pctx, err := NewProviderContext(ctx, service)
	if err != nil {
		return nil, err
	}
	sendProvider := &MySendProvider{pctx: pctx}
	pctx.RegisterSendProvider(200, 1, 1, 1, sendProvider)

	return sendProvider, nil
}

func (provider *MySendProvider) OnSend(msg *Message, transaction SendTransaction) error {
	if msg != nil {
		fmt.Println("\t$$$$$ sendHandler receive: ", string(msg.Body))
		provider.nbmsg += 1
	} else {
		fmt.Println("receive: nil")
	}
	return nil
}

// Test TCP transport Send Interaction using the high level API
func TestSendProvider(t *testing.T) {
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

	op1, err := consumer.NewSendOperation(provider.pctx.Uri, 200, 1, 1, 1)
	op1.Send([]byte("message1"))

	op2, err := consumer.NewSendOperation(provider.pctx.Uri, 200, 1, 1, 1)
	op2.Send([]byte("message2"))

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

type MySubmitProvider struct {
	pctx  *ProviderContext
	nbmsg int
}

func NewSubmitProvider(ctx *Context, service string) (*MySubmitProvider, error) {
	pctx, err := NewProviderContext(ctx, service)
	if err != nil {
		return nil, err
	}
	submitProvider := &MySubmitProvider{pctx: pctx}
	pctx.RegisterSubmitProvider(200, 1, 1, 1, submitProvider)

	return submitProvider, nil
}

func (provider *MySubmitProvider) OnSubmit(msg *Message, transaction SubmitTransaction) error {
	if msg != nil {
		fmt.Println("\t$$$$$ submitHandler receive: ", string(msg.Body))
		provider.nbmsg += 1
		transaction.Ack(nil)
	} else {
		fmt.Println("receive: nil")
	}
	return nil
}

// Test TCP transport Submit Interaction using the high level API
func TestSubmitProvider(t *testing.T) {
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

	op1, err := consumer.NewSubmitOperation(provider.pctx.Uri, 200, 1, 1, 1)
	if err != nil {
		t.Fatal("Error creating operation, ", err)
		return
	}
	_, err = op1.Submit([]byte("message1"))
	if err != nil {
		t.Fatal("Error during submit, ", err)
		return
	}
	fmt.Println("\t&&&&& Submit1: OK")

	op2, err := consumer.NewSubmitOperation(provider.pctx.Uri, 200, 1, 1, 1)
	if err != nil {
		t.Fatal("Error creating operation, ", err)
		return
	}
	_, err = op2.Submit([]byte("message2"))
	if err != nil {
		t.Fatal("Error during submit, ", err)
		return
	}
	fmt.Println("\t&&&&& Submit2: OK")

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

type MyRequestProvider struct {
	pctx  *ProviderContext
	nbmsg int
}

func NewRequestProvider(ctx *Context, service string) (*MyRequestProvider, error) {
	pctx, err := NewProviderContext(ctx, service)
	if err != nil {
		return nil, err
	}
	requestProvider := &MyRequestProvider{pctx: pctx}
	pctx.RegisterRequestProvider(200, 1, 1, 1, requestProvider)

	return requestProvider, nil
}

func (provider *MyRequestProvider) OnRequest(msg *Message, transaction RequestTransaction) error {
	if msg != nil {
		fmt.Println("\t$$$$$ requestHandler receive: ", string(msg.Body))
		provider.nbmsg += 1
		transaction.Reply([]byte("reply message"), nil)
	} else {
		fmt.Println("receive: nil")
	}
	return nil
}

// Test TCP transport Request Interaction using the high level API
func TestRequestProvider(t *testing.T) {
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

	op1, err := consumer.NewRequestOperation(provider.pctx.Uri, 200, 1, 1, 1)
	if err != nil {
		t.Fatal("Error creating operation, ", err)
		return
	}
	ret1, err := op1.Request([]byte("message1"))
	if err != nil {
		t.Fatal("Error during request, ", err)
		return
	}
	fmt.Println("\t&&&&& Request1: OK, ", string(ret1.Body))

	op2, err := consumer.NewRequestOperation(provider.pctx.Uri, 200, 1, 1, 1)
	if err != nil {
		t.Fatal("Error creating operation, ", err)
		return
	}
	ret2, err := op2.Request([]byte("message2"))
	if err != nil {
		t.Fatal("Error during request, ", err)
		return
	}
	fmt.Println("\t&&&&& Request2: OK, ", string(ret2.Body))

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

type MyInvokeProvider struct {
	pctx  *ProviderContext
	nbmsg int
}

func NewInvokeProvider(ctx *Context, service string) (*MyInvokeProvider, error) {
	pctx, err := NewProviderContext(ctx, service)
	if err != nil {
		return nil, err
	}
	invokeProvider := &MyInvokeProvider{pctx: pctx}
	pctx.RegisterInvokeProvider(200, 1, 1, 1, invokeProvider)

	return invokeProvider, nil
}

func (provider *MyInvokeProvider) OnInvoke(msg *Message, transaction InvokeTransaction) error {
	if msg != nil {
		fmt.Println("\t$$$$$ invokeProvider receive: ", string(msg.Body))
		transaction.Ack(nil)
		provider.nbmsg += 1
		time.Sleep(250 * time.Millisecond)
		//		transaction.Reply([]byte("reply message"), nil)
		transaction.Reply(msg.Body, nil)
	} else {
		fmt.Println("receive: nil")
	}
	return nil
}

// Test TCP transport Invoke Interaction using the high level API
func TestInvokeProvider(t *testing.T) {
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

	op1, err := consumer.NewInvokeOperation(provider.pctx.Uri, 200, 1, 1, 1)
	if err != nil {
		t.Fatal("Error creating operation, ", err)
		return
	}
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

	op2, err := consumer.NewInvokeOperation(provider.pctx.Uri, 200, 1, 1, 1)
	if err != nil {
		t.Fatal("Error creating operation, ", err)
		return
	}
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

type MyProgressProvider struct {
	pctx *ProviderContext
}

func NewProgressProvider(ctx *Context, service string) (*MyProgressProvider, error) {
	pctx, err := NewProviderContext(ctx, service)
	if err != nil {
		return nil, err
	}
	progressProvider := &MyProgressProvider{pctx: pctx}
	pctx.RegisterProgressProvider(200, 1, 1, 1, progressProvider)

	return progressProvider, nil
}

func (provider *MyProgressProvider) OnProgress(msg *Message, transaction ProgressTransaction) error {
	if msg != nil {
		fmt.Println("\t$$$$$ progressHandler receive: ", string(msg.Body))
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
func TestProgressProvider(t *testing.T) {
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

	op1, err := consumer.NewProgressOperation(provider.pctx.Uri, 200, 1, 1, 1)
	op1.Progress([]byte("message1"))
	fmt.Println("\t&&&&& Progress1: OK")

	updt, err := op1.GetUpdate()
	if err != nil {
		t.Error(err)
	}
	for updt != nil {
		nbmsg += 1
		fmt.Println("\t&&&&& Progress1: Update -> ", string(updt.Body))
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
	fmt.Println("\t&&&&& Progress1: Response -> ", string(rep.Body))

	time.Sleep(250 * time.Millisecond)
	provider_ctx.Close()
	consumer_ctx.Close()

	if nbmsg != 11 {
		t.Errorf("Receives %d messages, expect %d ", nbmsg, 2)
	}

	// Waits for socket close
	time.Sleep(250 * time.Millisecond)
}

// ########## ########## ########## ########## ########## ########## ########## ##########

type MyBrokerContext struct {
	pctx *ProviderContext
	subs SubscriberTransaction
}

func NewBroker(ctx *Context, service string) (*MyBrokerContext, error) {
	pctx, err := NewProviderContext(ctx, service)
	if err != nil {
		return nil, err
	}
	broker := &MyBrokerContext{pctx: pctx}
	// Registers the broker handler
	pctx.RegisterBroker(200, 1, 1, 1, broker)

	return broker, nil
}

func (broker *MyBrokerContext) OnRegister(msg *Message, tx SubscriberTransaction) error {
	fmt.Println("\t##########\n\t# OnRegister:")
	broker.subs = tx
	tx.AckRegister(nil)
	return nil
}

func (broker *MyBrokerContext) OnDeregister(msg *Message, tx SubscriberTransaction) error {
	fmt.Println("\t##########\n\t# OnDeregister:")
	broker.subs = nil
	tx.AckDeregister(nil)
	return nil
}

func (broker *MyBrokerContext) OnPublishRegister(msg *Message, tx PublisherTransaction) error {
	fmt.Println("\t##########\n\t# OnPublishRegister:")
	tx.AckRegister(nil)
	return nil
}

func (broker *MyBrokerContext) OnPublishDeregister(msg *Message, tx PublisherTransaction) error {
	fmt.Println("\t##########\n\t# OnPublishDeregister:")
	tx.AckDeregister(nil)
	return nil
}

func (broker *MyBrokerContext) OnPublish(msg *Message, tx PublisherTransaction) error {
	fmt.Println("\t##########\n\t# OnPublish:")
	if broker.subs != nil {
		broker.subs.Notify(msg.Body, nil)
	}
	return nil
}

// Test TCP transport Pub/Sub Interaction using the high level API
func TestPubSubProvider(t *testing.T) {
	broker_ctx, err := NewContext(broker_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}
	broker, err := NewBroker(broker_ctx, "broker")

	// TODO (AF): Creates broker

	pub_ctx, err := NewContext(publisher_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	publisher, err := NewOperationContext(pub_ctx, "publisher")
	if err != nil {
		t.Fatal("Error creating publisher, ", err)
		return
	}
	pubop, err := publisher.NewPublisherOperation(broker.pctx.Uri, 200, 1, 1, 1)
	// TODO (AF): Build PublishRegister message
	pubop.Register([]byte("register"))

	sub_ctx, err := NewContext(subscriber_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	subscriber, err := NewOperationContext(sub_ctx, "subscriber")
	if err != nil {
		t.Fatal("Error creating subscriber, ", err)
		return
	}
	subop, err := subscriber.NewSubscriberOperation(broker.pctx.Uri, 200, 1, 1, 1)
	// TODO (AF): Build Register message
	subop.Register([]byte("register"))

	pubop.Publish([]byte("publish #1"))
	pubop.Publish([]byte("publish #2"))

	// Try to get Notify
	r1, err := subop.GetNotify()
	fmt.Println("\t&&&&& Subscriber notified: OK, ", string(r1.Body))

	// Try to get Notify
	r2, err := subop.GetNotify()
	fmt.Println("\t&&&&& Subscriber notified: OK, ", string(r2.Body))

	pubop.Deregister([]byte("deregister"))

	subop.Deregister([]byte("deregister"))

	// Waits for socket close
	time.Sleep(250 * time.Millisecond)
}
