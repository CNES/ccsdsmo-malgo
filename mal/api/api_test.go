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
package api_test

import (
	"errors"
	"fmt"
	. "github.com/CNES/ccsdsmo-malgo/mal"
	. "github.com/CNES/ccsdsmo-malgo/mal/api"
	_ "github.com/CNES/ccsdsmo-malgo/mal/transport/tcp" // Needed to initialize TCP transport factory
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
			par, err := msg.DecodeLastParameter(NullString, false)
			fmt.Println("\t$$$$$ sendHandler receive: ", *par.(*String), err)
			provider.nbmsg += 1
		} else {
			fmt.Println("\tERROR sendHandler receive: nil")
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
	body := op1.NewBody()
	body.EncodeLastParameter(NewString("message1"), false)
	op1.Send(body)

	op2 := consumer.NewSendOperation(provider.cctx.Uri, 200, 1, 1, 1)
	body = op2.NewBody()
	body.EncodeLastParameter(NewString("message2"), false)
	op2.Send(body)

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

	body := consumer_ctx.NewBody()
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

// ########## ########## ########## ########## ########## ########## ########## ##########
// Test Request interaction

// Define RequestProvider

type TestRequestProvider struct {
	ctx   *Context
	cctx  *ClientContext
	nbmsg int
}

func newTestRequestProvider() (*TestRequestProvider, error) {
	ctx, err := NewContext(provider_url)
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
			par, err := msg.DecodeLastParameter(NullString, false)
			fmt.Println("\t$$$$$ requestHandler receive: ", *par.(*String), err)
			provider.nbmsg += 1
			body := t.NewBody()
			body.EncodeLastParameter(NewString("reply message"), false)
			transaction.Reply(body, false)
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
	body := op1.NewBody()
	body.EncodeLastParameter(NewString("message1"), false)
	ret1, err := op1.Request(body)
	if err != nil {
		t.Fatal("Error during request, ", err)
		return
	}
	par1, err := ret1.DecodeLastParameter(NullString, false)
	fmt.Println("\t&&&&& Request1: OK, ", *par1.(*String))

	op2 := consumer.NewRequestOperation(provider.cctx.Uri, 200, 1, 1, 1)
	body = op2.NewBody()
	body.EncodeLastParameter(NewString("message2"), false)
	ret2, err := op2.Request(body)
	if err != nil {
		t.Fatal("Error during request, ", err)
		return
	}
	par2, err := ret2.DecodeLastParameter(NullString, false)
	fmt.Println("\t&&&&& Request2: OK, ", *par2.(*String))

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
	ctx, err := NewContext(provider_url)
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
			par, err := msg.DecodeLastParameter(NullString, false)
			fmt.Println("\t$$$$$ invokeProvider receive: ", *par.(*String), err)
			transaction.Ack(nil, false)
			provider.nbmsg += 1
			time.Sleep(250 * time.Millisecond)

			// Note (AF): Be careful, the body of a previously decoded message should be used in a
			// newly message to send (the encoder is nil).

			// body := ctx.NewBody()
			//body.EncodeLastParameter(NewString("reply message"), false)

			transaction.Reply(msg.Body, false)
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

	body := consumer_ctx.NewBody()
	body.EncodeLastParameter(NewString("message1"), false)
	op1 := consumer.NewInvokeOperation(provider.cctx.Uri, 200, 1, 1, 1)
	_, err = op1.Invoke(body)
	if err != nil {
		t.Fatal("Error during invoke, ", err)
		return
	}

	r1, err := op1.GetResponse()
	if err != nil {
		t.Fatal("Error getting response, ", err)
		return
	}
	p1, err := r1.DecodeLastParameter(NullString, false)
	fmt.Println("\t&&&&& Invoke1: OK, ", *p1.(*String))

	body = consumer_ctx.NewBody()
	body.EncodeLastParameter(NewString("message2"), false)
	op2 := consumer.NewInvokeOperation(provider.cctx.Uri, 200, 1, 1, 1)
	_, err = op2.Invoke(body)
	if err != nil {
		t.Fatal("Error during invoke, ", err)
		return
	}

	r2, err := op2.GetResponse()
	if err != nil {
		t.Fatal("Error getting response, ", err)
		return
	}
	p2, err := r2.DecodeLastParameter(NullString, false)
	fmt.Println("\t&&&&& Invoke2: OK, ", *p2.(*String))

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
			par, err := msg.DecodeLastParameter(NullString, false)
			fmt.Println("\t$$$$$ progressHandler1 receive: ", *par.(*String), err)
			transaction := t.(ProgressTransaction)
			transaction.Ack(nil, false)
			for i := 0; i < 10; i++ {
				// Note (AF): Be careful, do not reuse a previously sent Body before you are
				// sure that it is no longer used.
				body := t.NewBody()
				body.EncodeLastParameter(NewString(fmt.Sprintf("messsage1.#%d", i)), false)
				transaction.Update(body, false)
			}
			body := t.NewBody()
			body.EncodeLastParameter(NewString("last message"), false)
			transaction.Reply(body, false)
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
			par, err := msg.DecodeLastParameter(NullString, false)
			fmt.Println("\t$$$$$ progressHandler2 receive: ", *par.(*String), err)
			transaction := t.(ProgressTransaction)
			transaction.Ack(nil, false)
			for i := 0; i < 5; i++ {
				// Note (AF): Be careful, do not reuse a previously sent Body before you are
				// sure that it is no longer used.
				body := t.NewBody()
				body.EncodeLastParameter(NewString(fmt.Sprintf("messsage2.#%d", i)), false)
				transaction.Update(body, false)
			}
			body := t.NewBody()
			body.EncodeLastParameter(NewString("last message2"), false)
			transaction.Reply(body, false)
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
	body := op1.NewBody()
	body.EncodeLastParameter(NewString("message1"), false)
	op1.Progress(body)
	fmt.Println("\t&&&&& Progress1: OK")

	updt, err := op1.GetUpdate()
	if err != nil {
		t.Error(err)
	}
	for updt != nil {
		nbmsg += 1
		p, err := updt.DecodeLastParameter(NullString, false)
		fmt.Println("\t&&&&& Progress1: Update -> ", *p.(*String))
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
	p, err := rep.DecodeLastParameter(NullString, false)
	fmt.Println("\t&&&&& Progress1: Response -> ", *p.(*String))

	if nbmsg != 11 {
		t.Errorf("Receives %d messages, expect %d ", nbmsg, 2)
	}

	// Call provider.Op2
	op2 := consumer.NewProgressOperation(provider.uri, 200, 1, 1, 2)
	body = op2.NewBody()
	body.EncodeLastParameter(NewString("message2"), false)
	op2.Progress(body)
	fmt.Println("\t&&&&& Progress2: OK")

	updt, err = op2.GetUpdate()
	if err != nil {
		t.Error(err)
	}
	for updt != nil {
		nbmsg += 1
		p, err := updt.DecodeLastParameter(NullString, false)
		fmt.Println("\t&&&&& Progress2: Update -> ", *p.(*String))
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
	p, err = rep.DecodeLastParameter(NullString, false)
	fmt.Println("\t&&&&& Progress2: Response -> ", *p.(*String))

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
	par, err := msg.DecodeLastParameter(NullString, false)
	fmt.Println("\t##########\n\t# OnRegister: ", *par.(*String), err)
	broker.subs = tx
	tx.AckRegister(nil, false)
	return nil
}

func (broker *TestPubSubProvider) OnDeregister(msg *Message, tx SubscriberTransaction) error {
	par, err := msg.DecodeLastParameter(NullString, false)
	fmt.Println("\t##########\n\t# OnDeregister: ", *par.(*String), err)
	broker.subs = nil
	tx.AckDeregister(nil, false)
	return nil
}

func (broker *TestPubSubProvider) OnPublishRegister(msg *Message, tx PublisherTransaction) error {
	par, err := msg.DecodeLastParameter(NullString, false)
	fmt.Println("\t##########\n\t# OnPublishRegister: ", *par.(*String), err)
	tx.AckRegister(nil, false)
	return nil
}

func (broker *TestPubSubProvider) OnPublish(msg *Message, tx PublisherTransaction) error {
	par, err := msg.DecodeLastParameter(NullString, false)
	fmt.Println("\t##########\n\t# OnPublish: ", *par.(*String), err)
	if broker.subs != nil {
		fmt.Println("\t##########\n\t# OnPublish: Notify")
		//		body := tx.NewBody()
		//		body.EncodeLastParameter(NewString("Notify message"), false)

		// Note (AF): Be careful, the body of a previously decoded message should not be used in a
		// newly message to send (the encoder is nil).

		broker.subs.Notify(msg.Body, false)
	}
	return nil
}

func (broker *TestPubSubProvider) OnPublishDeregister(msg *Message, tx PublisherTransaction) error {
	par, err := msg.DecodeLastParameter(NullString, false)
	fmt.Println("\t##########\n\t# OnPublishDeregister: ", *par.(*String), err)
	tx.AckDeregister(nil, false)
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
	body := pubop.NewBody()
	body.EncodeLastParameter(NewString("register#1"), false)
	// TODO (AF): Build PublishRegister message
	pubop.Register(body)

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

	body.Reset(true)
	body.EncodeLastParameter(NewString("register#2"), false)
	subop := subscriber.NewSubscriberOperation(broker.cctx.Uri, 200, 1, 1, 1)
	// TODO (AF): Build Register message
	subop.Register(body)

	body = pubop.NewBody()
	body.EncodeLastParameter(NewString("publish#1"), false)
	pubop.Publish(body)
	// Note (AF): Be careful, do not reuse a previously sent Body before you are
	// sure that it is no longer used (No acknowledge on publish)
	body = pubop.NewBody()
	body.EncodeLastParameter(NewString("publish#2"), false)
	pubop.Publish(body)

	// Try to get Notify
	r1, err := subop.GetNotify()
	p, err := r1.DecodeLastParameter(NullString, false)
	fmt.Println("\t&&&&& Subscriber notified: OK, ", *p.(*String))

	// Try to get Notify
	r2, err := subop.GetNotify()
	p, err = r2.DecodeLastParameter(NullString, false)
	fmt.Println("\t&&&&& Subscriber notified: OK, ", *p.(*String))

	body.Reset(true)
	body.EncodeLastParameter(NewString("deregister#1"), false)
	pubop.Deregister(body)
	body.Reset(true)
	body.EncodeLastParameter(NewString("deregister#2"), false)
	subop.Deregister(body)

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

	body := consumer_ctx.NewBody()
	body.EncodeLastParameter(NewString("message1"), false)
	op1 := consumer.NewSubmitOperation(provider.cctx.Uri, 200, 1, 1, 1)
	_, err = op1.Submit(body)
	if err != nil {
		t.Fatal("Error during submit, ", err)
		return
	}
	fmt.Println("\t&&&&& Submit1: OK")
	op1.Reset()
	body = consumer_ctx.NewBody()
	body.EncodeLastParameter(NewString("message2"), false)
	_, err = op1.Submit(body)
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
