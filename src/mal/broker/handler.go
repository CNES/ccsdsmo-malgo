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
package broker

import (
	. "mal"
	. "mal/api2"
	"mal/encoding/binary"
)

type BrokerSub struct {
	domain      IdentifierList
	session     SessionType
	sessionName Identifier
	serviceArea UShort
	Service     UShort
	operation   UShort
	entities    *EntityRequestList
	transaction SubscriberTransaction
}

type BrokerPub struct {
	domain      IdentifierList
	session     SessionType
	sessionName Identifier
	serviceArea UShort
	Service     UShort
	operation   UShort
	keys        *EntityKeyList
	transaction PublisherTransaction
}

type BrokerImpl struct {
	subs map[string]*BrokerSub
	pubs map[string]*BrokerPub
}

func (handler *BrokerImpl) register(msg *Message, transaction SubscriberTransaction) error {
	decoder := binary.NewBinaryDecoder(msg.Body)
	sub, err := DecodeSubscription(decoder)
	if err != nil {
		return err
	}

	subid := string(*msg.UriFrom) + string(sub.SubscriptionId)
	// Note (AF): Be careful the replacement of a subscription is an atomic operation.
	handler.subs[subid] = &BrokerSub{
		domain:      msg.Domain,
		session:     msg.Session,
		sessionName: msg.SessionName,
		serviceArea: msg.ServiceArea,
		Service:     msg.Service,
		operation:   msg.Operation,
		entities:    &sub.Entities,
		transaction: transaction,
	}

	return nil
}

func (handler *BrokerImpl) OnRegister(msg *Message, transaction SubscriberTransaction) error {
	err := handler.register(msg, transaction)
	return transaction.RegisterAck(err)
}

func (handler *BrokerImpl) deregister(msg *Message, transaction SubscriberTransaction) error {
	decoder := binary.NewBinaryDecoder(msg.Body)
	list, err := DecodeIdentifierList(decoder)
	if err != nil {
		return err
	}

	for _, id := range []*Identifier(*list) {
		subid := string(*msg.UriFrom) + string(*id)
		// TODDO (AF): May be we have to verify if the subscription exists.
		delete(handler.subs, string(subid))
	}
	return nil
}

func (handler *BrokerImpl) OnDeregister(msg *Message, transaction SubscriberTransaction) error {
	err := handler.deregister(msg, transaction)
	return transaction.DeregisterAck(err)
}

func (handler *BrokerImpl) publishRegister(msg *Message, transaction PublisherTransaction) error {
	decoder := binary.NewBinaryDecoder(msg.Body)
	list, err := DecodeEntityKeyList(decoder)
	if err != nil {
		return err
	}

	pubid := string(*msg.UriFrom)
	handler.pubs[pubid] = &BrokerPub{
		domain:      msg.Domain,
		session:     msg.Session,
		sessionName: msg.SessionName,
		serviceArea: msg.ServiceArea,
		Service:     msg.Service,
		operation:   msg.Operation,
		keys:        list,
		transaction: transaction,
	}

	return nil
}

func (handler *BrokerImpl) OnPublishRegister(msg *Message, transaction PublisherTransaction) error {
	err := handler.publishRegister(msg, transaction)
	return transaction.PublishRegisterAck(err)
}

func (handler *BrokerImpl) publishDeregister(msg *Message, transaction PublisherTransaction) error {
	pubid := string(*msg.UriFrom)
	// TODDO (AF): May be we have to verify if the publisher exists.
	delete(handler.pubs, string(pubid))

	return nil
}

func (handler *BrokerImpl) OnPublishDeregister(msg *Message, transaction PublisherTransaction) error {
	err := handler.publishDeregister(msg, transaction)
	return transaction.PublishDeregisterAck(err)
}

func (handler *BrokerImpl) publish(pub *Message, transaction PublisherTransaction) error {
	decoder := binary.NewBinaryDecoder(pub.Body)
	uhlist, err := DecodeUpdateHeaderList(decoder)
	if err != nil {
		return err
	}
	uvlist, err := DecodeUpdateList(decoder)
	if err != nil {
		return err
	}

	for id, sub := range handler.subs {
		buf := make([]byte, 0, 1024)
		encoder := binary.NewBinaryEncoder(buf)
		encoder.EncodeIdentifier(NewIdentifier(id))
		encoder.EncodeElement(uhlist)
		NewUInteger(uint32(len([]*Blob(*uvlist)))).Encode(encoder)
		for _, uv := range []*Blob(*uvlist) {
			encoder.WriteBody([]byte(*uv))
		}
		sub.transaction.Notify(encoder.Buffer(), nil)
	}
	return nil
}

func (handler *BrokerImpl) OnPublish(msg *Message, transaction PublisherTransaction) error {
	err := handler.publish(msg, transaction)
	if err != nil {
		return transaction.PublishError(err)
	}
	return nil
}
