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
	"errors"
	"fmt"
	. "mal"
	. "mal/api"
	_ "mal/transport/tcp" // Needed to initialize TCP transport factory
	"testing"
	"time"
)

const (
	provider_url   = "maltcp://127.0.0.1:16001"
	consumer_url   = "maltcp://127.0.0.1:16002"
	subscriber_url = "maltcp://127.0.0.1:16001"
	publisher_url  = "maltcp://127.0.0.1:16002"
	broker_url     = "maltcp://127.0.0.1:16003"
)

// ########## ########## ########## ########## ########## ########## ########## ##########
// Test Send interaction

// Define SendProvider

type TestSendProvider struct {
	ctx   *Context
	cctx  *ClientContext
	nbmsg int
}

func newTestSendProvider() (*TestSendProvider, error) {
	ctx, err := NewContext(provider_url)
	if err != nil {
		return nil, err
	}
	cctx, err := NewClientContext(ctx, "provider")
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
	provider.ctx.Close()
}

// Test TCP transport Send Interaction using the high level API
func TestSend(t *testing.T) {
	// Waits socket closing from previous test
	time.Sleep(250 * time.Millisecond)

	provider, err := newTestSendProvider()
	if err != nil {
		t.Fatal("Error creating provider, ", err)
		return
	}
	defer provider.close()

	consumer_ctx, err := NewContext(consumer_url)
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

	op1 := consumer.NewSendOperation(provider.cctx.Uri, 200, 1, 1, 1)
	op1.Send([]byte("message1"))

	op2 := consumer.NewSendOperation(provider.cctx.Uri, 200, 1, 1, 1)
	op2.Send([]byte("message2"))

	// Waits for message reception
	time.Sleep(250 * time.Millisecond)

	if provider.nbmsg != 2 {
		t.Errorf("Receives %d messages, expect %d ", provider.nbmsg, 2)
	}
}

// ########## ########## ########## ########## ########## ########## ########## ##########
// Test Submit interaction

// Define SubmitProvider

type TestSubmitProvider struct {
	ctx   *Context
	cctx  *ClientContext
	nbmsg int
}

func newTestSubmitProvider() (*TestSubmitProvider, error) {
	ctx, err := NewContext(provider_url)
	if err != nil {
		return nil, err
	}
	cctx, err := NewClientContext(ctx, "provider")
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
			transaction.Ack(nil)
		} else {
			fmt.Println("receive: nil")
		}
		return nil
	}
	cctx.RegisterSubmitHandler(200, 1, 1, 1, submitHandler)

	return provider, nil
}

func (provider *TestSubmitProvider) close() {
	provider.ctx.Close()
}

// Test TCP transport Submit Interaction using the high level API
func TestSubmit(t *testing.T) {
	// Waits socket closing from previous test
	time.Sleep(250 * time.Millisecond)

	provider, err := newTestSubmitProvider()
	if err != nil {
		t.Fatal("Error creating provider, ", err)
		return
	}
	defer provider.close()

	consumer_ctx, err := NewContext(consumer_url)
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

	op1 := consumer.NewSubmitOperation(provider.cctx.Uri, 200, 1, 1, 1)
	_, err = op1.Submit([]byte("message1"))
	if err != nil {
		t.Fatal("Error during submit, ", err)
		return
	}
	fmt.Println("\t&&&&& Submit1: OK")

	op2 := consumer.NewSubmitOperation(provider.cctx.Uri, 200, 1, 1, 1)
	_, err = op2.Submit([]byte("message2"))
	if err != nil {
		t.Fatal("Error during submit, ", err)
		return
	}
	fmt.Println("\t&&&&& Submit2: OK")

	if provider.nbmsg != 2 {
		t.Errorf("Receives %d messages, expect %d ", provider.nbmsg, 2)
	}
}

// ########## ########## ########## ########## ########## ########## ########## ##########
// Test Request interaction

// Define RequestProvider

type TestRequestProvider struct {
	ctx   *Context
	cctx  *ClientContext
	nbmsg int
}

func newTestRequestProvider() (*TestRequestProvider, error) {
	ctx, err := NewContext(nested1_provider2_url)
	if err != nil {
		return nil, err
	}
	cctx, err := NewClientContext(ctx, "provider")
	if err != nil {
		return nil, err
	}
	provider := &TestRequestProvider{ctx, cctx, 0}

	// Register handler
	requestHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			transaction := t.(RequestTransaction)
			fmt.Println("\t$$$$$ requestHandler receive: ", string(msg.Body))
			provider.nbmsg += 1
			transaction.Reply([]byte("reply message"), nil)
		} else {
			fmt.Println("receive: nil")
		}
		return nil
	}
	cctx.RegisterRequestHandler(200, 1, 1, 1, requestHandler)

	return provider, nil
}

func (provider *TestRequestProvider) close() {
	provider.ctx.Close()
}

// Test TCP transport Request Interaction using the high level API
func TestRequest(t *testing.T) {
	// Waits socket closing from previous test
	time.Sleep(250 * time.Millisecond)

	provider, err := newTestRequestProvider()
	if err != nil {
		t.Fatal("Error creating provider, ", err)
		return
	}
	defer provider.close()

	consumer_ctx, err := NewContext(consumer_url)
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

	op1 := consumer.NewRequestOperation(provider.cctx.Uri, 200, 1, 1, 1)
	ret1, err := op1.Request([]byte("message1"))
	if err != nil {
		t.Fatal("Error during request, ", err)
		return
	}
	fmt.Println("\t&&&&& Request1: OK, ", string(ret1.Body))

	op2 := consumer.NewRequestOperation(provider.cctx.Uri, 200, 1, 1, 1)
	ret2, err := op2.Request([]byte("message2"))
	if err != nil {
		t.Fatal("Error during request, ", err)
		return
	}
	fmt.Println("\t&&&&& Request2: OK, ", string(ret2.Body))

	if provider.nbmsg != 2 {
		t.Errorf("Receives %d messages, expect %d ", provider.nbmsg, 2)
	}
}

// ########## ########## ########## ########## ########## ########## ########## ##########
// Test Invoke interaction

// Define InvokeProvider

type TestInvokeProvider struct {
	ctx   *Context
	cctx  *ClientContext
	nbmsg int
}

func newTestInvokeProvider() (*TestInvokeProvider, error) {
	ctx, err := NewContext(nested1_provider2_url)
	if err != nil {
		return nil, err
	}
	cctx, err := NewClientContext(ctx, "provider")
	if err != nil {
		return nil, err
	}
	provider := &TestInvokeProvider{ctx, cctx, 0}

	// Register handler
	invokeHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			transaction := t.(InvokeTransaction)
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
	cctx.RegisterInvokeHandler(200, 1, 1, 1, invokeHandler)

	return provider, nil
}

func (provider *TestInvokeProvider) close() {
	provider.ctx.Close()
}

// Test TCP transport Invoke Interaction using the high level API
func TestInvoke(t *testing.T) {
	// Waits socket closing from previous test
	time.Sleep(250 * time.Millisecond)

	provider, err := newTestInvokeProvider()
	if err != nil {
		t.Fatal("Error creating provider, ", err)
		return
	}
	defer provider.close()

	consumer_ctx, err := NewContext(consumer_url)
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

	op1 := consumer.NewInvokeOperation(provider.cctx.Uri, 200, 1, 1, 1)
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

	op2 := consumer.NewInvokeOperation(provider.cctx.Uri, 200, 1, 1, 1)
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

	if provider.nbmsg != 2 {
		t.Errorf("Receives %d messages, expect %d ", provider.nbmsg, 2)
	}
}

// ########## ########## ########## ########## ########## ########## ########## ##########
// Test Progress interaction

// Define a Provider with 2 progress interaction

type TestProgressProvider struct {
	ctx   *Context
	cctx  *ClientContext
	uri   *URI
	nbmsg int
}

func newTestProgressProvider() (*TestProgressProvider, error) {
	ctx, err := NewContext(provider_url)
	if err != nil {
		return nil, err
	}

	cctx, err := NewClientContext(ctx, "provider")
	if err != nil {
		return nil, err
	}

	provider := &TestProgressProvider{ctx, cctx, cctx.Uri, 0}

	// Handler1
	progressHandler1 := func(msg *Message, t Transaction) error {
		provider.nbmsg += 1
		if msg != nil {
			fmt.Println("\t$$$$$ progressHandler1 receive: ", string(msg.Body))
			transaction := t.(ProgressTransaction)
			transaction.Ack(nil)
			for i := 0; i < 10; i++ {
				transaction.Update([]byte(fmt.Sprintf("messsage1.#%d", i)), nil)
			}
			transaction.Reply([]byte("last message1"), nil)
		} else {
			fmt.Println("receive: nil")
		}
		return nil
	}
	// Registers Progress handler
	cctx.RegisterProgressHandler(200, 1, 1, 1, progressHandler1)

	// Handler2
	progressHandler2 := func(msg *Message, t Transaction) error {
		provider.nbmsg += 1
		if msg != nil {
			fmt.Println("\t$$$$$ progressHandler2 receive: ", string(msg.Body))
			transaction := t.(ProgressTransaction)
			transaction.Ack(nil)
			for i := 0; i < 5; i++ {
				transaction.Update([]byte(fmt.Sprintf("messsage2.#%d", i)), nil)
			}
			transaction.Reply([]byte("last message2"), nil)
		} else {
			fmt.Println("receive: nil")
		}
		return nil
	}
	// Registers Progress handler
	cctx.RegisterProgressHandler(200, 1, 1, 2, progressHandler2)

	return provider, nil
}

func (provider *TestProgressProvider) close() {
	provider.ctx.Close()
}

// Test TCP transport Progress Interaction using the high level API
func TestProgress(t *testing.T) {
	// Waits socket closing from previous test
	time.Sleep(250 * time.Millisecond)

	provider, err := newTestProgressProvider()
	if err != nil {
		t.Fatal("Error creating provider, ", err)
		return
	}
	defer provider.close()

	consumer_ctx, err := NewContext(consumer_url)
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

	nbmsg := 0

	op1 := consumer.NewProgressOperation(provider.cctx.Uri, 200, 1, 1, 1)
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

	if nbmsg != 11 {
		t.Errorf("Receives %d messages, expect %d ", nbmsg, 2)
	}

	// Call provider.Op2
	op2 := consumer.NewProgressOperation(provider.uri, 200, 1, 1, 2)
	op2.Progress([]byte("message2"))
	fmt.Println("\t&&&&& Progress2: OK")

	updt, err = op2.GetUpdate()
	if err != nil {
		t.Error(err)
	}
	for updt != nil {
		nbmsg += 1
		fmt.Println("\t&&&&& Progress2: Update -> ", string(updt.Body))
		updt, err = op2.GetUpdate()
		if err != nil {
			t.Error(err)
		}
	}
	rep, err = op2.GetResponse()
	if err != nil {
		t.Error(err)
	}
	nbmsg += 1
	fmt.Println("\t&&&&& Progress2: Response -> ", string(rep.Body))

	if nbmsg != 17 {
		t.Errorf("Receives %d messages, expect %d ", nbmsg, 17)
	}

	if provider.nbmsg != 2 {
		t.Errorf("Provider receives %d messages, expect %d ", provider.nbmsg, 2)
	}
}

// ########## ########## ########## ########## ########## ########## ########## ##########

type TestPubSubProvider struct {
	ctx  *Context
	cctx *ClientContext
	subs SubscriberTransaction
}

func newTestPubSubProvider() (*TestPubSubProvider, error) {
	ctx, err := NewContext(broker_url)
	if err != nil {
		return nil, err
	}
	cctx, err := NewClientContext(ctx, "broker")
	if err != nil {
		return nil, err
	}
	broker := &TestPubSubProvider{ctx, cctx, nil}
	// Register handler
	brokerHandler := func(msg *Message, t Transaction) error {
		if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH_REGISTER {
			broker.OnPublishRegister(msg, t.(PublisherTransaction))
		} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH {
			broker.OnPublish(msg, t.(PublisherTransaction))
		} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH_DEREGISTER {
			broker.OnPublishDeregister(msg, t.(PublisherTransaction))
		} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_REGISTER {
			broker.OnRegister(msg, t.(SubscriberTransaction))
		} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_DEREGISTER {
			broker.OnDeregister(msg, t.(SubscriberTransaction))
		} else {
			return errors.New("Bad stage")
		}
		return nil
	}
	// Registers the broker handler
	cctx.RegisterBrokerHandler(200, 1, 1, 1, brokerHandler)

	return broker, nil
}

func (broker *TestPubSubProvider) close() {
	broker.ctx.Close()
}

func (broker *TestPubSubProvider) OnRegister(msg *Message, tx SubscriberTransaction) error {
	fmt.Println("\t##########\n\t# OnRegister:")
	broker.subs = tx
	tx.AckRegister(nil)
	return nil
}

func (broker *TestPubSubProvider) OnDeregister(msg *Message, tx SubscriberTransaction) error {
	fmt.Println("\t##########\n\t# OnDeregister:")
	broker.subs = nil
	tx.AckDeregister(nil)
	return nil
}

func (broker *TestPubSubProvider) OnPublishRegister(msg *Message, tx PublisherTransaction) error {
	fmt.Println("\t##########\n\t# OnPublishRegister:")
	tx.AckRegister(nil)
	return nil
}

func (broker *TestPubSubProvider) OnPublish(msg *Message, tx PublisherTransaction) error {
	fmt.Println("\t##########\n\t# OnPublish:")
	if broker.subs != nil {
		broker.subs.Notify(msg.Body, nil)
	}
	return nil
}

func (broker *TestPubSubProvider) OnPublishDeregister(msg *Message, tx PublisherTransaction) error {
	fmt.Println("\t##########\n\t# OnPublishDeregister:")
	tx.AckDeregister(nil)
	return nil
}

// Test TCP transport Pub/Sub Interaction using the high level API
func TestPubSub(t *testing.T) {
	// Waits socket closing from previous test
	time.Sleep(250 * time.Millisecond)

	broker, err := newTestPubSubProvider()
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}
	defer broker.close()

	// TODO (AF): Creates broker

	pub_ctx, err := NewContext(publisher_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}
	defer pub_ctx.Close()

	publisher, err := NewClientContext(pub_ctx, "publisher")
	if err != nil {
		t.Fatal("Error creating publisher, ", err)
		return
	}
	pubop := publisher.NewPublisherOperation(broker.cctx.Uri, 200, 1, 1, 1)
	// TODO (AF): Build PublishRegister message
	pubop.Register([]byte("register"))

	sub_ctx, err := NewContext(subscriber_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}
	defer sub_ctx.Close()

	subscriber, err := NewClientContext(sub_ctx, "subscriber")
	if err != nil {
		t.Fatal("Error creating subscriber, ", err)
		return
	}
	subop := subscriber.NewSubscriberOperation(broker.cctx.Uri, 200, 1, 1, 1)
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

// ########## ########## ########## ########## ########## ########## ########## ##########
// Test reuse of operation after reset

func TestReset(t *testing.T) {
	// Waits socket closing from previous test
	time.Sleep(250 * time.Millisecond)

	provider, err := newTestSubmitProvider()
	if err != nil {
		t.Fatal("Error creating provider, ", err)
		return
	}
	defer provider.close()

	consumer_ctx, err := NewContext(consumer_url)
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

	op1 := consumer.NewSubmitOperation(provider.cctx.Uri, 200, 1, 1, 1)
	_, err = op1.Submit([]byte("message1"))
	if err != nil {
		t.Fatal("Error during submit, ", err)
		return
	}
	fmt.Println("\t&&&&& Submit1: OK")
	op1.Reset()
	_, err = op1.Submit([]byte("message2"))
	if err != nil {
		t.Fatal("Error during submit, ", err)
		return
	}
	fmt.Println("\t&&&&& Submit2: OK")

	// Waits for message reception
	time.Sleep(250 * time.Millisecond)

	if provider.nbmsg != 2 {
		t.Errorf("Receives %d messages, expect %d ", provider.nbmsg, 2)
	}
}
