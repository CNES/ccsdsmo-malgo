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
package binary

import ()

var empty []byte

type BinaryBuffer struct {
	offset int
	buf    []byte
}

func (buffer *BinaryBuffer) Buffer() []byte {
	return buffer.buf
}

func (buffer *BinaryBuffer) Reset() {
	buffer.buf = buffer.buf[:0]
}

func (buffer *BinaryBuffer) Write(value byte) error {
	buffer.buf = append(buffer.buf, value)
	return nil
}

func (buffer *BinaryBuffer) Write16(value uint16) error {
	buffer.buf = append(buffer.buf, byte(value>>8), byte(value>>0))
	return nil
}

func (buffer *BinaryBuffer) Write32(value uint32) error {
	buffer.buf = append(buffer.buf, byte(value>>24), byte(value>>16), byte(value>>8), byte(value>>0))
	return nil
}

func (buffer *BinaryBuffer) Write64(value uint64) error {
	buffer.buf = append(buffer.buf,
		byte(value>>56), byte(value>>48), byte(value>>40), byte(value>>32),
		byte(value>>24), byte(value>>16), byte(value>>8), byte(value>>0))
	return nil
}

func (buffer *BinaryBuffer) WriteBytes(value []byte) error {
	buffer.buf = append(buffer.buf, value...)
	return nil
}

func (buffer *BinaryBuffer) Read() (byte, error) {
	b := buffer.buf[buffer.offset]
	buffer.offset += 1
	return b, nil
}

func (buffer *BinaryBuffer) Read16() (uint16, error) {
	s := uint16(buffer.buf[buffer.offset+1]) | uint16(buffer.buf[buffer.offset])<<8
	buffer.offset += 2
	return s, nil
}

func (buffer *BinaryBuffer) Read32() (uint32, error) {
	i := uint32(buffer.buf[buffer.offset+3]) | uint32(buffer.buf[buffer.offset+2])<<8 |
		uint32(buffer.buf[buffer.offset+1])<<16 | uint32(buffer.buf[buffer.offset])<<24
	buffer.offset += 4
	return i, nil
}

func (buffer *BinaryBuffer) Read64() (uint64, error) {
	l := uint64(buffer.buf[buffer.offset+7]) | uint64(buffer.buf[buffer.offset+6])<<8 |
		uint64(buffer.buf[buffer.offset+5])<<16 | uint64(buffer.buf[buffer.offset+4])<<24 |
		uint64(buffer.buf[buffer.offset+3])<<32 | uint64(buffer.buf[buffer.offset+2])<<40 |
		uint64(buffer.buf[buffer.offset+1])<<48 | uint64(buffer.buf[buffer.offset])<<56
	buffer.offset += 8
	return l, nil
}

func (buffer *BinaryBuffer) ReadBytes(buf []byte) error {
	_ = buffer.buf[buffer.offset+len(buf)-1] // bounds check hint to compiler; see golang.org/issue/14808
	copy(buf, buffer.buf[buffer.offset:])
	buffer.offset += len(buf)
	return nil
}

func (buffer *BinaryBuffer) Remaining() ([]byte, error) {
	if buffer.offset == len(buffer.buf) {
		return empty, nil
	} else {
		_ = buffer.buf[buffer.offset] // bounds check hint to compiler; see golang.org/issue/14808
		off := buffer.offset
		buffer.offset = len(buffer.buf)
		return buffer.buf[off:], nil
	}
}
