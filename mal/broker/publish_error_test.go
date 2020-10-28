/**
 * MIT License
 *
 * Copyright (c) 2020 CNES
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
	. "github.com/CNES/ccsdsmo-malgo/mal"
	. "github.com/CNES/ccsdsmo-malgo/mal/api"
	. "github.com/CNES/ccsdsmo-malgo/mal/broker"
	_ "github.com/CNES/ccsdsmo-malgo/mal/transport/tcp" // Needed to initialize TCP transport factory
	"testing"
	"time"
)

const (
	test3_varint = false

	test3_broker_url      = "maltcp://127.0.0.1:16000"
	test3_subscriber1_url = "maltcp://127.0.0.1:16001"
	test3_publisher1_url  = "maltcp://127.0.0.1:16002"
	test3_subscriber2_url = "maltcp://127.0.0.1:16003"
	test3_publisher2_url  = "maltcp://127.0.0.1:16004"
)

var (
	test3_running bool = true

	test3_broker_ctx  *Context
	test3_broker      *BrokerHandler
	test3_pub1_ctx    *Context
	test3_publisher1  *ClientContext
	test3_pub2_ctx    *Context
	test3_publisher2  *ClientContext
	test3_sub1_ctx    *Context
	test3_subscriber1 *ClientContext
	test3_sub2_ctx    *Context
	test3_subscriber2 *ClientContext

	test3_sub1_not_cpt  int = 0
	test3_sub1_updt_cpt int = 0
	test3_sub2_not_cpt  int = 0
	test3_sub2_updt_cpt int = 0

	test3_subid1 = Identifier("MySubscription1")
	test3_subid2 = Identifier("MySubscription2")
)

func newTest3Broker() error {
	var err error
	test3_broker_ctx, err = NewContext(test3_broker_url)
	if err != nil {
		return err
	}

	cctx, err := NewClientContext(test3_broker_ctx, "broker")
	if err != nil {
		return err
	}

	updtHandler := NewBlobUpdateValueHandler()
	if test3_varint {
		test3_broker, err = NewBroker(cctx, updtHandler, 200, 1, 1, 1)
	} else {
		test3_broker, err = NewBroker(cctx, updtHandler, 200, 1, 1, 1)
	}
	if err != nil {
		return err
	}

	return nil
}

func closeTest3Broker() {
	test3_broker.Close()
	test3_broker_ctx.Close()
}

func test3Pub1(t *testing.T) {
	var err error
	test3_pub1_ctx, err = NewContext(test3_publisher1_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}
	defer test3_pub1_ctx.Close()

	test3_publisher1, err = NewClientContext(test3_pub1_ctx, "publisher1")
	if err != nil {
		t.Fatal("Error creating publisher, ", err)
		return
	}
	defer test3_publisher1.Close()
	test3_publisher1.SetDomain(IdentifierList([]*Identifier{NewIdentifier("spacecraft1"), NewIdentifier("payload"), NewIdentifier("camera1")}))

	pubop := test3_publisher1.NewPublisherOperation(test3_broker.Uri(), 200, 1, 1, 1)
	pbody := pubop.NewBody()

	ekpub1 := &EntityKey{NewIdentifier("key1"), NewLong(1), NewLong(1), NewLong(1)}
	ekpub2 := &EntityKey{NewIdentifier("key2"), NewLong(2), NewLong(2), NewLong(2)}
	//	ekpub3 := &EntityKey{NewIdentifier("key3"), NewLong(0), NewLong(0), NewLong(0)}
	var eklist = EntityKeyList([]*EntityKey{ekpub1, ekpub2})
	pbody.EncodeLastParameter(&eklist, false)
	pubop.Register(pbody)
	fmt.Printf("pubop.Register OK\n")
	// Register is synchronous, we can reuse body
	pbody.Reset(true)

	// Publish a first update
	updthdr1 := &UpdateHeader{*TimeNow(), *test3_publisher1.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
	updthdr2 := &UpdateHeader{*TimeNow(), *test3_publisher1.Uri, MAL_UPDATETYPE_CREATION, *ekpub2}
	updthdr3 := &UpdateHeader{*TimeNow(), *test3_publisher1.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
	updtHdrlist1 := UpdateHeaderList([]*UpdateHeader{updthdr1, updthdr2, updthdr3})

	updt1 := &Blob{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	updt2 := &Blob{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
	updt3 := &Blob{0, 1}
	updtlist1 := BlobList([]*Blob{updt1, updt2, updt3})

	pbody1 := pubop.NewBody()
	pbody1.EncodeParameter(&updtHdrlist1)
	pbody1.EncodeLastParameter(&updtlist1, false)
	pubop.Publish(pbody1)

	fmt.Printf("pubop.Publish OK\n")

	time.Sleep(100 * time.Millisecond)

	// Publish a second update
	updthdr4 := &UpdateHeader{*TimeNow(), *test3_publisher1.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
	updthdr5 := &UpdateHeader{*TimeNow(), *test3_publisher1.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
	updtHdrlist2 := UpdateHeaderList([]*UpdateHeader{updthdr4, updthdr5})

	updt4 := &Blob{2, 3}
	updt5 := &Blob{4, 5, 6}
	// Number of values > number of keys => PublishError
	updtlist2 := BlobList([]*Blob{updt4, updt5, updt5})

	pbody2 := pubop.NewBody()
	pbody2.EncodeParameter(&updtHdrlist2)
	pbody2.EncodeLastParameter(&updtlist2, false)
	pubop.Publish(pbody2)

	time.Sleep(500 * time.Millisecond)
	merr, err := pubop.GetPublishError()
	if err != nil {
		t.Fatal("Error getting PublishError, ", err)
	}
	if merr == nil {
		t.Fatal("Error getting PublishError, no waiting PE message", nil)
	} else {
		pe_id, err := merr.DecodeParameter(NullUInteger)
		if err != nil {
			t.Fatal("Error decoding PublishError id, ", err)
		}
		pe_str, err := merr.DecodeLastParameter(NullString, false)
		if err != nil {
			t.Fatal("Error decoding PublishError str, ", err)
		}

		fmt.Printf("\t$$$$$ Publisher1 Get a PublishError: %v, %d, %s\n", merr, *pe_id.(*UInteger), *pe_str.(*String))
	}

	// Deregisters publisher
	pubop.Deregister(nil)
	fmt.Printf("Publisher#1 end\n")
}

func test3Pub2(t *testing.T) {
	var err error
	test3_pub2_ctx, err = NewContext(test3_publisher2_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
	}
	defer test3_pub2_ctx.Close()

	test3_publisher2, err = NewClientContext(test3_pub2_ctx, "publisher2")
	if err != nil {
		t.Fatal("Error creating publisher, ", err)
	}
	defer test3_publisher2.Close()
	test3_publisher2.SetDomain(IdentifierList([]*Identifier{NewIdentifier("spacecraft1"), NewIdentifier("payload"), NewIdentifier("camera2")}))

	pubop := test3_publisher2.NewPublisherOperation(test3_broker.Uri(), 200, 1, 1, 1)
	pbody := pubop.NewBody()

	ekpub1 := &EntityKey{NewIdentifier("key1"), NewLong(1), NewLong(1), NewLong(1)}
	ekpub2 := &EntityKey{NewIdentifier("key2"), NewLong(1), NewLong(1), NewLong(1)}
	var eklist = EntityKeyList([]*EntityKey{ekpub1, ekpub2})
	pbody.EncodeLastParameter(&eklist, false)

	pubop.Register(pbody)
	fmt.Printf("pubop.Register OK\n")
	// Register is synchronous, we can reuse body
	pbody.Reset(true)

	// Publish a first update
	updthdr1 := &UpdateHeader{*TimeNow(), *test3_publisher2.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
	updthdr2 := &UpdateHeader{*TimeNow(), *test3_publisher2.Uri, MAL_UPDATETYPE_CREATION, *ekpub2}
	updthdr3 := &UpdateHeader{*TimeNow(), *test3_publisher2.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
	updtHdrlist1 := UpdateHeaderList([]*UpdateHeader{updthdr1, updthdr2, updthdr3})

	updt1 := &Blob{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	updt2 := &Blob{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
	updt3 := &Blob{0, 1}
	updtlist1 := BlobList([]*Blob{updt1, updt2, updt3})

	pbody1 := pubop.NewBody()
	pbody1.EncodeParameter(&updtHdrlist1)
	pbody1.EncodeLastParameter(&updtlist1, false)
	pubop.Publish(pbody1)

	fmt.Printf("pubop.Publish OK\n")

	time.Sleep(100 * time.Millisecond)

	// Publish a second update
	updthdr4 := &UpdateHeader{*TimeNow(), *test3_publisher2.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
	updthdr5 := &UpdateHeader{*TimeNow(), *test3_publisher2.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
	updtHdrlist2 := UpdateHeaderList([]*UpdateHeader{updthdr4, updthdr5})

	updt4 := &Blob{2, 3}
	updt5 := &Blob{4, 5, 6}
	updtlist2 := BlobList([]*Blob{updt4, updt5})

	pbody2 := pubop.NewBody()
	pbody2.EncodeParameter(&updtHdrlist2)
	pbody2.EncodeLastParameter(&updtlist2, false)
	pubop.Publish(pbody2)

	// Deregisters publisher
	pubop.Deregister(nil)
	fmt.Printf("Publisher#1 end\n")
}

var test3_subop1 SubscriberOperation
var test3_sbody1 Body

func newTest3Sub1() error {
	var err error
	test3_sub1_ctx, err = NewContext(test3_subscriber1_url)
	if err != nil {
		return err
	}
	test3_subscriber1, err = NewClientContext(test3_sub1_ctx, "subscriber1")
	if err != nil {
		return err
	}
	test3_subscriber1.SetDomain(IdentifierList([]*Identifier{NewIdentifier("spacecraft1"), NewIdentifier("payload")}))

	subop1 = test3_subscriber1.NewSubscriberOperation(test3_broker.Uri(), 200, 1, 1, 1)
	sbody1 = subop1.NewBody()

	domains := IdentifierList([]*Identifier{NewIdentifier("*")})
	eksub := &EntityKey{NewIdentifier("key1"), NewLong(0), NewLong(0), NewLong(0)}
	var erlist = EntityRequestList([]*EntityRequest{
		&EntityRequest{
			&domains, true, true, true, true, EntityKeyList([]*EntityKey{eksub}),
		},
	})
	subs := &Subscription{test3_subid1, erlist}
	sbody1.EncodeLastParameter(subs, false)

	subop1.Register(sbody1)
	fmt.Printf("subop.Register OK\n")
	// Register is synchronous, we can clear buffer
	sbody1.Reset(true)

	return nil
}

func runTest3Sub1(t *testing.T) {
	for test3_running == true {
		// Try to get Notify
		r1, err := subop1.GetNotify()
		if err != nil {
			fmt.Printf("Subscriber#1, Error in GetNotify: %v\n", err)
			break
		}
		fmt.Printf("\t&&&&& Subscriber1 notified: %d\n", r1.TransactionId)
		test3_sub1_not_cpt += 1

		id, err := r1.DecodeParameter(NullIdentifier)
		updtHdrlist, err := r1.DecodeParameter(NullUpdateHeaderList)
		updtlist, err := r1.DecodeLastParameter(NullBlobList, false)
		test3_sub1_updt_cpt += len(*updtlist.(*BlobList))
		fmt.Printf("\t&&&&& Subscriber1 notified: OK, %s \n\t%+v \n\t%#v\n\n", id, updtHdrlist, updtlist)
	}

	if (test3_sub1_not_cpt != 3) || (test3_sub1_updt_cpt != 6) {
		t.Errorf("Subscriber#1, bad counters: %d %d", test3_sub1_not_cpt, test3_sub1_updt_cpt)
	}

	// Deregisters subscriber

	idlist := IdentifierList([]*Identifier{&test3_subid1})
	sbody1.EncodeLastParameter(&idlist, false)
	subop1.Deregister(sbody1)
	fmt.Printf("\t&&&&&Subscriber#1, Deregistered\n")
	sbody1.Reset(true)

	test3_subscriber1.Close()
	test3_sub1_ctx.Close()
}

var test3_subop2 SubscriberOperation
var test3_sbody2 Body

func newTest3Sub2() error {
	var err error
	test3_sub2_ctx, err = NewContext(test3_subscriber2_url)
	if err != nil {
		return err
	}
	test3_subscriber2, err = NewClientContext(test3_sub2_ctx, "subscriber2")
	if err != nil {
		return err
	}
	test3_subscriber2.SetDomain(IdentifierList([]*Identifier{NewIdentifier("spacecraft1"), NewIdentifier("payload")}))

	subop2 = test3_subscriber2.NewSubscriberOperation(test3_broker.Uri(), 200, 1, 1, 1)
	sbody2 = subop2.NewBody()

	domains := IdentifierList([]*Identifier{NewIdentifier("camera2")})
	eksub1 := &EntityKey{NewIdentifier("key1"), NewLong(0), NewLong(0), NewLong(0)}
	eksub2 := &EntityKey{NewIdentifier("key2"), NewLong(0), NewLong(0), NewLong(0)}
	var erlist = EntityRequestList([]*EntityRequest{
		&EntityRequest{
			&domains, true, true, true, true, EntityKeyList([]*EntityKey{eksub1, eksub2}),
		},
	})
	subs := &Subscription{test3_subid2, erlist}
	sbody2.EncodeLastParameter(subs, false)

	subop2.Register(sbody2)
	fmt.Printf("subop.Register OK\n")
	// Register is synchronous, we can clear buffer
	sbody2.Reset(true)

	return nil
}

func runTest3Sub2(t *testing.T) {
	for test3_running == true {
		// Try to get Notify
		r1, err := subop2.GetNotify()
		if err != nil {
			fmt.Printf("Subscriber#1, Error in GetNotify: %v\n", err)
			break
		}
		fmt.Printf("\t&&&&& Subscriber2 notified: %d\n", r1.TransactionId)
		test3_sub2_not_cpt += 1

		id, err := r1.DecodeParameter(NullIdentifier)
		updtHdrlist, err := r1.DecodeParameter(NullUpdateHeaderList)
		updtlist, err := r1.DecodeLastParameter(NullBlobList, false)
		test3_sub2_updt_cpt += len(*updtlist.(*BlobList))
		fmt.Printf("\t&&&&& Subscriber2 notified: OK, %s \n\t%+v \n\t%#v\n\n", id, updtHdrlist, updtlist)
	}

	if (test3_sub2_not_cpt != 2) || (test3_sub2_updt_cpt != 5) {
		t.Errorf("Subscriber#2, bad counters: %d %d", test3_sub2_not_cpt, test3_sub2_updt_cpt)
	}

	// Deregisters subscriber

	idlist := IdentifierList([]*Identifier{&test3_subid2})
	sbody1.EncodeLastParameter(&idlist, false)
	subop1.Deregister(sbody1)
	fmt.Printf("\t&&&&&Subscriber#2, Deregistered\n")
	sbody1.Reset(true)

	test3_subscriber2.Close()
	test3_sub2_ctx.Close()
}

func Test3PubSub(t *testing.T) {
	// Waits socket closing from previous test
	time.Sleep(250 * time.Millisecond)

	// Creates the broker
	err := newTest3Broker()
	if err != nil {
		t.Fatal("Error creating broker, ", err)
		return
	}
	defer closeTest3Broker()

	// Creates the subscribers and registers it
	err = newTest3Sub1()
	if err != nil {
		t.Fatal("Error creating subscriber#1, ", err)
	}
	go runTest3Sub1(t)

	err = newTest3Sub2()
	if err != nil {
		t.Fatal("Error creating subscriber#2, ", err)
	}
	go runTest3Sub2(t)

	// Waits for subscribers (notify reception)
	time.Sleep(500 * time.Millisecond)

	// Creates the publishers and registers it
	go test3Pub1(t)
	go test3Pub2(t)

	// Waits for subscribers (notify reception)
	time.Sleep(2000 * time.Millisecond)

	fmt.Printf("##### Finish: %d %d\n", test3_sub1_not_cpt, test3_sub1_updt_cpt)
	fmt.Printf("##### Finish: %d %d\n", test3_sub2_not_cpt, test3_sub2_updt_cpt)

	subop1.Interrupt()
	subop2.Interrupt()

	// Wait for subscribers (closing)
	time.Sleep(1000 * time.Millisecond)
}
