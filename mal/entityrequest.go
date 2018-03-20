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
package mal

import ()

// ################################################################################
// Defines MAL EntityRequest type
// ################################################################################

type EntityRequest struct {
	SubDomain     *IdentifierList
	AllAreas      Boolean
	AllServices   Boolean
	AllOperations Boolean
	OnlyOnChange  Boolean
	EntityKeys    EntityKeyList
}

var (
	NullEntityRequest *EntityRequest = nil
)

func NewEntityRequest() *EntityRequest {
	return new(EntityRequest)
}

// ================================================================================
// Defines MAL EntityRequest type as a MAL Composite

func (key *EntityRequest) Composite() Composite {
	return key
}

// ================================================================================
// Defines MAL EntityRequest type as a MAL Element

const MAL_ENTITY_REQUEST_TYPE_SHORT_FORM Integer = 0x18
const MAL_ENTITY_REQUEST_SHORT_FORM Long = 0x1000001000018

// Registers MAL EntityRequest type for polymorpsism handling
func init() {
	RegisterMALElement(MAL_ENTITY_REQUEST_SHORT_FORM, NullEntityRequest)
}

// Returns the absolute short form of the element type.
func (*EntityRequest) GetShortForm() Long {
	return MAL_ENTITY_REQUEST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*EntityRequest) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*EntityRequest) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*EntityRequest) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*EntityRequest) GetTypeShortForm() Integer {
	return MAL_ENTITY_REQUEST_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (request *EntityRequest) Encode(encoder Encoder) error {
	err := encoder.EncodeNullableElement(request.SubDomain)
	if err != nil {
		return err
	}
	err = encoder.EncodeBoolean(&request.AllAreas)
	if err != nil {
		return err
	}
	err = encoder.EncodeBoolean(&request.AllServices)
	if err != nil {
		return err
	}
	err = encoder.EncodeBoolean(&request.AllOperations)
	if err != nil {
		return err
	}
	err = encoder.EncodeBoolean(&request.OnlyOnChange)
	if err != nil {
		return err
	}
	return request.EntityKeys.Encode(encoder)
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (request *EntityRequest) Decode(decoder Decoder) (Element, error) {
	return DecodeEntityRequest(decoder)
}

// Decodes an instance of EntityRequest using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded EntityRequest instance.
func DecodeEntityRequest(decoder Decoder) (*EntityRequest, error) {
	var subDomain *IdentifierList
	element, err := decoder.DecodeNullableElement(subDomain)
	if err != nil {
		return nil, err
	}
	subDomain = element.(*IdentifierList)
	allAreas, err := decoder.DecodeBoolean()
	if err != nil {
		return nil, err
	}
	allServices, err := decoder.DecodeBoolean()
	if err != nil {
		return nil, err
	}
	allOperations, err := decoder.DecodeBoolean()
	if err != nil {
		return nil, err
	}
	onlyOnChange, err := decoder.DecodeBoolean()
	if err != nil {
		return nil, err
	}
	entityKeys, err := DecodeEntityKeyList(decoder)
	if err != nil {
		return nil, err
	}
	var request = EntityRequest{
		SubDomain:     subDomain,
		AllAreas:      *allAreas,
		AllServices:   *allServices,
		AllOperations: *allOperations,
		OnlyOnChange:  *onlyOnChange,
		EntityKeys:    *entityKeys,
	}
	return &request, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (request *EntityRequest) CreateElement() Element {
	// TODO (AF):
	//	return new(EntityRequest)
	return NewEntityRequest()
}

func (request *EntityRequest) IsNull() bool {
	return request == nil
}

func (*EntityRequest) Null() Element {
	return NullEntityRequest
}
