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
package mal

import (
	"errors"
	"github.com/CNES/ccsdsmo-malgo/mal/debug"
)

var rlogger debug.Logger = debug.GetLogger("mal.registry")
var _MALElementRegistry = map[int64]Element{}

var NullElement Element = nil

func RegisterMALElement(shortForm Long, element Element) error {
	rlogger.Debugf("MALElementRegistry.RegisterMALElement: %x", (int64)(shortForm))
	_, ok := _MALElementRegistry[(int64)(shortForm)]
	if ok {
		rlogger.Errorf("MALElementRegistry.RegisterMALElement: %x already registered", (int64)(shortForm))
		return errors.New("MALElementRegistry.RegisterMALElement: element already registered")
	}
	_MALElementRegistry[(int64)(shortForm)] = element
	return nil
}

func LookupMALElement(shortForm Long) (Element, error) {
	rlogger.Debugf("MALElementRegistry.LookupMALElement: %x", (int64)(shortForm))
	element, ok := _MALElementRegistry[(int64)(shortForm)]
	if !ok {
		rlogger.Errorf("MALElementRegistry.LookupMALElement: unknown %x element", (int64)(shortForm))
		return nil, errors.New("MALElementRegistry.LookupMALElement: unknown element")
	}
	return element, nil
}

func DeregisterMALElement(shortForm Long) error {
	rlogger.Debugf("MALElementRegistry.DeregisterMALElement: %x", (int64)(shortForm))
	_, ok := _MALElementRegistry[(int64)(shortForm)]
	if !ok {
		rlogger.Errorf("MALElementRegistry.DeregisterMALElement: %x not registered", (int64)(shortForm))
		return errors.New("MALElementRegistry.DeregisterMALElement: element not registered")
	}
	delete(_MALElementRegistry, (int64)(shortForm))
	return nil
}

// The Element interface represents the MAL Element type.
type Element interface {
	// Returns the absolute short form of the element type.
	GetShortForm() Long

	// Returns the number of the area this element type belongs to.
	GetAreaNumber() UShort

	// Returns the version of the area this element type belongs to.
	GetAreaVersion() UOctet

	// Returns the number of the service this element type belongs to.
	GetServiceNumber() UShort

	// Return the relative short form of the element type.
	GetTypeShortForm() Integer

	// Encodes this element using the supplied encoder.
	// @param encoder The encoder to use, must not be null.
	Encode(encoder Encoder) error

	// Decodes an instance of this element type using the supplied decoder.
	// @param decoder The decoder to use, must not be null.
	// @return the decoded instance, may be not the same instance as this Element.
	Decode(decoder Decoder) (Element, error)

	// TODO (AF): It seems that this method is no longer needed (to verify)
	// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
	CreateElement() Element

	IsNull() bool
	Null() Element
}

func GetListShortForm(element Element) Long {
	if element == nil {
		return -1
	}
	shortForm := element.GetShortForm()
	// this short form should never represent a list type
	typePart := ULong(shortForm & 0xFFFFFF)
	if typePart > 0x7FFFFF {
		return -1
	}
	return Long((ULong(shortForm) & 0xFFFFFFFFFF000000) | ULong((-Long(typePart))&0xFFFFFF))
}

func GetListItemShortForm(element Element) Long {
	if element == nil {
		return -1
	}
	shortForm := element.GetShortForm()
	// this short form should represent a list type
	typePart := ULong(shortForm & 0xFFFFFF)
	if typePart < 0x800000 {
		return -1
	}
	return Long((ULong(shortForm) & 0xFFFFFFFFFF000000) | ULong((-Long(typePart))&0xFFFFFF))
}

func CreateElement(shortform Long) Element {
	null, err := LookupMALElement(shortform)
	if err != nil {
		return nil
	}
	return null.CreateElement()
}
