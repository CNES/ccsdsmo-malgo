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

type EventProvider struct {
	factory EncodingFactory
	cctx    *ClientContext
	pub     PublisherOperation
}

// Creates a publisher object using the current registered provider set for the PubSub operation monitorEvent.
func NewEventProvider(factory EncodingFactory, cctx *ClientContext, broker *URI) (*EventProvider, error) {
	pub := cctx.NewPublisherOperation(broker,
		COM_AREA_NUMBER, COM_AREA_VERSION,
		COM_EVENT_SERVICE_NUMBER, COM_EVENT_MONITOR_EVENT_OP_NUMBER)
	provider := &EventProvider{factory: factory, cctx: cctx, pub: pub}
	return provider, nil
}

func (provider *EventProvider) monitorEventRegister(keys *EntityKeyList) error {
	encoder := provider.factory.NewEncoder(make([]byte, 0, 8192))
	keys.Encode(encoder)
	msg, err := provider.pub.Register(encoder.Body())
	if err != nil {
		// TODO (AF): Get errors in reply then log a message.
		logger.Errorf("EventProvider.monitorEventRegister error: %v, %v", msg, err)
		return err
	}
	return nil
}

func (provider *EventProvider) monitorEventPublish(updtHdr *UpdateHeaderList, updtDetails *ObjectDetailsList, updtValues ElementList) error {
	encoder := provider.factory.NewEncoder(make([]byte, 0, 8192))
	updtHdr.Encode(encoder)
	updtDetails.Encode(encoder)
	encoder.EncodeAbstractElement(updtValues)
	msg, err := provider.pub.Register(encoder.Body())
	if err != nil {
		// TODO (AF): Get errors in reply then log a message.
		logger.Errorf("EventProvider.monitorEventPublish error: %v, %v", msg, err)
		return err
	}
	return nil
}

func (provider *EventProvider) monitorEventDeregister() error {
	msg, err := provider.pub.Register(nil)
	if err != nil {
		// TODO (AF): Get errors in reply then log a message.
		logger.Errorf("EventProvider.monitorEventDeregister error: %v, %v", msg, err)
		return err
	}
	return nil
}

func (provider *EventProvider) Close() error {
	return provider.cctx.Close()
}
