/**
 * MIT License
 *
 * Copyright (c) 2017 CNES
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

// Note (AF): This API is now deprecated, a merged API allows to use either handler
// based and interface based interfaces in a same way.

import (
	"fmt"
	. "github.com/ccsdsmo/malgo/src/mal"
	. "github.com/ccsdsmo/malgo/src/mal/api2"
	_ "github.com/ccsdsmo/malgo/src/mal/transport/tcp" // Needed to initialize TCP transport factory
	"testing"
	"time"
)

const (
	subscriber_url = "github.com/ccsdsmo/malgo/src/maltcp://127.0.0.1:16001"
	publisher_url  = "github.com/ccsdsmo/malgo/src/maltcp://127.0.0.1:16002"
	broker_url     = "github.com/ccsdsmo/malgo/src/maltcp://127.0.0.1:16003"
)

// ########## ########## ########## ########## ########## ########## ########## ##########

type BrokerContext struct {
	pctx *ProviderContext
}

func NewBroker(ctx *Context, service string) (*BrokerContext, error) {
	pctx, err := NewProviderContext(ctx, service)
	if err != nil {
		return nil, err
	}
	broker := &BrokerContext{pctx: pctx}
	// Registers the broker handler
	pctx.RegisterBrokerHandler(200, 1, 1, 1, broker)

	return broker, nil
}

func (handler *BrokerContext) OnRegister(msg *Message, tx SubscriberTransaction) error {
	fmt.Println("\n\n##########\n# OnRegister:\n##########\n\n")
	tx.RegisterAck(nil)
	return nil
}

func (handler *BrokerContext) OnDeregister(msg *Message, tx SubscriberTransaction) error {
	fmt.Println("\n\n##########\n# OnDeregister:\n##########\n\n")
	tx.DeregisterAck(nil)
	return nil
}

func (handler *BrokerContext) OnPublishRegister(msg *Message, tx PublisherTransaction) error {
	fmt.Println("\n\n##########\n# OnPublishRegister:\n##########\n\n")
	tx.RegisterAck(nil)
	return nil
}

func (handler *BrokerContext) OnPublishDeregister(msg *Message, tx PublisherTransaction) error {
	fmt.Println("\n\n##########\n# OnPublishDeregister:\n##########\n\n")
	tx.DeregisterAck(nil)
	return nil
}

func (handler *BrokerContext) OnPublish(msg *Message, tx PublisherTransaction) error {
	fmt.Println("\n\n##########\n# OnPublish:\n##########\n\n")
	return nil
}

// Test TCP transport Pub/Sub Interaction using the high level API
func TestPubSub(t *testing.T) {
	broker_ctx, err := NewContext(broker_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}
	broker, err := NewBroker(broker_ctx, "broker")
	t.Log("Broker: ", broker)

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
	pubop, err := publisher.NewPublisherOperation(200, 1, 1, 1)
	// TODO (AF): Build PublishRegister message
	pubregmsg := &Message{
		UriTo: broker.pctx.Uri,
		Body:  []byte("publish register"),
	}
	pubop.Register(pubregmsg)

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
	subop, err := subscriber.NewSubscriberOperation(200, 1, 1, 1)
	// TODO (AF): Build Register message
	regmsg := &Message{
		UriTo: broker.pctx.Uri,
		Body:  []byte("register"),
	}
	subop.Register(regmsg)

	//	nbmsg := 0

	pubmsg1 := &Message{
		UriTo: broker.pctx.Uri,
		Body:  []byte("publish #1"),
	}
	pubop.Publish(pubmsg1)

	pubmsg2 := &Message{
		UriTo: broker.pctx.Uri,
		Body:  []byte("publish #2"),
	}
	pubop.Publish(pubmsg2)

	pubderegmsg := &Message{
		UriTo: broker.pctx.Uri,
		Body:  []byte("publish deregister"),
	}
	pubop.Deregister(pubderegmsg)

	deregmsg := &Message{
		UriTo: broker.pctx.Uri,
		Body:  []byte("deregister"),
	}
	subop.Deregister(deregmsg)

	// Waits for socket close
	time.Sleep(250 * time.Millisecond)
}
