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
	"github.com/ccsdsmo/malgo/src/mal/encoding/binary"
)

var empty []byte

type SplitBinaryBuffer struct {
	binary.BinaryBuffer
	Bitfield_idx uint
	Bitfield_len uint
	Bitfield     []byte
}

// Returns a slice containing all encoded datas
func (buffer *SplitBinaryBuffer) Body() []byte {
	//	var length int = 4 + len(buffer.Bitfield) + len(buffer.Buf)
	var bflen int = ((((int)(buffer.Bitfield_len) - 1) / 8) + 1)
	var length int = 4 + bflen + len(buffer.Buf)
	buf := make([]byte, 0, length)
	if buffer.Bitfield_len == 0 {
		buf = binary.WriteUVarInt(0, buf)
	} else {
		buf = binary.WriteUVarInt(uint64(bflen), buf)
		buf = append(buf, buffer.Bitfield[0:bflen]...)
	}
	buf = append(buf, buffer.Buf...)

	return buf
}

// Reset the buffer allowing to reuse it to write anew
func (buffer *SplitBinaryBuffer) Reset(write bool) {
	if write {
		buffer.Buf = buffer.Buf[:0]
		buffer.Bitfield_len = 0
		buffer.Bitfield = buffer.Bitfield[:0]
	}
	buffer.Offset = 0
	buffer.Bitfield_idx = 0
}

func (buffer *SplitBinaryBuffer) WriteFlag(value bool) error {
	if (buffer.Bitfield_idx % 8) == 0 {
		// Ensures that the buffer has space for next bit
		buffer.Bitfield = append(buffer.Bitfield, 0)
	}
	if value {
		buffer.Bitfield_len = buffer.Bitfield_idx + 1
		buffer.Bitfield[buffer.Bitfield_idx>>3] |= 1 << (buffer.Bitfield_idx & 7)
	}
	buffer.Bitfield_idx += 1
	return nil
}

func (buffer *SplitBinaryBuffer) ReadFlag() (bool, error) {
	if buffer.Bitfield_idx >= buffer.Bitfield_len {
		return false, nil
	}
	res := (buffer.Bitfield[buffer.Bitfield_idx>>3] & (1 << (buffer.Bitfield_idx & 0x7))) != 0
	buffer.Bitfield_idx += 1
	return res, nil
}
