/**
 * MIT License
 *
 * Copyright (c) 2017 - 2019 CNES
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
	. "github.com/CNES/ccsdsmo-malgo/mal"
	. "github.com/CNES/ccsdsmo-malgo/mal/api"
	"github.com/CNES/ccsdsmo-malgo/mal/debug"
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
	service     UShort
	operation   UShort
	entities    *EntityRequestList
	transaction SubscriberTransaction
}

func subkey(urifrom string, subid string) string {
	// Conforming to 3.5.6.3.d, the URI of the consumer and the subscription identifier
	// shall form the unique identifier of the subscription.
	return urifrom + "/" + subid
}

func (sub *BrokerSub) domainMatches(domain IdentifierList, subdomain *IdentifierList) bool {
	// See MAL specification 3.5.6.5 e,f,g p 3-57

	logger.Debugf("Broker.domainMatches [%v + %v] -> %v", sub.domain, subdomain, domain)

	// e) The domain of the update message shall match the domain of the subscription message.
	// f) If the subscription EntityRequest included a subDomain field, then this shall be appended
	//    to the domain of the subscription message to make the complete domain for that request.
	// g) The final Identifier of the subDomain may be the wildcard character ‘*’.

	var required []*Identifier
	var all bool = false

	if subdomain == nil {
		required = sub.domain
	} else {
		required = make([]*Identifier, 0, len(sub.domain)+len(*subdomain))
		required = append(required, sub.domain...)
		required = append(required, *subdomain...)
		if (*(required)[len(required)-1]) == "*" {
			all = true
			required = required[:len(required)-1]
		}
	}
	logger.Debugf("Broker.domainMatches %v, %v", required, all)
	if len(domain) < len(required) {
		logger.Debugf("Broker.domainMatches #1 !NOK! -> %d < %d", len(domain), len(required))
		return false
	}

	for idx, name := range ([]*Identifier)(required) {
		if *name != *([]*Identifier)(domain)[idx] {
			logger.Debugf("Broker.domainMatches #2 %d %s != %s !NOK!", idx, *name, *([]*Identifier)(domain)[idx])
			return false
		}
	}

	if len(domain) > len(required) {
		logger.Debugf("Broker.domainMatches #3 !NOK! -> %v", all)
		return all
	}

	return true
}

func (sub *BrokerSub) matches(msg *Message, key *EntityKey) bool {
	// See MAL specification 3.5.6.5 e,f,g p 3-57
	logger.Debugf("Broker.matches? -> %s", sub.subid)

	if (msg.Session != sub.session) || (msg.SessionName != sub.sessionName) {
		// h) The session types and names must match.
		logger.Debugf("Broker.matches #1 !NOK!")
		return false
	}

	// Evaluates all requests of the subscription
	for _, request := range ([]*EntityRequest)(*sub.entities) {
		if !sub.domainMatches(msg.Domain, request.SubDomain) {
			logger.Debugf("Broker.matches #2 !NOK!")
			continue
		}
		if !request.AllAreas && msg.ServiceArea != sub.serviceArea {
			// j) The area identifiers must match unless the subscription specified True in the allAreas
			//    field of the EntityRequest, in which case they shall be ignored.
			logger.Debugf("Broker.matches #3 !NOK!")
			continue
		}
		if !request.AllServices && msg.Service != sub.service {
			// k) The service identifiers must match unless the subscription specified True in the
			//    allServices field of the EntityRequest, in which case they shall be ignored.
			logger.Debugf("Broker.matches #4 !NOK!")
			continue
		}
		if !request.AllOperations && msg.Operation != sub.operation {
			// l) The operation identifiers must match unless the subscription specified True in the
			// allOperations field of the EntityRequest, in which case they shall be ignored.
			logger.Debugf("Broker.matches #5 !NOK!")
			continue
		}

		// Search for a matching key in the current request
		for _, rkey := range ([]*EntityKey)(request.EntityKeys) {
			// a) A sub-key specified in the EntityKey structure shall take one of three types of value:
			//    an actual value, a NULL value, and the special wildcard value of ‘*’ (for the first subkey
			//    only) or zero (for the other three sub-keys).
			// b) If a sub-key contains a specific value it shall only match a sub-key that contains the
			//    same value. This includes an empty ‘’ value for the first sub-key. The matches are
			//    case sensitive.
			// c) If a sub-key contains a NULL value it shall only match a sub-key that contains
			//    NULL.
			// d) If a sub-key contains the wildcard value it shall match a sub-key that contains any
			//    value including NULL.

			logger.Debugf("Broker.matches request -> %s", *rkey)
			logger.Debugf("Broker.matches update -> %s", key)

			if rkey.Match(key) {
				logger.Debugf("Broker.matches #6 OK")
				return true
			}
		}
		// There is no matching key in this entity request
	}

	// There is no matching key in this subscription
	return false
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
	// TODO (AF): Is it needed ? used ?
	transaction PublisherTransaction
}

// TODO (AF): Creates a client interface to handle broker implementation

type BrokerHandler struct {
	cctx *ClientContext

	updtHandler UpdateValueHandler

	encoding EncodingFactory

	// Map of all active subscribers
	subs map[string]*BrokerSub
	// Map o fall active publishers
	pubs map[string]*BrokerPub
}

type UpdateValueHandler interface {
	DecodeUpdateValueList(body Body) error
	UpdateValueListSize() int
	AppendValue(idx int)
	EncodeUpdateValueList(body Body) error
	ResetValues()
}

// ################################################################################
// Implements an UpdateValueHandler for Blob update value type

type BlobUpdateValueHandler struct {
	list   *BlobList
	values BlobList
}

func NewBlobUpdateValueHandler() *BlobUpdateValueHandler {
	return new(BlobUpdateValueHandler)
}

func (handler *BlobUpdateValueHandler) DecodeUpdateValueList(body Body) error {
	p, err := body.DecodeLastParameter(NullBlobList, false)
	//	list, err := DecodeBlobList(decoder)
	if err != nil {
		return err
	}
	list := p.(*BlobList)
	logger.Infof("Broker.Publish, DecodeUpdateValueList -> %d %v", len([]*Blob(*list)), list)

	handler.list = list
	handler.values = BlobList(make([]*Blob, 0, handler.list.Size()))

	return nil
}

func (handler *BlobUpdateValueHandler) UpdateValueListSize() int {
	return handler.list.Size()
}

func (handler *BlobUpdateValueHandler) AppendValue(idx int) {
	handler.values = append(handler.values, ([]*Blob)(*handler.list)[idx])
}

func (handler *BlobUpdateValueHandler) EncodeUpdateValueList(body Body) error {
	//	err := handler.values.Encode(encoder)
	err := body.EncodeLastParameter(&handler.values, false)
	if err != nil {
		return err
	}
	handler.values = handler.values[:0]
	return nil
}

func (handler *BlobUpdateValueHandler) ResetValues() {
	handler.values = handler.values[:0]
}

// ################################################################################

func NewBroker(cctx *ClientContext, updtHandler UpdateValueHandler, encoding EncodingFactory) (*BrokerHandler, error) {
	subs := make(map[string]*BrokerSub)
	pubs := make(map[string]*BrokerPub)
	broker := &BrokerHandler{cctx, updtHandler, encoding, subs, pubs}

	brokerHandler := func(msg *Message, t Transaction) error {
		//		fmt.Println("##########", msg.Body)
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

func (handler *BrokerHandler) Uri() *URI {
	return handler.cctx.Uri
}

// Gets the underlying ClientContext used by the broker.
func (handler *BrokerHandler) ClientContext() *ClientContext {
	return handler.cctx
}

func (handler *BrokerHandler) Close() {
	// TODO (AF): Removes all remaining subscribers and publishers
	handler.cctx.Close()
}

func (handler *BrokerHandler) register(msg *Message, transaction SubscriberTransaction) error {
	p, err := msg.DecodeLastParameter(NullSubscription, false)
	//	decoder := handler.encoding.NewDecoder(msg.Body)
	//	sub, err := DecodeSubscription(decoder)
	if err != nil {
		return err
	}
	sub := p.(*Subscription)

	subkey := subkey(string(*msg.UriFrom), string(sub.SubscriptionId))
	logger.Infof("Broker.Register: %t -> %t", subkey, sub.Entities)

	// Note (AF): Be careful the replacement of a subscription should be an atomic operation.
	handler.subs[subkey] = &BrokerSub{
		subid:       sub.SubscriptionId,
		domain:      msg.Domain,
		session:     msg.Session,
		sessionName: msg.SessionName,
		serviceArea: msg.ServiceArea,
		service:     msg.Service,
		operation:   msg.Operation,
		entities:    &sub.Entities,
		transaction: transaction,
	}

	return nil
}

func (handler *BrokerHandler) OnRegister(msg *Message, transaction SubscriberTransaction) error {
	err := handler.register(msg, transaction)
	if err != nil {
		return transaction.AckRegister(nil, true)
	} else {
		// TODO (AF): Builds and encode error structure (cf. 3.5.6.11.3)
		return transaction.AckRegister(nil, false)
	}
}

func (handler *BrokerHandler) deregister(msg *Message, transaction SubscriberTransaction) error {
	p, err := msg.DecodeLastParameter(NullIdentifierList, false)
	//	decoder := handler.encoding.NewDecoder(msg.Body)
	//	list, err := DecodeIdentifierList(decoder)
	if err != nil {
		return err
	}
	list := p.(*IdentifierList)

	for _, id := range []*Identifier(*list) {
		subkey := subkey(string(*msg.UriFrom), string(*id))
		logger.Infof("Broker.Deregister: %v", subkey)
		// TODDO (AF): May be we have to verify if the subscriber is registered.
		delete(handler.subs, string(subkey))
	}
	return nil
}

func (handler *BrokerHandler) OnDeregister(msg *Message, transaction SubscriberTransaction) error {
	err := handler.deregister(msg, transaction)
	if err == nil {
		// TODO (AF): Logs an error message
	}
	return transaction.AckDeregister(nil, true)
}

func (handler *BrokerHandler) publishRegister(msg *Message, transaction PublisherTransaction) error {
	p, err := msg.DecodeLastParameter(NullEntityKeyList, false)
	//	decoder := handler.encoding.NewDecoder(msg.Body)
	//	list, err := DecodeEntityKeyList(decoder)
	if err != nil {
		return err
	}
	list := p.(*EntityKeyList)

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

func (handler *BrokerHandler) OnPublishRegister(msg *Message, transaction PublisherTransaction) error {
	err := handler.publishRegister(msg, transaction)
	if err != nil {
		// TODO (AF): Builds and encode reply
		return transaction.AckRegister(nil, true)
	} else {
		// TODO (AF): Builds and encode error structure (cf 3.5.6.11.6)
		return transaction.AckRegister(nil, false)
	}
}

func (handler *BrokerHandler) publishDeregister(msg *Message, transaction PublisherTransaction) error {
	pubid := string(*msg.UriFrom)
	logger.Infof("Broker.PublishDeregister: %v", pubid)
	// TODDO (AF): May be we have to verify if the publisher is registered.
	delete(handler.pubs, string(pubid))

	return nil
}

func (handler *BrokerHandler) OnPublishDeregister(msg *Message, transaction PublisherTransaction) error {
	err := handler.publishDeregister(msg, transaction)
	if err == nil {
		// TODO (AF): Logs an error message
	}
	return transaction.AckDeregister(nil, true)
}

func (handler *BrokerHandler) publish(pub *Message, transaction PublisherTransaction) error {
	logger.Debugf("Broker.Publish -> %v", pub)

	p1, err := pub.Body.DecodeParameter(NullUpdateHeaderList)
	//	decoder := handler.encoding.NewDecoder(pub.Body)
	//	uhlist, err := DecodeUpdateHeaderList(decoder)
	if err != nil {
		return err
	}
	uhlist := p1.(*UpdateHeaderList)
	logger.Infof("Broker.Publish, DecodeUpdateHeaderList -> %+v", uhlist)

	handler.updtHandler.DecodeUpdateValueList(pub.Body)
	if err != nil {
		// TODO (AF): Returns a PublishError MAL message to publisher
		return err
	}
	logger.Infof("Broker.Publish, DecodeUpdateList -> %d", handler.updtHandler.UpdateValueListSize())

	if len(*uhlist) != handler.updtHandler.UpdateValueListSize() {
		return errors.New("Bad header and value list lengths")
	}

	// TODO (AF): Verify the publication validity see 3.5.6.8 e, f

	for _, sub := range handler.subs {
		var headers UpdateHeaderList = make([]*UpdateHeader, 0, len(*uhlist))
		for idx, hdr := range *uhlist {
			if sub.matches(pub, &hdr.Key) {
				logger.Warnf("Broker.Publish match !!")
				// Adds the update to the notify message for this subscription
				headers = append(headers, hdr)
				handler.updtHandler.AppendValue(idx)
			}
		}
		if len(headers) == 0 {
			// there is no update matching this subscription
			handler.updtHandler.ResetValues()
			continue
		}

		body := transaction.NewBody()
		//		buf := make([][]byte, 1)
		//		buf[0] = make([]byte, 0, 1024)
		//		encoder := handler.encoding.NewEncoder(buf)
		//		encoder.EncodeIdentifier(&sub.subid)
		body.EncodeParameter(&sub.subid)
		//		headers.Encode(encoder)
		body.EncodeParameter(&headers)
		handler.updtHandler.EncodeUpdateValueList(body)
		//		sub.transaction.Notify(encoder.Body(), false)
		sub.transaction.Notify(body, false)
	}
	return nil
}

func (handler *BrokerHandler) OnPublish(msg *Message, transaction PublisherTransaction) error {
	err := handler.publish(msg, transaction)
	if err != nil {
		// TODO (AF): Returns error
		//		return transaction.PublishError(err)
		return err
	}
	return nil
}
