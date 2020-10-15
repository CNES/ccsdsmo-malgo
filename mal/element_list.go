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
package mal

// ################################################################################
// Defines MAL ElementList type
// It corresponds to a list of abstract elements.
// It should be implemented by all list of elements.
// ################################################################################

type ElementList interface {
	Composite
	Size() int
	GetElementAt(i int) Element
	AppendElement(element Element)
}

// ################################################################################
// Defines a MAL ElementList generic implementation
// Be careful, ElementList is not a MAL type: It does not implement Element and
// there is no corresponding short-form.
// ################################################################################

type ElementListImpl []Element

var (
	NullElementList *ElementListImpl = nil
)

func NewElementListImpl(size int) *ElementListImpl {
	var list ElementListImpl = ElementListImpl(make([]Element, size))
	return &list
}

// ================================================================================
// Defines ElementListImpl type as an ElementList
// TODO Composite interface

func (list *ElementListImpl) Size() int {
	if list != nil {
		return len(*list)
	}
	return -1
}

func (list *ElementListImpl) GetElementAt(i int) Element {
	if list != nil {
		if i < list.Size() {
			return (*list)[i]
		}
		return nil
	}
	return nil
}

func (list *ElementListImpl) AppendElement(element Element) {
	if list != nil {
		// TODO check element type to be the same as the type of the list elements
		*list = append(*list, element)
	}
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (list *ElementListImpl) Encode(encoder Encoder) error {
	return encoder.EncodeElementList(*list)
}

// Decodes an instance of BlobList using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded BlobList instance.
func DecodeElementList(decoder Decoder) (*ElementListImpl, error) {
	list, err := decoder.DecodeElementList()
	if err != nil {
		return nil, err
	}
	var list2 ElementListImpl = ElementListImpl(list)
	return &list2, nil
}
