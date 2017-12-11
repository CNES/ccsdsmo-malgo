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

const MAL_ATTRIBUTE_AREA_NUMBER UShort = 0x01
const MAL_ATTRIBUTE_AREA_VERSION UOctet = 0x01
const MAL_ATTRIBUTE_AREA_SERVICE_NUMBER UShort = 0x0

const (
// TODO (AF): To remove.
//	MAL_BLOB_ATTRIBUTE_TAG int = 0
//	MAL_BOOLEAN_ATTRIBUTE_TAG
//	MAL_DURATION_ATTRIBUTE_TAG
//	MAL_FLOAT_ATTRIBUTE_TAG
//	MAL_DOUBLE_ATTRIBUTE_TAG
//	MAL_IDENTIFIER_ATTRIBUTE_TAG
//	MAL_OCTET_ATTRIBUTE_TAG
//	MAL_UOCTET_ATTRIBUTE_TAG
//	MAL_SHORT_ATTRIBUTE_TAG
//	MAL_USHORT_ATTRIBUTE_TAG
//	MAL_INTEGER_ATTRIBUTE_TAG
//	MAL_UINTEGER_ATTRIBUTE_TAG
//	MAL_LONG_ATTRIBUTE_TAG
//	MAL_ULONG_ATTRIBUTE_TAG
//	MAL_STRING_ATTRIBUTE_TAG
//	MAL_TIME_ATTRIBUTE_TAG
//	MAL_FINETIME_ATTRIBUTE_TAG
//	MAL_URI_ATTRIBUTE_TAG
)

type Attribute interface {
	Element
	attribute() Attribute
}
