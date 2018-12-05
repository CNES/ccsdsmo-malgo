/**
 * MIT License
 *
 * Copyright (c) 2017 - 2018 CNES
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
	. "github.com/CNES/ccsdsmo-malgo/mal"
)

const (
	TRUE  byte = 1
	FALSE byte = 0
)

var (
	FixedBinaryEncodingFactory  *FixedBinaryEncoding  = nil
	VarintBinaryEncodingFactory *VarintBinaryEncoding = nil
)

type FixedBinaryEncoding struct{}

func (*FixedBinaryEncoding) NewEncoder(buf []byte) Encoder {
	return NewBinaryEncoder(buf, false)
}

func (*FixedBinaryEncoding) NewDecoder(buf []byte) Decoder {
	return NewBinaryDecoder(buf, false)
}

type VarintBinaryEncoding struct{}

func (*VarintBinaryEncoding) NewEncoder(buf []byte) Encoder {
	return NewBinaryEncoder(buf, true)
}

func (*VarintBinaryEncoding) NewDecoder(buf []byte) Decoder {
	return NewBinaryDecoder(buf, true)
}
