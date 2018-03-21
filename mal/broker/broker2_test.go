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
	test2_varint = true

	test2_broker_url      = "maltcp://127.0.0.1:16000"
	test2_subscriber1_url = "maltcp://127.0.0.1:16001"
	test2_publisher1_url  = "maltcp://127.0.0.1:16002"
	test2_subscriber2_url = "maltcp://127.0.0.1:16003"
	test2_publisher2_url  = "maltcp://127.0.0.1:16004"
)

var (
	running bool = true

	test2_broker_ctx  *Context
	test2_broker      *BrokerHandler
	test2_pub1_ctx    *Context
	test2_publisher1  *ClientContext
	test2_pub2_ctx    *Context
	test2_publisher2  *ClientContext
	test2_sub1_ctx    *Context
	test2_subscriber1 *ClientContext
	test2_sub2_ctx    *Context
	test2_subscriber2 *ClientContext

	test2_sub1_not_cpt  int = 0
	test2_sub1_updt_cpt int = 0
	test2_sub2_not_cpt  int = 0
	test2_sub2_updt_cpt int = 0

	subid1 = Identifier("MySubscription1")
	subid2 = Identifier("MySubscription2")
)

func newTest2Broker() error {
	var err error
	test2_broker_ctx, err = NewContext(test2_broker_url)
	if err != nil {
		return err
	}

	cctx, err := NewClientContext(test2_broker_ctx, "broker")
	if err != nil {
		return err
	}

	updtHandler := NewBlobUpdateValueHandler()
	test2_broker, err = NewBroker(cctx, updtHandler, binary.VarintBinaryEncodingFactory)
	if err != nil {
		return err
	}

	return nil
}

func closeTest2Broker() {
	test2_broker.Close()
	test2_broker_ctx.Close()
}

func test2Pub1(t *testing.T) {
	var err error
	test2_pub1_ctx, err = NewContext(test2_publisher1_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}
	defer test2_pub1_ctx.Close()

	test2_publisher1, err = NewClientContext(test2_pub1_ctx, "publisher1")
	if err != nil {
		t.Fatal("Error creating publisher, ", err)
		return
	}
	defer test2_publisher1.Close()
	test2_publisher1.SetDomain(IdentifierList([]*Identifier{NewIdentifier("spacecraft1"), NewIdentifier("payload"), NewIdentifier("camera1")}))

	encoder := binary.NewBinaryEncoder(make([]byte, 0, 8192), test2_varint)

	pubop := test2_publisher1.NewPublisherOperation(test2_broker.Uri(), 200, 1, 1, 1)

	ekpub1 := &EntityKey{NewIdentifier("key1"), NewLong(1), NewLong(1), NewLong(1)}
	ekpub2 := &EntityKey{NewIdentifier("key2"), NewLong(2), NewLong(2), NewLong(2)}
	var eklist = EntityKeyList([]*EntityKey{ekpub1, ekpub2})
	eklist.Encode(encoder)

	pubop.Register(encoder.Body())
	encoder.Out.Reset(true)

	fmt.Printf("pubop.Register OK\n")

	// Publish a first update
	updthdr1 := &UpdateHeader{*TimeNow(), *test2_publisher1.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
	updthdr2 := &UpdateHeader{*TimeNow(), *test2_publisher1.Uri, MAL_UPDATETYPE_CREATION, *ekpub2}
	updthdr3 := &UpdateHeader{*TimeNow(), *test2_publisher1.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
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
	updthdr4 := &UpdateHeader{*TimeNow(), *test2_publisher1.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
	updthdr5 := &UpdateHeader{*TimeNow(), *test2_publisher1.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
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

	// Deregisters publisher
	pubop.Deregister(nil)
	fmt.Printf("Publisher#1 end\n")
}

func test2Pub2(t *testing.T) {
	var err error
	test2_pub2_ctx, err = NewContext(test2_publisher2_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}
	defer test2_pub2_ctx.Close()

	test2_publisher2, err = NewClientContext(test2_pub2_ctx, "publisher2")
	if err != nil {
		t.Fatal("Error creating publisher, ", err)
		return
	}
	defer test2_publisher2.Close()
	test2_publisher2.SetDomain(IdentifierList([]*Identifier{NewIdentifier("spacecraft1"), NewIdentifier("payload"), NewIdentifier("camera2")}))

	encoder := binary.NewBinaryEncoder(make([]byte, 0, 8192), test2_varint)

	pubop := test2_publisher2.NewPublisherOperation(test2_broker.Uri(), 200, 1, 1, 1)

	ekpub1 := &EntityKey{NewIdentifier("key1"), NewLong(1), NewLong(1), NewLong(1)}
	ekpub2 := &EntityKey{NewIdentifier("key2"), NewLong(1), NewLong(1), NewLong(1)}
	var eklist = EntityKeyList([]*EntityKey{ekpub1, ekpub2})
	eklist.Encode(encoder)

	pubop.Register(encoder.Body())
	encoder.Out.Reset(true)

	fmt.Printf("pubop.Register OK\n")

	// Publish a first update
	updthdr1 := &UpdateHeader{*TimeNow(), *test2_publisher2.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
	updthdr2 := &UpdateHeader{*TimeNow(), *test2_publisher2.Uri, MAL_UPDATETYPE_CREATION, *ekpub2}
	updthdr3 := &UpdateHeader{*TimeNow(), *test2_publisher2.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
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
	updthdr4 := &UpdateHeader{*TimeNow(), *test2_publisher2.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
	updthdr5 := &UpdateHeader{*TimeNow(), *test2_publisher2.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
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

	// Deregisters publisher
	pubop.Deregister(nil)
	fmt.Printf("Publisher#1 end\n")
}

var subop1 SubscriberOperation

func newTest2Sub1() error {
	var err error
	test2_sub1_ctx, err = NewContext(test2_subscriber1_url)
	if err != nil {
		return err
	}
	test2_subscriber1, err = NewClientContext(test2_sub1_ctx, "subscriber1")
	if err != nil {
		return err
	}
	test2_subscriber1.SetDomain(IdentifierList([]*Identifier{NewIdentifier("spacecraft1"), NewIdentifier("payload")}))

	encoder := binary.NewBinaryEncoder(make([]byte, 0, 8192), test2_varint)

	subop1 = test2_subscriber1.NewSubscriberOperation(test2_broker.Uri(), 200, 1, 1, 1)

	domains := IdentifierList([]*Identifier{NewIdentifier("*")})
	eksub := &EntityKey{NewIdentifier("key1"), NewLong(0), NewLong(0), NewLong(0)}
	var erlist = EntityRequestList([]*EntityRequest{
		&EntityRequest{
			&domains, true, true, true, true, EntityKeyList([]*EntityKey{eksub}),
		},
	})
	subs := &Subscription{subid1, erlist}
	subs.Encode(encoder)

	subop1.Register(encoder.Body())
	encoder.Out.Reset(true)

	fmt.Printf("subop.Register OK\n")
	return nil
}

func runTest2Sub1(t *testing.T) {
	for running == true {
		// Try to get Notify
		r1, err := subop1.GetNotify()
		if err != nil {
			fmt.Printf("Subscriber#1, Error in GetNotify: %v\n", err)
			break
		}
		fmt.Printf("\t&&&&& Subscriber1 notified: %d\n", r1.TransactionId)
		test2_sub1_not_cpt += 1
		decoder := binary.NewBinaryDecoder(r1.Body, varint)
		id, err := decoder.DecodeIdentifier()
		updtHdrlist, err := DecodeUpdateHeaderList(decoder)
		updtlist, err := DecodeUpdateList(decoder)
		test2_sub1_updt_cpt += len(*updtlist)
		fmt.Printf("\t&&&&& Subscriber1 notified: OK, %s \n\t%+v \n\t%#v\n\n", *id, updtHdrlist, updtlist)
	}

	if (test2_sub1_not_cpt != 4) || (test2_sub1_updt_cpt != 8) {
		t.Errorf("Subscriber#1, bad counters: %d %d", test2_sub1_not_cpt, test2_sub1_updt_cpt)
	}

	// Deregisters subscriber
	encoder := binary.NewBinaryEncoder(make([]byte, 0, 8192), test2_varint)

	idlist := IdentifierList([]*Identifier{&subid1})
	idlist.Encode(encoder)
	subop1.Deregister(encoder.Body())
	encoder.Out.Reset(true)

	fmt.Printf("\t&&&&&Subscriber#1, Deregistered\n")

	test2_subscriber1.Close()
	test2_sub1_ctx.Close()
}

var subop2 SubscriberOperation

func newTest2Sub2() error {
	var err error
	test2_sub2_ctx, err = NewContext(test2_subscriber2_url)
	if err != nil {
		return err
	}
	test2_subscriber2, err = NewClientContext(test2_sub2_ctx, "subscriber2")
	if err != nil {
		return err
	}
	test2_subscriber2.SetDomain(IdentifierList([]*Identifier{NewIdentifier("spacecraft1"), NewIdentifier("payload")}))

	encoder := binary.NewBinaryEncoder(make([]byte, 0, 8192), test2_varint)

	subop2 = test2_subscriber2.NewSubscriberOperation(test2_broker.Uri(), 200, 1, 1, 1)

	domains := IdentifierList([]*Identifier{NewIdentifier("camera2")})
	eksub1 := &EntityKey{NewIdentifier("key1"), NewLong(0), NewLong(0), NewLong(0)}
	eksub2 := &EntityKey{NewIdentifier("key2"), NewLong(0), NewLong(0), NewLong(0)}
	var erlist = EntityRequestList([]*EntityRequest{
		&EntityRequest{
			&domains, true, true, true, true, EntityKeyList([]*EntityKey{eksub1, eksub2}),
		},
	})
	subs := &Subscription{subid2, erlist}
	subs.Encode(encoder)

	subop2.Register(encoder.Body())
	encoder.Out.Reset(true)

	fmt.Printf("subop.Register OK\n")
	return nil
}

func runTest2Sub2(t *testing.T) {
	for running == true {
		// Try to get Notify
		r1, err := subop2.GetNotify()
		if err != nil {
			fmt.Printf("Subscriber#1, Error in GetNotify: %v\n", err)
			break
		}
		fmt.Printf("\t&&&&& Subscriber2 notified: %d\n", r1.TransactionId)
		test2_sub2_not_cpt += 1
		decoder := binary.NewBinaryDecoder(r1.Body, varint)
		id, err := decoder.DecodeIdentifier()
		updtHdrlist, err := DecodeUpdateHeaderList(decoder)
		updtlist, err := DecodeUpdateList(decoder)
		test2_sub2_updt_cpt += len(*updtlist)
		fmt.Printf("\t&&&&& Subscriber2 notified: OK, %s \n\t%+v \n\t%#v\n\n", *id, updtHdrlist, updtlist)
	}

	if (test2_sub2_not_cpt != 2) || (test2_sub2_updt_cpt != 5) {
		t.Errorf("Subscriber#2, bad counters: %d %d", test2_sub2_not_cpt, test2_sub2_updt_cpt)
	}

	// Deregisters subscriber
	encoder := binary.NewBinaryEncoder(make([]byte, 0, 8192), test2_varint)

	idlist := IdentifierList([]*Identifier{&subid2})
	idlist.Encode(encoder)
	subop2.Deregister(encoder.Body())
	encoder.Out.Reset(true)

	fmt.Printf("\t&&&&&Subscriber#2, Deregistered\n")

	test2_subscriber2.Close()
	test2_sub2_ctx.Close()
}

func Test2PubSub(t *testing.T) {
	// Waits socket closing from previous test
	time.Sleep(250 * time.Millisecond)

	// Creates the broker
	err := newTest2Broker()
	if err != nil {
		t.Fatal("Error creating broker, ", err)
		return
	}
	defer closeTest2Broker()

	// Creates the subscribers and registers it
	err = newTest2Sub1()
	if err != nil {
		t.Fatal("Error creating subscriber#1, ", err)
	}
	go runTest2Sub1(t)

	err = newTest2Sub2()
	if err != nil {
		t.Fatal("Error creating subscriber#2, ", err)
	}
	go runTest2Sub2(t)

	// Waits for subscribers (notify reception)
	time.Sleep(500 * time.Millisecond)

	// Creates the publishers and registers it
	go test2Pub1(t)
	go test2Pub2(t)

	// Waits for subscribers (notify reception)
	time.Sleep(500 * time.Millisecond)

	fmt.Printf("##### Finish: %d %d\n", test2_sub1_not_cpt, test2_sub1_updt_cpt)
	fmt.Printf("##### Finish: %d %d\n", test2_sub2_not_cpt, test2_sub2_updt_cpt)

	subop1.Interrupt()
	subop2.Interrupt()

	// Wait for subscribers (closing)
	time.Sleep(1000 * time.Millisecond)
}
