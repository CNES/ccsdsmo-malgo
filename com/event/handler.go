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
	"errors"
	. "github.com/CNES/ccsdsmo-malgo/com"
	. "github.com/CNES/ccsdsmo-malgo/mal"
	. "github.com/CNES/ccsdsmo-malgo/mal/api"
	. "github.com/CNES/ccsdsmo-malgo/mal/broker"
	"github.com/CNES/ccsdsmo-malgo/mal/debug"
)

var (
	logger debug.Logger = debug.GetLogger("com.event")
)

type EventHandler struct {
	cctx   *ClientContext
	broker *BrokerHandler
}

func NewEventHandler(factory EncodingFactory, cctx *ClientContext) (*EventHandler, error) {
	updtHandler := NewEventUpdateValueHandler()
	broker, err := NewBroker(cctx, updtHandler, factory)
	if err != nil {
		return nil, err
	}
	handler := &EventHandler{cctx, broker}
	return handler, nil
}

func (handler *EventHandler) Close() error {
	handler.broker.Close()
	return handler.cctx.Close()
}

// ################################################################################
// Implements UpdateValueHandler for Event specific broker.
// Handles updates containing an ObjectDetail and an Element depending of the event.

type EventUpdateValueHandler struct {
	list1   *ObjectDetailsList
	values1 ObjectDetailsList
	list2   []Element
	values2 []Element
}

func NewEventUpdateValueHandler() *EventUpdateValueHandler {
	return new(EventUpdateValueHandler)
}

func (handler *EventUpdateValueHandler) DecodeUpdateValueList(decoder Decoder) error {
	list1, err := DecodeObjectDetailsList(decoder)
	if err != nil {
		return err
	}
	logger.Infof("Broker.Publish, DecodeUpdateValueList -> %d, %v", len([]*ObjectDetails(*list1)), list1)

	list2, err := decoder.DecodeElementList()
	if err != nil {
		return err
	}
	logger.Infof("Broker.Publish, DecodeUpdateValueList -> %d, %v", len(list2), list2)

	if len([]*ObjectDetails(*list1)) != len(list2) {
		return errors.New("Bad list length")
	}

	handler.list1 = list1
	handler.values1 = ObjectDetailsList(make([]*ObjectDetails, 0, len([]*ObjectDetails(*list1))))

	handler.list2 = list2
	handler.values2 = make([]Element, 0, len(list2))

	return nil
}

func (handler *EventUpdateValueHandler) UpdateValueListSize() int {
	return len([]*ObjectDetails(*handler.list1))
}

func (handler *EventUpdateValueHandler) AppendValue(idx int) {
	handler.values1 = append(handler.values1, ([]*ObjectDetails)(*handler.list1)[idx])
	handler.values2 = append(handler.values2, handler.list2[idx])
}

func (handler *EventUpdateValueHandler) EncodeUpdateValueList(encoder Encoder) error {
	err := encoder.EncodeElement(&handler.values1)
	if err != nil {
		return err
	}
	handler.values1 = handler.values1[:0]

	err = encoder.EncodeElementList(handler.values2)
	if err != nil {
		return err
	}
	handler.values2 = handler.values2[:0]

	return nil
}

func (handler *EventUpdateValueHandler) ResetValues() {
	handler.values1 = handler.values1[:0]
	handler.values2 = handler.values2[:0]
}
