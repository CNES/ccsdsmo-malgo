/**
 * MIT License
 *
 * Copyright (c) 2018 CNES
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
package broker_test

import (
	"fmt"
	. "github.com/ccsdsmo/malgo/mal"
	. "github.com/ccsdsmo/malgo/mal/api"
	. "github.com/ccsdsmo/malgo/mal/broker"
	"github.com/ccsdsmo/malgo/mal/encoding/binary"
	_ "github.com/ccsdsmo/malgo/mal/transport/tcp" // Needed to initialize TCP transport factory
	"testing"
	"time"
)

const (
	varint = true

	subscriber_url = "maltcp://127.0.0.1:16001"
	publisher_url  = "maltcp://127.0.0.1:16002"
	broker_url     = "maltcp://127.0.0.1:16003"
)

func TestPubSub(t *testing.T) {
	// Waits socket closing from previous test
	time.Sleep(250 * time.Millisecond)

	// Creates the broker

	ctx, err := NewContext(broker_url)
	if err != nil {
		t.Fatal("Error creating broker context, ", err)
		return
	}

	broker, err := NewBroker(ctx, "broker")
	if err != nil {
		t.Fatal("Error creating broker, ", err)
		return
	}
	defer broker.Close()

	encoder := binary.NewBinaryEncoder(make([]byte, 0, 8192), varint)

	// Creates the publisher and registers it

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
	defer publisher.Close()
	publisher.SetDomain(IdentifierList([]*Identifier{NewIdentifier("spacecraft1"), NewIdentifier("payload"), NewIdentifier("camera")}))

	pubop := publisher.NewPublisherOperation(broker.Uri(), 200, 1, 1, 1)

	ekpub1 := &EntityKey{NewIdentifier("key1"), NewLong(1), NewLong(1), NewLong(1)}
	ekpub2 := &EntityKey{NewIdentifier("key2"), NewLong(1), NewLong(1), NewLong(1)}
	var eklist = EntityKeyList([]*EntityKey{ekpub1, ekpub2})
	eklist.Encode(encoder)

	pubop.Register(encoder.Body())
	encoder.Out.Reset(true)

	fmt.Printf("pubop.Register OK\n")

	// Creates the subscriber and registers it

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
	defer subscriber.Close()
	subscriber.SetDomain(IdentifierList([]*Identifier{NewIdentifier("spacecraft1"), NewIdentifier("payload")}))

	subop := subscriber.NewSubscriberOperation(broker.Uri(), 200, 1, 1, 1)

	domains := IdentifierList([]*Identifier{NewIdentifier("*")})
	eksub := &EntityKey{NewIdentifier("key1"), NewLong(0), NewLong(0), NewLong(0)}
	var erlist = EntityRequestList([]*EntityRequest{
		&EntityRequest{
			&domains, true, true, true, true, EntityKeyList([]*EntityKey{eksub}),
		},
	})
	var subid = Identifier("MySubscription")
	subs := &Subscription{subid, erlist}
	subs.Encode(encoder)

	subop.Register(encoder.Body())
	encoder.Out.Reset(true)

	fmt.Printf("subop.Register OK\n")

	// Publish a first update
	updthdr1 := &UpdateHeader{*TimeNow(), *publisher.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
	updthdr2 := &UpdateHeader{*TimeNow(), *publisher.Uri, MAL_UPDATETYPE_CREATION, *ekpub2}
	updthdr3 := &UpdateHeader{*TimeNow(), *publisher.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
	updtHdrlist1 := UpdateHeaderList([]*UpdateHeader{updthdr1, updthdr2, updthdr3})

	updt1 := &Blob{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	updt2 := &Blob{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
	updt3 := &Blob{0, 1}
	updtlist1 := UpdateList([]*Blob{updt1, updt2, updt3})

	updtHdrlist1.Encode(encoder)
	updtlist1.Encode(encoder)
	body1 := encoder.Body()
	encoder.Out.Reset(true)
	//	fmt.Printf("\n\nBody=%p %v\n\n", body1, body1)
	pubop.Publish(body1)

	fmt.Printf("pubop.Publish OK\n")

	time.Sleep(100 * time.Millisecond)

	// Publish a second update
	updthdr4 := &UpdateHeader{*TimeNow(), *publisher.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
	updthdr5 := &UpdateHeader{*TimeNow(), *publisher.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
	updtHdrlist2 := UpdateHeaderList([]*UpdateHeader{updthdr4, updthdr5})

	updt4 := &Blob{2, 3}
	updt5 := &Blob{4, 5, 6}
	updtlist2 := UpdateList([]*Blob{updt4, updt5})

	updtHdrlist2.Encode(encoder)
	updtlist2.Encode(encoder)
	body2 := encoder.Body()
	encoder.Out.Reset(true)
	//	fmt.Printf("\n\nBody=%p %v\n\n", body2, body2)
	pubop.Publish(body2)

	// Try to get Notify
	r1, err := subop.GetNotify()
	fmt.Printf("\t&&&&& Subscriber notified: %d\n", r1.TransactionId)
	decoder := binary.NewBinaryDecoder(r1.Body, varint)
	id, err := decoder.DecodeIdentifier()
	updtHdrlist, err := DecodeUpdateHeaderList(decoder)
	updtlist, err := DecodeUpdateList(decoder)
	fmt.Printf("\t&&&&& Subscriber notified: OK, %s \n\t%+v \n\t%#v\n\n", *id, updtHdrlist, updtlist)

	// Try to get Notify
	r2, err := subop.GetNotify()
	fmt.Printf("\t&&&&& Subscriber notified: %d\n", r2.TransactionId)
	decoder = binary.NewBinaryDecoder(r2.Body, varint)
	id, err = decoder.DecodeIdentifier()
	updtHdrlist, err = DecodeUpdateHeaderList(decoder)
	updtlist, err = DecodeUpdateList(decoder)
	fmt.Printf("\t&&&&& Subscriber notified: OK, %s \n\t%+v \n\t%#v\n\n", *id, updtHdrlist, updtlist)

	// Deregisters publisher
	pubop.Deregister(nil)

	// Deregisters subscriber

	idlist := IdentifierList([]*Identifier{&subid})
	idlist.Encode(encoder)
	subop.Deregister(encoder.Body())
	encoder.Out.Reset(true)

	// Waits for socket close
	time.Sleep(250 * time.Millisecond)
}
