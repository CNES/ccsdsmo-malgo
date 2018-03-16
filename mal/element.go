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
package mal

import ()

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
