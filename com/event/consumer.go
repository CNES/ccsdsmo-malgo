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
package event

import (
	. "github.com/CNES/ccsdsmo-malgo/com"
	. "github.com/CNES/ccsdsmo-malgo/mal"
	. "github.com/CNES/ccsdsmo-malgo/mal/api"
)

type EventConsumer struct {
	factory EncodingFactory
	cctx    *ClientContext
	subs    SubscriberOperation
}

func NewEventConsumer(factory EncodingFactory, cctx *ClientContext, broker *URI) (*EventConsumer, error) {
	subs := cctx.NewSubscriberOperation(broker,
		COM_AREA_NUMBER, COM_AREA_VERSION,
		COM_EVENT_SERVICE_NUMBER, COM_EVENT_MONITOR_EVENT_OP_NUMBER)
	consumer := &EventConsumer{factory: factory, cctx: cctx, subs: subs}
	return consumer, nil
}

func (consumer *EventConsumer) monitorEventRegister(sub Subscription) error {
	encoder := consumer.factory.NewEncoder(make([]byte, 0, 8192))
	sub.Encode(encoder)
	msg, err := consumer.subs.Register(encoder.Body())
	if err != nil {
		// TODO (AF): Get errors in reply then log a message.
		logger.Errorf("EventConsumer.monitorEventRegister error: %v, %v", msg, err)
		return err
	}
	return nil
}

// Get events.
// Last out parameter of type ElementList should be cast in *XList using out.(*XList)) where X is the type
// of wanted elements.
func (consumer *EventConsumer) monitorEventGetNotify() (*Message, *Identifier, *UpdateHeaderList, *ObjectDetails, ElementList, error) {
	msg, err := consumer.subs.GetNotify()
	if err != nil {
		// TODO (AF): Get errors in reply then log a message.
		return nil, nil, nil, nil, nil, err
	}

	decoder := consumer.factory.NewDecoder(msg.Body)
	id, err := decoder.DecodeIdentifier()
	updtHdrlist, err := DecodeUpdateHeaderList(decoder)
	updtDetailslist, err := DecodeObjectDetails(decoder)
	updtElementlist, err := decoder.DecodeAbstractElement()

	logger.Debugf("EventConsumer.monitorEventGetNotify: %s, %v, %v, %v", id, updtHdrlist, updtDetailslist, updtElementlist)

	return msg, id, updtHdrlist, updtDetailslist, updtElementlist.(ElementList), nil
}

func (consumer *EventConsumer) monitorEventDeregister(subids IdentifierList) error {
	encoder := consumer.factory.NewEncoder(make([]byte, 0, 8192))
	subids.Encode(encoder)
	msg, err := consumer.subs.Deregister(encoder.Body())
	if err != nil {
		// TODO (AF): Get error in reply then log a message.
		logger.Errorf("EventConsumer.monitorEventDeregister error: %v, %v", msg, err)
		return err
	}
	return nil
}

func (consumer *EventConsumer) Close() error {
	return consumer.cctx.Close()
}
