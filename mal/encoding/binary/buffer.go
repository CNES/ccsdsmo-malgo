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
package binary

import (
	"errors"
)

var empty []byte

type BinaryBuffer struct {
	Offset int
	Buf    []byte
}

// Returns the slice containing all encoded datas.
func (buffer *BinaryBuffer) Body() []byte {
	// No longer duplicate the buffer
	//	buf := make([]byte, len(buffer.Buf))
	//	copy(buf, buffer.Buf)
	return buffer.Buf
}

// Reset the buffer allowing to reuse it.
// If write parameter is true the internal slice is cleaned and the buffer
// can be write anew.
func (buffer *BinaryBuffer) Reset(write bool) {
	if write {
		buffer.Buf = buffer.Buf[:0]
	}
	buffer.Offset = 0
}

func (buffer *BinaryBuffer) Write(value byte) error {
	buffer.Buf = append(buffer.Buf, value)
	return nil
}

func (buffer *BinaryBuffer) Write16(value uint16) error {
	buffer.Buf = append(buffer.Buf, byte(value>>8), byte(value>>0))
	return nil
}

func (buffer *BinaryBuffer) Write32(value uint32) error {
	buffer.Buf = append(buffer.Buf, byte(value>>24), byte(value>>16), byte(value>>8), byte(value>>0))
	return nil
}

func (buffer *BinaryBuffer) Write64(value uint64) error {
	buffer.Buf = append(buffer.Buf,
		byte(value>>56), byte(value>>48), byte(value>>40), byte(value>>32),
		byte(value>>24), byte(value>>16), byte(value>>8), byte(value>>0))
	return nil
}

// Writes an unsigned varint as defined in 5.25 section of the specification.
func (buffer *BinaryBuffer) WriteUVarInt(value uint64) error {
	buffer.Buf = WriteUVarInt(value, buffer.Buf)
	return nil
}

func WriteUVarInt(value uint64, buf []byte) []byte {
	for (value & 0xFFFFFFFFFFFFFF80) != 0 {
		buf = append(buf, byte((value&0x7F)|0x80))
		value >>= 7
	}
	buf = append(buf, byte(value&0x7F))
	return buf
}

func (buffer *BinaryBuffer) WriteFlag(value bool) error {
	if value {
		buffer.Buf = append(buffer.Buf, TRUE)
	} else {
		buffer.Buf = append(buffer.Buf, FALSE)
	}
	return nil
}

func (buffer *BinaryBuffer) WriteBytes(value []byte) error {
	buffer.Buf = append(buffer.Buf, value...)
	return nil
}

func (buffer *BinaryBuffer) Read() (byte, error) {
	b := buffer.Buf[buffer.Offset]
	buffer.Offset += 1
	return b, nil
}

func (buffer *BinaryBuffer) Read16() (uint16, error) {
	s := uint16(buffer.Buf[buffer.Offset+1]) | uint16(buffer.Buf[buffer.Offset])<<8
	buffer.Offset += 2
	return s, nil
}

func (buffer *BinaryBuffer) Read32() (uint32, error) {
	i := uint32(buffer.Buf[buffer.Offset+3]) | uint32(buffer.Buf[buffer.Offset+2])<<8 |
		uint32(buffer.Buf[buffer.Offset+1])<<16 | uint32(buffer.Buf[buffer.Offset])<<24
	buffer.Offset += 4
	return i, nil
}

func (buffer *BinaryBuffer) Read64() (uint64, error) {
	l := uint64(buffer.Buf[buffer.Offset+7]) | uint64(buffer.Buf[buffer.Offset+6])<<8 |
		uint64(buffer.Buf[buffer.Offset+5])<<16 | uint64(buffer.Buf[buffer.Offset+4])<<24 |
		uint64(buffer.Buf[buffer.Offset+3])<<32 | uint64(buffer.Buf[buffer.Offset+2])<<40 |
		uint64(buffer.Buf[buffer.Offset+1])<<48 | uint64(buffer.Buf[buffer.Offset])<<56
	buffer.Offset += 8
	return l, nil
}

// Reads an unsigned varint as defined in 5.25 section of the specification.
func (buffer *BinaryBuffer) ReadUVarInt() (uint64, error) {
	value := ReadUVarInt(buffer.Buf, &buffer.Offset)
	return value, nil
}

func ReadUVarInt(buf []byte, offset *int) uint64 {
	var value uint64 = 0
	var i uint = 0

	for {
		b := buf[*offset]
		*offset += 1
		value |= uint64(b&0x7F) << i
		if (b & 0x80) == 0 {
			break
		}
		i += 7
	}

	return value
}

func (buffer *BinaryBuffer) ReadFlag() (bool, error) {
	b := buffer.Buf[buffer.Offset]
	buffer.Offset += 1
	if b == FALSE {
		return false, nil
	} else if b == TRUE {
		return true, nil
	}
	return false, errors.New("Bad boolean encoding: " + string(b))
}

func (buffer *BinaryBuffer) ReadBytes(buf []byte) error {
	_ = buffer.Buf[buffer.Offset+len(buf)-1] // bounds check hint to compiler; see golang.org/issue/14808
	copy(buf, buffer.Buf[buffer.Offset:])
	buffer.Offset += len(buf)
	return nil
}

// Returns the part of buffer that still needs to be decoded
func (buffer *BinaryBuffer) Remaining() []byte {
	if buffer.Offset == len(buffer.Buf) {
		return empty
	} else {
		_ = buffer.Buf[buffer.Offset] // bounds check hint to compiler; see golang.org/issue/14808
		off := buffer.Offset
		buffer.Offset = len(buffer.Buf)
		return buffer.Buf[off:]
	}
}
