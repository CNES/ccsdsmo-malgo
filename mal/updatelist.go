/**
 * MIT License
 *
 * Copyright (c) 2017  2018 CNES
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

// ################################################################################
// Defines MAL UpdateList type
// ################################################################################

/**
 * This type is needed to create generic broker (cf. CCSDS 524.1-B-1 3.5.3.3.1 p 3-15).
 */
type UpdateList []*Blob

var (
	NullUpdateList *UpdateList = nil
)

func NewUpdateList(size int) *UpdateList {
	var updates UpdateList = UpdateList(make([]*Blob, size))
	return &updates
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (list *UpdateList) Encode(encoder Encoder) error {
	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*Blob(*list)))))
	if err != nil {
		return err
	}
	for _, e := range []*Blob(*list) {
		encoder.EncodeNullableBlob(e)
	}
	return nil
}

// Decodes an instance of UpdateList using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded UpdateList instance.
func DecodeUpdateList(decoder Decoder) (*UpdateList, error) {
	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := UpdateList(make([]*Blob, int(*size)))
	for i := 0; i < len(list); i++ {
		list[i], err = decoder.DecodeNullableBlob()
		if err != nil {
			return nil, err
		}
	}
	return &list, nil
}
