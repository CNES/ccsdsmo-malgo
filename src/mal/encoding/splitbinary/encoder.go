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
package splitbinary

import (
	"mal/encoding/binary"
)

import ()

type SplitBinaryEncoder struct {
	binary.BinaryEncoder
}

// Creates a new encoder using a slice with sufficient capacity to encode datas.
// If the slice is not empty the encoded datas are append afterwards.
func NewSplitBinaryEncoder(buf []byte, bitfield []byte) *SplitBinaryEncoder {
	buffer := new(SplitBinaryBuffer)
	buffer.Offset = 0
	buffer.Buf = buf
	buffer.Bitfield_idx = 0
	buffer.Bitfield_len = 0
	buffer.Bitfield = bitfield

	encoder := &SplitBinaryEncoder{
		binary.BinaryEncoder{
			Varint: true,
			Out:    buffer,
		},
	}

	encoder.GenEncoder.Encoder = encoder
	return encoder
}

// Returns a slice containing all encoded data as needed to be sent
func (encoder *SplitBinaryEncoder) Body() []byte {
	return encoder.Out.(*SplitBinaryBuffer).Body()
}

//func (encoder *SplitBinaryEncoder) Bitfield() []byte {
//	return encoder.Out.Buffer()
//}

// TODO (AF): To remove
//func (encoder *SplitBinaryEncoder) InitSplitBinaryEncoder(varint bool) {
//	encoder.InitBinaryEncoder(varint)
//}
