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
	// create a body for the operation call
	body := consumer.subs.NewBody()

	err := body.EncodeLastParameter(&sub, false)
	if err != nil {
		logger.Errorf("EventConsumer.monitorEventRegister error: %v", err)
		return err
	}

	msg, err := consumer.subs.Register(body)
	if err != nil {
		handleError("EventConsumer.monitorEventRegister error", err, msg)
		return err
	}
	return nil
}

// Get events.
// Last out parameter of type ElementList should be cast in *XList using out.(*XList)) where X is the type
// of wanted elements.
func (consumer *EventConsumer) monitorEventGetNotify() (*Message, *Identifier, *UpdateHeaderList, *ObjectDetailsList, ElementList, error) {
	msg, err := consumer.subs.GetNotify()
	if err != nil {
		handleError("EventConsumer.monitorEventGetNotify error", err, msg)
		return nil, nil, nil, nil, nil, err
	}

	param, err := msg.DecodeParameter(NullIdentifier)
	if err != nil {
		logger.Errorf("EventConsumer.monitorEventGetNotify error: decoding Identifier, %v", err)
		return nil, nil, nil, nil, nil, err
	}
	id := param.(*Identifier)
	param, err = msg.DecodeParameter(NullUpdateHeaderList)
	if err != nil {
		logger.Errorf("EventConsumer.monitorEventGetNotify error: decoding UpdateHeaderList, %v", err)
		return nil, nil, nil, nil, nil, err
	}
	updtHdrlist := param.(*UpdateHeaderList)
	param, err = msg.DecodeParameter(NullObjectDetailsList)
	if err != nil {
		logger.Errorf("EventConsumer.monitorEventGetNotify error: decoding ObjectDetailsList, %v", err)
		return nil, nil, nil, nil, nil, err
	}
	updtDetailslist := param.(*ObjectDetailsList)
	param, err = msg.DecodeLastParameter(nil, true)
	if err != nil {
		logger.Errorf("EventConsumer.monitorEventGetNotify error: decoding ElementList, %v", err)
		return nil, nil, nil, nil, nil, err
	}
	updtElementlist := param.(ElementList)

	logger.Debugf("EventConsumer.monitorEventGetNotify: %s, %v, %v, %v", id, updtHdrlist, updtDetailslist, updtElementlist)

	return msg, id, updtHdrlist, updtDetailslist, updtElementlist, nil
}

func (consumer *EventConsumer) monitorEventDeregister(subids IdentifierList) error {
	// create a body for the operation call
	body := consumer.subs.NewBody()

	err := body.EncodeLastParameter(&subids, false)
	if err != nil {
		logger.Errorf("EventConsumer.monitorEventDeregister error: %v", err)
		return err
	}

	msg, err := consumer.subs.Deregister(body)
	if err != nil {
		handleError("EventConsumer.monitorEventDeregister error", err, msg)
		return err
	}

	return nil
}

func (consumer *EventConsumer) Close() error {
	var err error = nil
	if consumer.subs != nil {
		err = consumer.subs.Close()
	}
	return err
}
