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
	// create a body for the operation call
	body := provider.pub.NewBody()
	err := body.EncodeLastParameter(keys, false)
	if err != nil {
		logger.Errorf("EventProvider.monitorEventRegister error: %v", err)
		return err
	}
	msg, err := provider.pub.Register(body)
	if err != nil {
		handleError("EventProvider.monitorEventRegister error", err, msg)
		return err
	}
	return nil
}

func (provider *EventProvider) monitorEventPublish(updtHdr *UpdateHeaderList, updtDetails *ObjectDetailsList, updtValues ElementList) error {
	// create a body for the operation call
	body := provider.pub.NewBody()
	err := body.EncodeParameter(updtHdr)
	if err != nil {
		logger.Errorf("EventProvider.monitorEventPublish error: %v", err)
		return err
	}
	err = body.EncodeParameter(updtDetails)
	if err != nil {
		logger.Errorf("EventProvider.monitorEventPublish error: %v", err)
		return err
	}
	err = body.EncodeLastParameter(updtValues, true)
	if err != nil {
		logger.Errorf("EventProvider.monitorEventPublish error: %v", err)
		return err
	}
	err = provider.pub.Publish(body)
	if err != nil {
		logger.Errorf("EventProvider.monitorEventPublish error: %v", err)
		return err
	}
	return nil
}

func (provider *EventProvider) monitorEventDeregister() error {
	msg, err := provider.pub.Deregister(nil)
	if err != nil {
		handleError("EventProvider.monitorEventDeregister error", err, msg)
		return err
	}
	return nil
}

func (provider *EventProvider) Close() error {
	return provider.cctx.Close()
}
