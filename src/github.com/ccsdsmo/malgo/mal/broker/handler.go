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
package broker

import (
	"errors"
	. "github.com/ccsdsmo/malgo/mal"
	. "github.com/ccsdsmo/malgo/mal/api"
	"github.com/ccsdsmo/malgo/mal/debug"
	"github.com/ccsdsmo/malgo/mal/encoding/binary"
)

const (
	varint = true
)

var (
	logger debug.Logger = debug.GetLogger("mal.broker")
)

// Structure used to memorize a subscriber registration
type BrokerSub struct {
	subid       Identifier
	domain      IdentifierList
	session     SessionType
	sessionName Identifier
	serviceArea UShort
	Service     UShort
	operation   UShort
	entities    *EntityRequestList
	transaction SubscriberTransaction
}

func subkey(urifrom string, subid string) string {
	return urifrom + "/" + subid
}

// Structure used to memorize a publisher registration
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

// TODO (AF): Creates a client interface to handle broker implementation

type BrokerImpl struct {
	ctx  *Context
	cctx *ClientContext

	// Map of all active subscribers
	subs map[string]*BrokerSub
	// Map o fall active publishers
	pubs map[string]*BrokerPub
}

func NewBroker(ctx *Context, name string) (*BrokerImpl, error) {
	cctx, err := NewClientContext(ctx, name)
	if err != nil {
		return nil, err
	}

	subs := make(map[string]*BrokerSub)
	pubs := make(map[string]*BrokerPub)
	broker := &BrokerImpl{ctx, cctx, subs, pubs}

	brokerHandler := func(msg *Message, t Transaction) error {
		if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH_REGISTER {
			broker.OnPublishRegister(msg, t.(PublisherTransaction))
		} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH {
			broker.OnPublish(msg, t.(PublisherTransaction))
		} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_PUBLISH_DEREGISTER {
			broker.OnPublishDeregister(msg, t.(PublisherTransaction))
		} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_REGISTER {
			broker.OnRegister(msg, t.(SubscriberTransaction))
		} else if msg.InteractionStage == MAL_IP_STAGE_PUBSUB_DEREGISTER {
			broker.OnDeregister(msg, t.(SubscriberTransaction))
		} else {
			return errors.New("Bad stage")
		}
		return nil
	}
	// Registers the broker handler
	cctx.RegisterBrokerHandler(200, 1, 1, 1, brokerHandler)

	return broker, nil
}

func (handler *BrokerImpl) Uri() *URI {
	return handler.cctx.Uri
}

func (handler *BrokerImpl) Close() {
	// TODO (AF): Removes all remaining subscribers and publishers
	handler.cctx.Close()
}

func (handler *BrokerImpl) register(msg *Message, transaction SubscriberTransaction) error {
	decoder := binary.NewBinaryDecoder(msg.Body, varint)
	sub, err := DecodeSubscription(decoder)
	if err != nil {
		return err
	}
	subkey := subkey(string(*msg.UriFrom), string(sub.SubscriptionId))
	logger.Infof("Broker.Register: %t -> %t", subkey, sub.Entities)

	// Note (AF): Be careful the replacement of a subscription should be an atomic operation.
	handler.subs[subkey] = &BrokerSub{
		subid:       sub.SubscriptionId,
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
	if err != nil {
		// TODO (AF): Builds and encode reply
		return transaction.AckRegister(nil, true)
	} else {
		// TODO (AF): Builds and encode reply
		return transaction.AckRegister(nil, false)
	}
}

func (handler *BrokerImpl) deregister(msg *Message, transaction SubscriberTransaction) error {
	decoder := binary.NewBinaryDecoder(msg.Body, varint)
	list, err := DecodeIdentifierList(decoder)
	if err != nil {
		return err
	}

	for _, id := range []*Identifier(*list) {
		subkey := subkey(string(*msg.UriFrom), string(*id))
		logger.Infof("Broker.Deregister: %v", subkey)
		// TODDO (AF): May be we have to verify if the subscriber is registered.
		delete(handler.subs, string(subkey))
	}
	return nil
}

func (handler *BrokerImpl) OnDeregister(msg *Message, transaction SubscriberTransaction) error {
	err := handler.deregister(msg, transaction)
	if err != nil {
		// TODO (AF): Builds and encode reply
		return transaction.AckDeregister(nil, true)
	} else {
		// TODO (AF): Builds and encode reply
		return transaction.AckDeregister(nil, false)
	}
}

func (handler *BrokerImpl) publishRegister(msg *Message, transaction PublisherTransaction) error {
	decoder := binary.NewBinaryDecoder(msg.Body, varint)
	list, err := DecodeEntityKeyList(decoder)
	if err != nil {
		return err
	}

	logger.Infof("Broker.PublishRegister: %t", list)

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
	if err != nil {
		// TODO (AF): Builds and encode reply
		return transaction.AckRegister(nil, true)
	} else {
		// TODO (AF): Builds and encode reply
		return transaction.AckRegister(nil, false)
	}
}

func (handler *BrokerImpl) publishDeregister(msg *Message, transaction PublisherTransaction) error {
	pubid := string(*msg.UriFrom)
	logger.Infof("Broker.PublishDeregister: %v", pubid)
	// TODDO (AF): May be we have to verify if the publisher is registered.
	delete(handler.pubs, string(pubid))

	return nil
}

func (handler *BrokerImpl) OnPublishDeregister(msg *Message, transaction PublisherTransaction) error {
	err := handler.publishDeregister(msg, transaction)
	if err != nil {
		// TODO (AF): Builds and encode reply
		return transaction.AckDeregister(nil, true)
	} else {
		// TODO (AF): Builds and encode reply
		return transaction.AckDeregister(nil, false)
	}
}

func (handler *BrokerImpl) publish(pub *Message, transaction PublisherTransaction) error {
	logger.Debugf("Broker.Publish -> %v", pub)

	decoder := binary.NewBinaryDecoder(pub.Body, varint)
	uhlist, err := DecodeUpdateHeaderList(decoder)
	if err != nil {
		return err
	}
	logger.Infof("Broker.Publish, DecodeUpdateHeaderList -> %+v", uhlist)
	uvlist, err := DecodeUpdateList(decoder)
	if err != nil {
		return err
	}
	logger.Infof("Broker.Publish, DecodeUpdateList -> %d %v", len([]*Blob(*uvlist)), uvlist)

	for _, sub := range handler.subs {
		buf := make([]byte, 0, 1024)
		encoder := binary.NewBinaryEncoder(buf, varint)
		encoder.EncodeIdentifier(&sub.subid)
		// encoder.EncodeElement(uhlist)
		// Encode update value list
		//		NewUInteger(uint32(len([]*Blob(*uvlist)))).Encode(encoder)
		//		for _, uv := range []*Blob(*uvlist) {
		//			encoder.WriteBody([]byte(*uv))
		//		}
		uhlist.Encode(encoder)
		uvlist.Encode(encoder)
		sub.transaction.Notify(encoder.Body(), false)
	}
	return nil
}

func (handler *BrokerImpl) OnPublish(msg *Message, transaction PublisherTransaction) error {
	err := handler.publish(msg, transaction)
	if err != nil {
		// TODO (AF): Returns error
		//		return transaction.PublishError(err)
		return err
	}
	return nil
}
