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

import (
	. "mal"
	. "mal/api2"
	_ "mal/transport/tcp" // Needed to initialize TCP transport factory
	"testing"
	"time"
)

const (
	//	maltcp1 = "maltcp://127.0.0.1:16001"
	//	maltcp2 = "maltcp://127.0.0.1:16002"
	maltcp3 = "maltcp://127.0.0.1:16003"
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
	// TODO (AF): Registers the broker handler
	//	pctx.RegisterBrokerHandler(200, 1, 1, 1, broker)

	return broker, nil
}

func (handler *BrokerContext) OnRegister(msg *Message, transaction SubscriberTransaction) error {
	return nil
}

func (handler *BrokerContext) OnDeregister(msg *Message, transaction SubscriberTransaction) error {
	return nil
}

func (handler *BrokerContext) OnPublishRegister(msg *Message, transaction PublisherTransaction) error {
	return nil
}

func (handler *BrokerContext) OnPublishDeregister(msg *Message, transaction PublisherTransaction) error {
	return nil
}

func (handler *BrokerContext) OnPublish(msg *Message, transaction PublisherTransaction) error {
	return nil
}

// Test TCP transport Pub/Sub Interaction using the high level API
func TestPubSub(t *testing.T) {
	ctx1, err := NewContext(maltcp1)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}
	broker, err := NewBroker(ctx1, "broker")

	// TODO (AF): Creates broker

	ctx2, err := NewContext(maltcp2)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	publisher, err := NewOperationContext(ctx2, "publisher")
	if err != nil {
		t.Fatal("Error creating publisher, ", err)
		return
	}
	pubop, err := publisher.NewPublisherOperation(200, 1, 1, 1)
	// TODO (AF): Build PublishRegister message
	pubmsg := &Message{}
	pubop.PublishRegister(pubmsg)

	ctx3, err := NewContext(maltcp3)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	subscriber, err := NewOperationContext(ctx3, "subscriber")
	if err != nil {
		t.Fatal("Error creating subscriber, ", err)
		return
	}
	subop, err := subscriber.NewSubscriberOperation(200, 1, 1, 1)
	// TODO (AF): Build Register message
	submsg := &Message{}
	subop.Register(submsg)

	//	nbmsg := 0

	// Waits for socket close
	time.Sleep(250 * time.Millisecond)
}
