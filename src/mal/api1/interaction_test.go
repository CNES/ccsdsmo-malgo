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
package api1_test

import (
	"fmt"
	. "github.com/ccsdsmo/malgo/src/mal"
	. "github.com/ccsdsmo/malgo/src/mal/api1"
	_ "github.com/ccsdsmo/malgo/src/mal/transport/tcp" // Needed to initialize TCP transport factory
	"testing"
	"time"
)

const (
	provider_url   = "github.com/ccsdsmo/malgo/src/maltcp://127.0.0.1:16001"
	consumer_url   = "github.com/ccsdsmo/malgo/src/maltcp://127.0.0.1:16002"
	publisher_url  = "github.com/ccsdsmo/malgo/src/maltcp://127.0.0.1:16001"
	subscriber_url = "github.com/ccsdsmo/malgo/src/maltcp://127.0.0.1:16002"
	broker_url     = "github.com/ccsdsmo/malgo/src/maltcp://127.0.0.1:16003"
)

// ########## ########## ########## ########## ########## ########## ########## ##########

// Test TCP transport Send Interaction using the high level API
func TestSend(t *testing.T) {
	provider_ctx, err := NewContext(provider_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	provider, err := NewProviderContext(provider_ctx, "provider")
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

	// Declares Send handler function here so nbmsg local variable is accessible throught
	// the closure.
	sendHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			fmt.Println("\n\n $$$$$ sendHandler receive: ", string(msg.Body), "\n\n")
			nbmsg += 1
		} else {
			fmt.Println("receive: nil")
		}
		return nil
	}
	// Registers Send handler
	provider.RegisterSendHandler(200, 1, 1, 1, sendHandler)

	op1, err := consumer.NewSendOperation(200, 1, 1, 1)
	msg1 := &Message{
		UriTo: provider.Uri,
		Body:  []byte("message1"),
	}
	op1.Send(msg1)

	op2, err := consumer.NewSendOperation(200, 1, 1, 1)
	msg2 := &Message{
		UriTo: provider.Uri,
		Body:  []byte("message2"),
	}
	op2.Send(msg2)

	time.Sleep(250 * time.Millisecond)
	provider_ctx.Close()
	consumer_ctx.Close()

	if nbmsg != 2 {
		t.Errorf("Receives %d messages, expect %d ", nbmsg, 2)
	}

	// Waits for socket close
	time.Sleep(250 * time.Millisecond)
}

// ########## ########## ########## ########## ########## ########## ########## ##########

// Test TCP transport Submit Interaction using the high level API
func TestSubmit(t *testing.T) {
	provider_ctx, err := NewContext(provider_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	provider, err := NewProviderContext(provider_ctx, "provider")
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

	// Declares handler function here so nbmsg local variable is accessible throught
	// the closure.
	submitHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			fmt.Println("\n\n $$$$$ submitHandler receive: ", string(msg.Body), "\n\n")
			nbmsg += 1
			transaction := t.(SubmitTransaction)
			transaction.Ack(nil)
		} else {
			fmt.Println("receive: nil")
		}
		return nil
	}
	// Registers Submit handler
	provider.RegisterSubmitHandler(200, 1, 1, 1, submitHandler)

	op1, err := consumer.NewSubmitOperation(200, 1, 1, 1)
	msg1 := &Message{
		UriTo: provider.Uri,
		Body:  []byte("message1"),
	}
	op1.Submit(msg1)

	fmt.Println("\n\n &&&&& Submit1: OK\n\n")

	op2, err := consumer.NewSubmitOperation(200, 1, 1, 1)
	msg2 := &Message{
		UriTo: provider.Uri,
		Body:  []byte("message2"),
	}
	op2.Submit(msg2)

	fmt.Println("\n\n &&&&& Submit2: OK\n\n")

	time.Sleep(250 * time.Millisecond)
	provider_ctx.Close()
	consumer_ctx.Close()

	if nbmsg != 2 {
		t.Errorf("Receives %d messages, expect %d ", nbmsg, 2)
	}

	// Waits for socket close
	time.Sleep(250 * time.Millisecond)
}

// ########## ########## ########## ########## ########## ########## ########## ##########

// Test TCP transport Request Interaction using the high level API
func TestRequest(t *testing.T) {
	provider_ctx, err := NewContext(provider_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	provider, err := NewProviderContext(provider_ctx, "provider")
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

	// Declares handler function here so nbmsg local variable is accessible throught
	// the closure.

	requestHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			fmt.Println("\n\n $$$$$ requestHandler receive: ", string(msg.Body), "\n\n")
			nbmsg += 1
			transaction := t.(RequestTransaction)
			transaction.Reply([]byte("reply message"), nil)
		} else {
			fmt.Println("receive: nil")
		}
		return nil
	}

	provider.RegisterRequestHandler(200, 1, 1, 1, requestHandler)

	op1, err := consumer.NewRequestOperation(200, 1, 1, 1)
	msg1 := &Message{
		UriTo: provider.Uri,
		Body:  []byte("message1"),
	}
	op1.Request(msg1)

	fmt.Println("\n\n &&&&& Request1: OK\n\n")

	op2, err := consumer.NewRequestOperation(200, 1, 1, 1)
	msg2 := &Message{
		UriTo: provider.Uri,
		Body:  []byte("message2"),
	}
	op2.Request(msg2)

	fmt.Println("\n\n &&&&& Request2: OK\n\n")

	time.Sleep(250 * time.Millisecond)
	provider_ctx.Close()
	consumer_ctx.Close()

	if nbmsg != 2 {
		t.Errorf("Receives %d messages, expect %d ", nbmsg, 2)
	}

	// Waits for socket close
	time.Sleep(250 * time.Millisecond)
}

// ########## ########## ########## ########## ########## ########## ########## ##########

// Test TCP transport Invoke Interaction using the high level API
func TestInvoke(t *testing.T) {
	provider_ctx, err := NewContext(provider_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	provider, err := NewProviderContext(provider_ctx, "provider")
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

	// Declares handler function here so nbmsg local variable is accessible throught
	// the closure.

	invokeHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			fmt.Println("\n\n $$$$$ invokeHandler receive: ", string(msg.Body), "\n\n")
			nbmsg += 1
			transaction := t.(InvokeTransaction)
			transaction.Ack(nil)
			transaction.Reply([]byte("reply message"), nil)
		} else {
			fmt.Println("receive: nil")
		}
		return nil
	}

	provider.RegisterInvokeHandler(200, 1, 1, 1, invokeHandler)

	op1, err := consumer.NewInvokeOperation(200, 1, 1, 1)
	msg1 := &Message{
		UriTo: provider.Uri,
		Body:  []byte("message1"),
	}
	op1.Invoke(msg1)

	r1, err := op1.GetResponse()
	fmt.Println("\n\n &&&&& Invoke1: OK\n", string(r1.Body), "\n\n")

	op2, err := consumer.NewInvokeOperation(200, 1, 1, 1)
	msg2 := &Message{
		UriTo: provider.Uri,
		Body:  []byte("message2"),
	}
	op2.Invoke(msg2)

	r2, err := op1.GetResponse()
	fmt.Println("\n\n &&&&& Invoke2: OK\n", string(r2.Body), "\n\n")

	time.Sleep(250 * time.Millisecond)
	provider_ctx.Close()
	consumer_ctx.Close()

	if nbmsg != 2 {
		t.Errorf("Receives %d messages, expect %d ", nbmsg, 2)
	}

	// Waits for socket close
	time.Sleep(250 * time.Millisecond)
}

// ########## ########## ########## ########## ########## ########## ########## ##########

// Test TCP transport Progress Interaction using the high level API
func TestProgress(t *testing.T) {
	provider_ctx, err := NewContext(provider_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	provider, err := NewProviderContext(provider_ctx, "provider")
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

	// Declares handler function here so nbmsg local variable is accessible throught
	// the closure.
	progressHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			fmt.Println("\n\n $$$$$ progressHandler receive: ", string(msg.Body), "\n\n")
			transaction := t.(ProgressTransaction)
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
	// Registers Progress handler
	provider.RegisterProgressHandler(200, 1, 1, 1, progressHandler)

	op1, err := consumer.NewProgressOperation(200, 1, 1, 1)
	msg1 := &Message{
		UriTo: provider.Uri,
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

// ########## ########## ########## ########## ########## ########## ########## ##########

// Test TCP transport Pub/Sub Interaction using the high level API
func TestPubSub(t *testing.T) {
	pub_ctx, err := NewContext(publisher_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	publisher, err := NewOperationContext(pub_ctx, "publisher")
	if err != nil {
		t.Fatal("Error creating provider, ", err)
		return
	}

	sub_ctx, err := NewContext(subscriber_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	subscriber, err := NewOperationContext(sub_ctx, "subscriber")
	if err != nil {
		t.Fatal("Error creating consumer, ", err)
		return
	}

	broker_ctx, err := NewContext(broker_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	broker, err := NewProviderContext(broker_ctx, "broker")
	if err != nil {
		t.Fatal("Error creating consumer, ", err)
		return
	}

	//	nbmsg := 0

	// Declares handler functions here so nbmsg local variable is accessible throught
	// the closure.

	// Be careful: handling of a unique global subscription !!
	var subs SubscriberTransaction = nil
	brokerHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH_REGISTER {
				transaction := t.(PublisherTransaction)
				fmt.Println("\n\n $$$$$ publisherHandler receive: PUBLISH_REGISTER", string(msg.Body), "==> ", msg.TransactionId, "\n", "\n\n")
				transaction.AckRegister(nil)
			} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH {
				//				transaction := t.(PublisherTransaction)
				fmt.Println("\n\n $$$$$ publisherHandler receive: PUBLISH", string(msg.Body), "==> ", msg.TransactionId, "\n", "\n\n")
				// TODO (AF): We should verify that the publisher is registered
				if subs != nil {
					subs.Notify(msg.Body, nil)
				}
			} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH_DEREGISTER {
				transaction := t.(PublisherTransaction)
				fmt.Println("\n\n $$$$$ publisherHandler receive: PUBLISH_DEREGISTER", string(msg.Body), "==> ", msg.TransactionId, "\n", "\n\n")
				transaction.AckDeregister(nil)
			} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_REGISTER {
				transaction := t.(SubscriberTransaction)
				fmt.Println("\n\n $$$$$ subscriberHandler receive: REGISTER", string(msg.Body), "==> ", msg.TransactionId, "\n", "\n\n")
				subs = transaction
				transaction.AckRegister(nil)
			} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_DEREGISTER {
				transaction := t.(SubscriberTransaction)
				fmt.Println("\n\n $$$$$ subscriberHandler receive: DEREGISTER", string(msg.Body), "==> ", msg.TransactionId, "\n", "\n\n")
				subs = nil
				transaction.AckDeregister(nil)
			} else {
				fmt.Println("\n\n $$$$$ publisherHandler receive: Bad message ", msg.InteractionStage, " -> ", "==> ", msg.TransactionId, string(msg.Body), "\n\n")
			}
		} else {
			fmt.Println("receive: nil")
		}
		return nil
	}
	err = broker.RegisterBrokerHandler(200, 1, 1, 1, brokerHandler)
	if err != nil {
		t.Fatal("Error registering publisher handler, ", err)
		return
	}

	// Initiates Publisher operation and do register
	op1, err := publisher.NewPublisherOperation(200, 1, 1, 1)
	if err != nil {
		t.Fatal("Error creating publisher operation, ", err)
		return
	}
	msg1 := &Message{
		UriTo: broker.Uri,
		Body:  []byte("register"),
	}
	ack, err := op1.Register(msg1)
	if err != nil {
		t.Fatal("Error during publish register operation, ", err)
		return
	}
	fmt.Println("\n\n &&&&& Publisher registered: ", string(ack.Body), "\n\n")

	// Initiates Subscriber operation and do register
	op2, err := subscriber.NewSubscriberOperation(200, 1, 1, 1)
	if err != nil {
		t.Fatal("Error creating subscriber operation, ", err)
		return
	}
	msg2 := &Message{
		UriTo: broker.Uri,
		Body:  []byte("register"),
	}
	ack, err = op2.Register(msg2)
	if err != nil {
		t.Fatal("Error during register operation, ", err)
		return
	}
	fmt.Println("\n\n &&&&& Subscriber registered: ", string(ack.Body), "\n\n")

	// Do publish
	msg3 := &Message{
		UriTo: broker.Uri,
		Body:  []byte("publish"),
	}
	err = op1.Publish(msg3)
	if err != nil {
		t.Fatal("Error during publish operation, ", err)
		return
	}
	fmt.Println("\n\n &&&&& Publisher publish: OK\n\n")

	// Try to get Notify
	r, err := op2.GetNotify()
	fmt.Println("\n\n &&&&& Subscriber notified: \n", string(r.Body), "OK\n\n")

	// Do Deregister
	msg4 := &Message{
		UriTo: broker.Uri,
		Body:  []byte("deregister"),
	}
	ack, err = op1.Deregister(msg4)
	if err != nil {
		t.Fatal("Error during publish deregister operation, ", err)
		return
	}
	fmt.Println("\n\n &&&&& Publisher deregister: ", string(ack.Body), "\n\n")

	// Do Deregister
	msg5 := &Message{
		UriTo: broker.Uri,
		Body:  []byte("deregister"),
	}
	ack, err = op2.Deregister(msg5)
	if err != nil {
		t.Fatal("Error during deregister operation, ", err)
		return
	}
	fmt.Println("\n\n &&&&& Subscriber deregister: ", string(ack.Body), "\n\n")

	time.Sleep(250 * time.Millisecond)
	pub_ctx.Close()
	sub_ctx.Close()
	broker_ctx.Close()

	//	if nbmsg != 11 {
	//		t.Errorf("Receives %d messages, expect %d ", nbmsg, 2)
	//	}

	// Waits for socket close
	time.Sleep(250 * time.Millisecond)
}
