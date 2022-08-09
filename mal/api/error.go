/**
 * MIT License
 *
 * Copyright (c) 2020 CNES
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
package api

import (
	"fmt"
	"github.com/CNES/ccsdsmo-malgo/mal"
)

/*
 * This structure represents an error as it can be returned by an operation according to the CCSDS MAL specification.
 * It includes a code, represented by a UInteger, and extra information, represented by an abstract Element.
 * The MalError structure also conforms to the error interface of the standard errors go package.
 * A MalError may then be returned when an error is expected.
 *
 * The MalError structure is generally not used in the malgo api functions, as the concrete type of the Element
 * in the extra information is specific to service operations, and cannot be known by generic mal functions.
 */
type MalError struct {
	Code      mal.UInteger
	ExtraInfo mal.Element
}

func NewMalError(code mal.UInteger, extraInfo ...mal.Element) *MalError {
	var e mal.Element = nil
	if len(extraInfo) == 1 {
		e = extraInfo[0]
	}
	return &MalError{code, e}
}

func (e *MalError) Error() string {
	return fmt.Sprintf("MAL error %v", e.Code)
}
