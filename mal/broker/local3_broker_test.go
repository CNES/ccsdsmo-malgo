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

/**
 * Test the embedded broker through the LocalBroker interface.
 */

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
	lt3_brokerpub_url   = "maltcp://127.0.0.1:16010"
	lt3_subscriber1_url = "maltcp://127.0.0.1:16011"
	lt3_subscriber2_url = "maltcp://127.0.0.1:16013"
)

var (
	lt3_running bool = true

	lt3_broker         *LocalBroker
	lt3_updtHandler    UpdateValueHandler
	lt3_brokerpub_ctx  *Context
	lt3_brokerpub_cctx *ClientContext

	lt3_sub1_ctx    *Context
	lt3_subscriber1 *ClientContext

	lt3_sub2_ctx    *Context
	lt3_subscriber2 *ClientContext

	lt3_sub1_not_cpt  int = 0
	lt3_sub1_updt_cpt int = 0

	lt3_sub2_not_cpt  int = 0
	lt3_sub2_updt_cpt int = 0

	lt3_subid1 = Identifier("MySubscription1")
	lt3_subid2 = Identifier("MySubscription2")
)

func closeLocalTest2BrokerPub() {
	lt3_brokerpub_cctx.Close()
	lt3_brokerpub_ctx.Close()
}

func newLocalTest2BrokerPub(t *testing.T) {
	var err error
	lt3_brokerpub_ctx, err = NewContext(lt3_brokerpub_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	lt3_brokerpub_cctx, err = NewClientContext(lt3_brokerpub_ctx, "brokerpub")
	if err != nil {
		t.Fatal("Error creating publisher, ", err)
		return
	}
	lt3_brokerpub_cctx.SetDomain(IdentifierList([]*Identifier{NewIdentifier("spacecraft1"), NewIdentifier("payload"), NewIdentifier("camera1")}))

	// Creates local broker
	lt3_updtHandler = NewGenericUpdateValueHandler(NullBlobList)
	lt3_broker, err = NewLocalBroker(lt3_brokerpub_cctx, lt3_updtHandler, 200, 1, 1, 1)
	if err != nil {
		t.Fatal("Error creating broker, ", err)
	}
}

func localTest2Pub1(t *testing.T) {
	ekpub1 := &EntityKey{NewIdentifier("key1"), NewLong(1), NewLong(1), NewLong(1)}
	ekpub2 := &EntityKey{NewIdentifier("key2"), NewLong(2), NewLong(2), NewLong(2)}
	var eklist = EntityKeyList([]*EntityKey{ekpub1, ekpub2})

	lt3_broker.PublishRegister(&eklist)
	fmt.Printf("pubop.Register OK\n")

	// Publish a first update
	updthdr1 := &UpdateHeader{*TimeNow(), *lt3_brokerpub_cctx.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
	updthdr2 := &UpdateHeader{*TimeNow(), *lt3_brokerpub_cctx.Uri, MAL_UPDATETYPE_CREATION, *ekpub2}
	updthdr3 := &UpdateHeader{*TimeNow(), *lt3_brokerpub_cctx.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
	updtHdrlist1 := UpdateHeaderList([]*UpdateHeader{updthdr1, updthdr2, updthdr3})

	updt1 := &Blob{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	updt2 := &Blob{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
	updt3 := &Blob{0, 1}
	updtlist1 := BlobList([]*Blob{updt1, updt2, updt3})

	lt3_broker.Publish(&updtHdrlist1, &updtlist1)
	fmt.Printf("pubop.Publish OK\n")

	time.Sleep(100 * time.Millisecond)

	// Publish a second update
	updthdr4 := &UpdateHeader{*TimeNow(), *lt3_brokerpub_cctx.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
	updthdr5 := &UpdateHeader{*TimeNow(), *lt3_brokerpub_cctx.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
	updtHdrlist2 := UpdateHeaderList([]*UpdateHeader{updthdr4, updthdr5})

	updt4 := &Blob{2, 3}
	updt5 := &Blob{4, 5, 6}
	updtlist2 := BlobList([]*Blob{updt4, updt5})

	lt3_broker.Publish(&updtHdrlist2, &updtlist2)
	fmt.Printf("pubop.Publish OK\n")

	// Deregisters publisher
	lt3_broker.PublishDeregister()

	fmt.Printf("Publisher#1 end\n")
}

var lt3_subop1 SubscriberOperation
var lt3_sbody1 Body

func newLocalTest2Sub1() error {
	var err error
	lt3_sub1_ctx, err = NewContext(lt3_subscriber1_url)
	if err != nil {
		return err
	}
	lt3_subscriber1, err = NewClientContext(lt3_sub1_ctx, "subscriber1")
	if err != nil {
		return err
	}
	lt3_subscriber1.SetDomain(IdentifierList([]*Identifier{NewIdentifier("spacecraft1"), NewIdentifier("payload")}))

	lt3_subop1 = lt3_subscriber1.NewSubscriberOperation(lt3_broker.Uri(), 200, 1, 1, 1)
	lt3_sbody1 = lt3_subop1.NewBody()

	domains := IdentifierList([]*Identifier{NewIdentifier("*")})
	eksub := &EntityKey{NewIdentifier("key1"), NewLong(0), NewLong(0), NewLong(0)}
	var erlist = EntityRequestList([]*EntityRequest{
		&EntityRequest{
			&domains, true, true, true, true, EntityKeyList([]*EntityKey{eksub}),
		},
	})
	subs := &Subscription{lt3_subid1, erlist}
	lt3_sbody1.EncodeLastParameter(subs, false)

	lt3_subop1.Register(lt3_sbody1)
	fmt.Printf("subop.Register OK\n")
	// Register is synchronous, we can clear buffer
	lt3_sbody1.Reset(true)

	return nil
}

func runLocalTest2Sub1(t *testing.T) {
	for lt3_running == true {
		// Try to get Notify
		r1, err := lt3_subop1.GetNotify()
		if err != nil {
			fmt.Printf("Subscriber#1, Error in GetNotify: %v\n", err)
			break
		}
		fmt.Printf("\t&&&&& Subscriber1 notified: %d\n", r1.TransactionId)
		lt3_sub1_not_cpt += 1

		id, err := r1.DecodeParameter(NullIdentifier)
		updtHdrlist, err := r1.DecodeParameter(NullUpdateHeaderList)
		updtlist, err := r1.DecodeLastParameter(NullBlobList, false)
		lt3_sub1_updt_cpt += len(*updtlist.(*BlobList))
		fmt.Printf("\t&&&&& Subscriber1 notified: OK, %s \n\t%+v \n\t%#v\n\n", id, updtHdrlist, updtlist)
	}

	if (lt3_sub1_not_cpt != 2) || (lt3_sub1_updt_cpt != 4) {
		t.Errorf("Subscriber#1, bad counters: %d %d", lt3_sub1_not_cpt, lt3_sub1_updt_cpt)
	}

	// Deregisters subscriber

	idlist := IdentifierList([]*Identifier{&lt3_subid1})
	lt3_sbody1.EncodeLastParameter(&idlist, false)
	lt3_subop1.Deregister(lt3_sbody1)
	fmt.Printf("\t&&&&&Subscriber#1, Deregistered\n")
	lt3_sbody1.Reset(true)

	lt3_subscriber1.Close()
	lt3_sub1_ctx.Close()
}

var lt3_subop2 SubscriberOperation
var lt3_sbody2 Body

func newLocalTest2Sub2() error {
	var err error
	lt3_sub2_ctx, err = NewContext(lt3_subscriber2_url)
	if err != nil {
		return err
	}
	lt3_subscriber2, err = NewClientContext(lt3_sub2_ctx, "subscriber2")
	if err != nil {
		return err
	}
	lt3_subscriber2.SetDomain(IdentifierList([]*Identifier{NewIdentifier("spacecraft1"), NewIdentifier("payload")}))

	lt3_subop2 = lt3_subscriber2.NewSubscriberOperation(lt3_broker.Uri(), 200, 1, 1, 1)
	lt3_sbody2 = lt3_subop2.NewBody()

	domains := IdentifierList([]*Identifier{NewIdentifier("camera1")})
	eksub1 := &EntityKey{NewIdentifier("key1"), NewLong(0), NewLong(0), NewLong(0)}
	eksub2 := &EntityKey{NewIdentifier("key2"), NewLong(0), NewLong(0), NewLong(0)}
	var erlist = EntityRequestList([]*EntityRequest{
		&EntityRequest{
			&domains, true, true, true, true, EntityKeyList([]*EntityKey{eksub1, eksub2}),
		},
	})
	subs := &Subscription{lt3_subid2, erlist}
	lt3_sbody2.EncodeLastParameter(subs, false)

	lt3_subop2.Register(lt3_sbody2)
	fmt.Printf("subop.Register OK\n")
	// Register is synchronous, we can clear buffer
	lt3_sbody2.Reset(true)

	return nil
}

func runLocalTest2Sub2(t *testing.T) {
	for lt3_running == true {
		// Try to get Notify
		r1, err := lt3_subop2.GetNotify()
		if err != nil {
			fmt.Printf("Subscriber#1, Error in GetNotify: %v\n", err)
			break
		}
		fmt.Printf("\t&&&&& Subscriber2 notified: %d\n", r1.TransactionId)
		lt3_sub2_not_cpt += 1

		id, err := r1.DecodeParameter(NullIdentifier)
		updtHdrlist, err := r1.DecodeParameter(NullUpdateHeaderList)
		updtlist, err := r1.DecodeLastParameter(NullBlobList, false)
		lt3_sub2_updt_cpt += len(*updtlist.(*BlobList))
		fmt.Printf("\t&&&&& Subscriber2 notified: OK, %s \n\t%+v \n\t%#v\n\n", id, updtHdrlist, updtlist)
	}

	if (lt3_sub2_not_cpt != 2) || (lt3_sub2_updt_cpt != 5) {
		t.Errorf("Subscriber#2, bad counters: %d %d", lt3_sub2_not_cpt, lt3_sub2_updt_cpt)
	}

	// Deregisters subscriber

	idlist := IdentifierList([]*Identifier{&lt3_subid2})
	lt3_sbody1.EncodeLastParameter(&idlist, false)
	lt3_subop1.Deregister(lt3_sbody1)
	fmt.Printf("\t&&&&&Subscriber#2, Deregistered\n")
	lt3_sbody1.Reset(true)

	lt3_subscriber2.Close()
	lt3_sub2_ctx.Close()
}

func Test2LocalPubSub(t *testing.T) {
	// Waits socket closing from previous test
	time.Sleep(250 * time.Millisecond)

	// Creates the broker
	newLocalTest2BrokerPub(t)
	defer closeLocalTest2BrokerPub()

	// Creates the subscribers and registers it
	err := newLocalTest2Sub1()
	if err != nil {
		t.Fatal("Error creating subscriber#1, ", err)
	}
	go runLocalTest2Sub1(t)

	err = newLocalTest2Sub2()
	if err != nil {
		t.Fatal("Error creating subscriber#2, ", err)
	}
	go runLocalTest2Sub2(t)

	// Waits for subscribers (notify reception)
	time.Sleep(1000 * time.Millisecond)

	// Creates the publishers and registers it
	go localTest2Pub1(t)

	// Waits for subscribers (notify reception)
	time.Sleep(1000 * time.Millisecond)

	fmt.Printf("##### Finish: %d %d\n", lt3_sub1_not_cpt, lt3_sub1_updt_cpt)
	fmt.Printf("##### Finish: %d %d\n", lt3_sub2_not_cpt, lt3_sub2_updt_cpt)

	lt3_subop1.Interrupt()
	lt3_subop2.Interrupt()

	// Wait for subscribers (closing)
	time.Sleep(1000 * time.Millisecond)
}
