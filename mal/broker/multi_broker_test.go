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
	"sync"
	"testing"
	"time"
)

const (
	mt1_brokerpub_url   = "maltcp://127.0.0.1:16010"
	mt1_subscriber1_url = "maltcp://127.0.0.1:16011"
	mt1_subscriber2_url = "maltcp://127.0.0.1:16013"
)

var (
	mt1_running bool = true
	mt1_wg sync.WaitGroup

	mt1_broker         *LocalBroker
	mt1_updtHandler    UpdateValueHandler
	mt1_brokerpub_ctx  *Context
	mt1_brokerpub_cctx *ClientContext

	mt1_sub1_ctx    *Context
	mt1_subscriber1 *ClientContext

	mt1_sub2_ctx    *Context
	mt1_subscriber2 *ClientContext

	mt1_sub1_not_cpt  int = 0
	mt1_sub1_updt_cpt int = 0

	mt1_sub2_not_cpt  int = 0
	mt1_sub2_updt_cpt int = 0

	mt1_subid1 = Identifier("MySubscription1")
	mt1_subid2 = Identifier("MySubscription2")
)

func closeLocalTestM1BrokerPub() {
	mt1_brokerpub_cctx.Close()
	mt1_brokerpub_ctx.Close()
}

func newLocalTestM1BrokerPub(t *testing.T) {
	var err error
	mt1_brokerpub_ctx, err = NewContext(mt1_brokerpub_url)
	if err != nil {
		t.Fatal("Error creating context, ", err)
		return
	}

	mt1_brokerpub_cctx, err = NewClientContext(mt1_brokerpub_ctx, "brokerpub")
	if err != nil {
		t.Fatal("Error creating publisher, ", err)
		return
	}
	mt1_brokerpub_cctx.SetDomain(IdentifierList([]*Identifier{NewIdentifier("spacecraft1"), NewIdentifier("payload"), NewIdentifier("camera1")}))

	// Creates local broker
	mt1_updtHandler = NewGenericUpdateValueHandler(NullIntegerList, NullStringList)
	mt1_broker, err = NewLocalBroker(mt1_brokerpub_cctx, mt1_updtHandler, 200, 1, 1, 1)
	if err != nil {
		t.Fatal("Error creating broker, ", err)
	}
}

func localTestM1Pub1(t *testing.T) {
	defer mt1_wg.Done()
	ekpub1 := &EntityKey{NewIdentifier("key1"), NewLong(1), NewLong(1), NewLong(1)}
	ekpub2 := &EntityKey{NewIdentifier("key2"), NewLong(2), NewLong(2), NewLong(2)}
	var eklist = EntityKeyList([]*EntityKey{ekpub1, ekpub2})

	mt1_broker.PublishRegister(&eklist)
	fmt.Printf("pubop.Register OK\n")

	// Publish a first update
	updthdr1 := &UpdateHeader{*TimeNow(), *mt1_brokerpub_cctx.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
	updthdr2 := &UpdateHeader{*TimeNow(), *mt1_brokerpub_cctx.Uri, MAL_UPDATETYPE_CREATION, *ekpub2}
	updthdr3 := &UpdateHeader{*TimeNow(), *mt1_brokerpub_cctx.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
	updtHdrlist1 := UpdateHeaderList([]*UpdateHeader{updthdr1, updthdr2, updthdr3})

	updtlistI1 := IntegerList([]*Integer{NewInteger(1), NewInteger(2), NewInteger(3)})
	updtlistS1 := StringList([]*String{NewString("a"), NewString("b"), NewString("c")})

	mt1_broker.Publish(&updtHdrlist1, &updtlistI1, &updtlistS1)
	fmt.Printf("pubop.Publish OK\n")

	time.Sleep(100 * time.Millisecond)

	// Publish a second update
	updthdr4 := &UpdateHeader{*TimeNow(), *mt1_brokerpub_cctx.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
	updthdr5 := &UpdateHeader{*TimeNow(), *mt1_brokerpub_cctx.Uri, MAL_UPDATETYPE_CREATION, *ekpub1}
	updtHdrlist2 := UpdateHeaderList([]*UpdateHeader{updthdr4, updthdr5})

	updtlistI2 := IntegerList([]*Integer{NewInteger(4), NewInteger(5)})
	updtlistS2 := StringList([]*String{NewString("d"), NewString("e")})

	mt1_broker.Publish(&updtHdrlist2, &updtlistI2, &updtlistS2)
	fmt.Printf("pubop.Publish OK\n")

	// Deregisters publisher
	mt1_broker.PublishDeregister()

	fmt.Printf("Publisher#1 end\n")
}

var mt1_subop1 SubscriberOperation
var mt1_sbody1 Body

func newLocalTestM1Sub1() error {
	var err error
	mt1_sub1_ctx, err = NewContext(mt1_subscriber1_url)
	if err != nil {
		return err
	}
	mt1_subscriber1, err = NewClientContext(mt1_sub1_ctx, "subscriber1")
	if err != nil {
		return err
	}
	mt1_subscriber1.SetDomain(IdentifierList([]*Identifier{NewIdentifier("spacecraft1"), NewIdentifier("payload")}))

	mt1_subop1 = mt1_subscriber1.NewSubscriberOperation(mt1_broker.Uri(), 200, 1, 1, 1)
	mt1_sbody1 = mt1_subop1.NewBody()

	domains := IdentifierList([]*Identifier{NewIdentifier("*")})
	eksub := &EntityKey{NewIdentifier("key1"), NewLong(0), NewLong(0), NewLong(0)}
	var erlist = EntityRequestList([]*EntityRequest{
		&EntityRequest{
			&domains, true, true, true, true, EntityKeyList([]*EntityKey{eksub}),
		},
	})
	subs := &Subscription{mt1_subid1, erlist}
	mt1_sbody1.EncodeLastParameter(subs, false)

	mt1_subop1.Register(mt1_sbody1)
	fmt.Printf("subop.Register OK\n")
	// Register is synchronous, we can clear buffer
	mt1_sbody1.Reset(true)

	return nil
}

func runLocalTestM1Sub1(t *testing.T) {
	defer mt1_wg.Done()
	for mt1_running == true {
		// Try to get Notify
		r1, err := mt1_subop1.GetNotify()
		if err != nil {
			fmt.Printf("Subscriber#1, Error in GetNotify: %v\n", err)
			break
		}
		fmt.Printf("\t&&&&& Subscriber1 notified: %d\n", r1.TransactionId)
		mt1_sub1_not_cpt += 1

		id, err := r1.DecodeParameter(NullIdentifier)
		updtHdrlist, err := r1.DecodeParameter(NullUpdateHeaderList)
		updtlistI, err := r1.DecodeParameter(NullIntegerList)
		updtlistS, err := r1.DecodeLastParameter(NullStringList, false)
		mt1_sub1_updt_cpt += len(*updtlistI.(*IntegerList))
		fmt.Printf("\t&&&&& Subscriber1 notified: OK, %s \n\t%+v \n\t%#v\n\t%#v\n\n", id, updtHdrlist, updtlistI, updtlistS)
	}

	if (mt1_sub1_not_cpt != 2) || (mt1_sub1_updt_cpt != 4) {
		t.Errorf("Subscriber#1, bad counters: %d %d", mt1_sub1_not_cpt, mt1_sub1_updt_cpt)
	}

	// Deregisters subscriber

	idlist := IdentifierList([]*Identifier{&mt1_subid1})
	mt1_sbody1.EncodeLastParameter(&idlist, false)
	mt1_subop1.Deregister(mt1_sbody1)
	fmt.Printf("\t&&&&&Subscriber#1, Deregistered\n")
	mt1_sbody1.Reset(true)

	mt1_subscriber1.Close()
	mt1_sub1_ctx.Close()
}

var mt1_subop2 SubscriberOperation
var mt1_sbody2 Body

func newLocalTestM1Sub2() error {
	var err error
	mt1_sub2_ctx, err = NewContext(mt1_subscriber2_url)
	if err != nil {
		return err
	}
	mt1_subscriber2, err = NewClientContext(mt1_sub2_ctx, "subscriber2")
	if err != nil {
		return err
	}
	mt1_subscriber2.SetDomain(IdentifierList([]*Identifier{NewIdentifier("spacecraft1"), NewIdentifier("payload")}))

	mt1_subop2 = mt1_subscriber2.NewSubscriberOperation(mt1_broker.Uri(), 200, 1, 1, 1)
	mt1_sbody2 = mt1_subop2.NewBody()

	domains := IdentifierList([]*Identifier{NewIdentifier("camera1")})
	eksub1 := &EntityKey{NewIdentifier("key1"), NewLong(0), NewLong(0), NewLong(0)}
	eksub2 := &EntityKey{NewIdentifier("key2"), NewLong(0), NewLong(0), NewLong(0)}
	var erlist = EntityRequestList([]*EntityRequest{
		&EntityRequest{
			&domains, true, true, true, true, EntityKeyList([]*EntityKey{eksub1, eksub2}),
		},
	})
	subs := &Subscription{mt1_subid2, erlist}
	mt1_sbody2.EncodeLastParameter(subs, false)

	mt1_subop2.Register(mt1_sbody2)
	fmt.Printf("subop.Register OK\n")
	// Register is synchronous, we can clear buffer
	mt1_sbody2.Reset(true)

	return nil
}

func runLocalTestM1Sub2(t *testing.T) {
	defer mt1_wg.Done()
	for mt1_running == true {
		// Try to get Notify
		r1, err := mt1_subop2.GetNotify()
		if err != nil {
			fmt.Printf("Subscriber#2, Error in GetNotify: %v\n", err)
			break
		}
		fmt.Printf("\t&&&&& Subscriber2 notified: %d\n", r1.TransactionId)
		mt1_sub2_not_cpt += 1

		id, err := r1.DecodeParameter(NullIdentifier)
		updtHdrlist, err := r1.DecodeParameter(NullUpdateHeaderList)
		updtlistI, err := r1.DecodeParameter(NullIntegerList)
		updtlistS, err := r1.DecodeLastParameter(NullStringList, false)
		mt1_sub2_updt_cpt += len(*updtlistI.(*IntegerList))
		fmt.Printf("\t&&&&& Subscriber2 notified: OK, %s \n\t%+v \n\t%#v\n\t%#v\n\n", id, updtHdrlist, updtlistI, updtlistS)
	}

	if (mt1_sub2_not_cpt != 2) || (mt1_sub2_updt_cpt != 5) {
		t.Errorf("Subscriber#2, bad counters: %d %d", mt1_sub2_not_cpt, mt1_sub2_updt_cpt)
	}

	// Deregisters subscriber

	idlist := IdentifierList([]*Identifier{&mt1_subid2})
	mt1_sbody2.EncodeLastParameter(&idlist, false)
	mt1_subop2.Deregister(mt1_sbody2)
	fmt.Printf("\t&&&&&Subscriber#2, Deregistered\n")
	mt1_sbody2.Reset(true)

	mt1_subscriber2.Close()
	mt1_sub2_ctx.Close()
}

func TestM1LocalPubSub(t *testing.T) {
	// Waits socket closing from previous test
	time.Sleep(1000 * time.Millisecond)

	// Creates the broker
	newLocalTestM1BrokerPub(t)
	defer closeLocalTestM1BrokerPub()

	// Creates the subscribers and registers it
	err := newLocalTestM1Sub1()
	if err != nil {
		t.Fatal("Error creating subscriber#1, ", err)
	}
	mt1_wg.Add(1)
	go runLocalTestM1Sub1(t)

	err = newLocalTestM1Sub2()
	if err != nil {
		t.Fatal("Error creating subscriber#2, ", err)
	}
	mt1_wg.Add(1)
	go runLocalTestM1Sub2(t)

	// Waits for subscribers (notify reception)
	time.Sleep(1000 * time.Millisecond)

	// Creates the publishers and registers it
	mt1_wg.Add(1)
	go localTestM1Pub1(t)

	// Waits for subscribers (notify reception)
	time.Sleep(1000 * time.Millisecond)

	fmt.Printf("##### Finish: %d %d\n", mt1_sub1_not_cpt, mt1_sub1_updt_cpt)
	fmt.Printf("##### Finish: %d %d\n", mt1_sub2_not_cpt, mt1_sub2_updt_cpt)

	mt1_subop1.Interrupt()
	mt1_subop2.Interrupt()

	// Wait for subscribers (closing)
	mt1_wg.Wait()
}
