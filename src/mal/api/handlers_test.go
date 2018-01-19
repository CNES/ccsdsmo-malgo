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

const (
	provider_url   = "maltcp://127.0.0.1:16001"
	consumer_url   = "maltcp://127.0.0.1:16002"
	publisher_url  = "maltcp://127.0.0.1:16001"
	subscriber_url = "maltcp://127.0.0.1:16002"
	broker_url     = "maltcp://127.0.0.1:16003"
)

// ########## ########## ########## ########## ########## ########## ########## ##########

// Test TCP transport Send Interaction using the high level API
func TestSend(t *testing.T) {
	provider_ctx, err := NewContext(provider_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	provider, err := NewHandlerContext(provider_ctx, "provider")
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
			fmt.Println("\t$$$$$ sendHandler receive: ", string(msg.Body))
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

	provider, err := NewHandlerContext(provider_ctx, "provider")
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
			fmt.Println("\t$$$$$ submitHandler receive: ", string(msg.Body))
			nbmsg += 1
			t.(SubmitTransaction).Ack(nil)
		} else {
			fmt.Println("receive: nil")
		}
		return nil
	}
	// Registers Submit handler
	provider.RegisterSubmitHandler(200, 1, 1, 1, submitHandler)

	op1, err := consumer.NewSubmitOperation(200, 1, 1, 1)
	if err != nil {
		t.Fatal("Error creating operation, ", err)
		return
	}
	msg1 := &Message{
		UriTo: provider.Uri,
		Body:  []byte("message1"),
	}
	_, err = op1.Submit(msg1)
	if err != nil {
		t.Fatal("Error during submit, ", err)
		return
	}
	fmt.Println("\t&&&&& Submit1: OK")

	op2, err := consumer.NewSubmitOperation(200, 1, 1, 1)
	if err != nil {
		t.Fatal("Error creating operation, ", err)
		return
	}
	msg2 := &Message{
		UriTo: provider.Uri,
		Body:  []byte("message2"),
	}
	_, err = op2.Submit(msg2)
	if err != nil {
		t.Fatal("Error during submit, ", err)
		return
	}
	fmt.Println("\t&&&&& Submit2: OK")

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

	provider, err := NewHandlerContext(provider_ctx, "provider")
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
			fmt.Println("\t$$$$$ requestHandler receive: ", string(msg.Body))
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
	if err != nil {
		t.Fatal("Error creating operation, ", err)
		return
	}
	msg1 := &Message{
		UriTo: provider.Uri,
		Body:  []byte("message1"),
	}
	ret1, err := op1.Request(msg1)
	if err != nil {
		t.Fatal("Error during request, ", err)
		return
	}
	fmt.Println("\t&&&&& Request1: OK, ", string(ret1.Body))

	op2, err := consumer.NewRequestOperation(200, 1, 1, 1)
	if err != nil {
		t.Fatal("Error creating operation, ", err)
		return
	}
	msg2 := &Message{
		UriTo: provider.Uri,
		Body:  []byte("message2"),
	}
	ret2, err := op2.Request(msg2)
	if err != nil {
		t.Fatal("Error during request, ", err)
		return
	}
	fmt.Println("\t&&&&& Request2: OK, ", string(ret2.Body))

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

	provider, err := NewHandlerContext(provider_ctx, "provider")
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
			fmt.Println("\t$$$$$ invokeHandler receive: ", string(msg.Body))
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
	if err != nil {
		t.Fatal("Error creating operation, ", err)
		return
	}
	msg1 := &Message{
		UriTo: provider.Uri,
		Body:  []byte("message1"),
	}
	_, err = op1.Invoke(msg1)
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

	op2, err := consumer.NewInvokeOperation(200, 1, 1, 1)
	if err != nil {
		t.Fatal("Error creating operation, ", err)
		return
	}
	msg2 := &Message{
		UriTo: provider.Uri,
		Body:  []byte("message2"),
	}
	_, err = op2.Invoke(msg2)
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

	provider, err := NewHandlerContext(provider_ctx, "provider")
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
			fmt.Println("\t$$$$$ progressHandler receive: ", string(msg.Body))
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

	broker, err := NewHandlerContext(broker_ctx, "broker")
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
				fmt.Println("\t$$$$$ publisherHandler receive: PUBLISH_REGISTER", string(msg.Body), "==> ", msg.TransactionId)
				transaction.AckRegister(nil)
			} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH {
				//				transaction := t.(PublisherTransaction)
				fmt.Println("\t$$$$$ publisherHandler receive: PUBLISH", string(msg.Body), "==> ", msg.TransactionId)
				// TODO (AF): We should verify that the publisher is registered
				if subs != nil {
					subs.Notify(msg.Body, nil)
				}
			} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH_DEREGISTER {
				transaction := t.(PublisherTransaction)
				fmt.Println("\t$$$$$ publisherHandler receive: PUBLISH_DEREGISTER", string(msg.Body), "==> ", msg.TransactionId)
				transaction.AckDeregister(nil)
			} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_REGISTER {
				transaction := t.(SubscriberTransaction)
				fmt.Println("\t$$$$$ subscriberHandler receive: REGISTER", string(msg.Body), "==> ", msg.TransactionId)
				subs = transaction
				transaction.AckRegister(nil)
			} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_DEREGISTER {
				transaction := t.(SubscriberTransaction)
				fmt.Println("\t$$$$$ subscriberHandler receive: DEREGISTER", string(msg.Body), "==> ", msg.TransactionId)
				subs = nil
				transaction.AckDeregister(nil)
			} else {
				fmt.Println("\t$$$$$ publisherHandler receive: Bad message ", msg.InteractionStage, " -> ", "==> ", msg.TransactionId, string(msg.Body))
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
	err = op1.Register(msg1)
	if err != nil {
		t.Fatal("Error during publish register operation, ", err)
		return
	}
	fmt.Println("\t&&&&& Publisher registered")

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
	err = op2.Register(msg2)
	if err != nil {
		t.Fatal("Error during register operation, ", err)
		return
	}
	fmt.Println("\t&&&&& Subscriber registered")

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
	fmt.Println("\t&&&&& Publisher publish: OK")

	// Try to get Notify
	r, err := op2.GetNotify()
	fmt.Println("\t&&&&& Subscriber notified: OK, ", string(r.Body))

	// Do Deregister
	msg4 := &Message{
		UriTo: broker.Uri,
		Body:  []byte("deregister"),
	}
	err = op1.Deregister(msg4)
	if err != nil {
		t.Fatal("Error during publish deregister operation, ", err)
		return
	}
	fmt.Println("\t&&&&& Publisher deregister")

	// Do Deregister
	msg5 := &Message{
		UriTo: broker.Uri,
		Body:  []byte("deregister"),
	}
	err = op2.Deregister(msg5)
	if err != nil {
		t.Fatal("Error during deregister operation, ", err)
		return
	}
	fmt.Println("\t&&&&& Subscriber deregister")

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
