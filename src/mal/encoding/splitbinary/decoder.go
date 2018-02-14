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

type SplitBinaryDecoder struct {
	binary.BinaryDecoder
}

// Creates a new decoder using a slice containing binary data to decode.
func NewSplitBinaryDecoder(data []byte) *SplitBinaryDecoder {
	// Get informations from incoming slice
	var buf []byte = nil
	var bitfield []byte = nil

	var offset int = 0

	bitfield_size := int(binary.ReadUVarInt(data, &offset))
	if bitfield_size > 0 {
		bitfield = data[offset : offset+bitfield_size]
		offset += bitfield_size
	}
	buf = data[offset:]

	buffer := new(SplitBinaryBuffer)
	buffer.Offset = 0
	buffer.Buf = buf
	buffer.Bitfield_idx = 0
	buffer.Bitfield_len = uint(bitfield_size * 8)
	buffer.Bitfield = bitfield

	decoder := &SplitBinaryDecoder{
		binary.BinaryDecoder{
			Varint: true,
			In:     buffer,
		},
	}
	decoder.GenDecoder.Decoder = decoder
	return decoder
}

// TODO (AF): Normaly not needed
// Creates a new decoder using various parameters allowing to create the corresponding buffer.
//func NewSplitBinaryDecoder2(buf []byte, bitfield []byte, bitfield_len uint) *SplitBinaryDecoder {
//	buffer := new(SplitBinaryBuffer)
//	buffer.Offset = 0
//	buffer.Buf = buf
//	buffer.Bitfield_idx = 0
//	buffer.Bitfield_len = bitfield_len
//	buffer.Bitfield = bitfield
//
//	decoder := &SplitBinaryDecoder{
//		binary.BinaryDecoder{
//			Varint: true,
//			In:     buffer,
//		},
//	}
//
//	decoder.GenDecoder.Decoder = decoder
//	return decoder
//}
