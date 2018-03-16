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
// Defines MAL EntityKey type
// ################################################################################

type EntityKey struct {
	FirstSubKey  *Identifier
	SecondSubKey *Long
	ThirdSubKey  *Long
	FourthSubKey *Long
}

var (
	NullEntityKey *EntityKey = nil
)

func NewEntityKey() *EntityKey {
	return new(EntityKey)
}

// Note (AF): This code can avoid the use of the costly reflect.DeepEqual  method.
// However it needs to be declared on Element and implemented by all types implementing
// Element.

//func (key EntityKey) Equals(other *EntityKey) bool {
//  // TODO (AF): Defines Equals with other parameter as Element
//	if (key.FirstSubKey != nil) && (other.FirstSubKey != nil) {
//		if *(key.FirstSubKey) != *(other.FirstSubKey) {
//			return false
//		}
//	} else if !((key.FirstSubKey == nil) && (other.FirstSubKey == nil)) {
//		return false
//	}
//	if (key.SecondSubKey != nil) && (other.SecondSubKey != nil) {
//		if *(key.SecondSubKey) != *(other.SecondSubKey) {
//			return false
//		}
//	} else if !((key.SecondSubKey == nil) && (other.SecondSubKey == nil)) {
//		return false
//	}
//	if (key.ThirdSubKey != nil) && (other.ThirdSubKey != nil) {
//		if *(key.ThirdSubKey) != *(other.ThirdSubKey) {
//			return false
//		}
//	} else if !((key.ThirdSubKey == nil) && (other.ThirdSubKey == nil)) {
//		return false
//	}
//	if (key.FourthSubKey != nil) && (other.FourthSubKey != nil) {
//		if *(key.FourthSubKey) != *(other.FourthSubKey) {
//			return false
//		}
//	} else if !((key.FourthSubKey == nil) && (other.FourthSubKey == nil)) {
//		return false
//	}
//	return true
//}

// ================================================================================
// Defines MAL EntityKey type as a MAL Composite

func (key *EntityKey) Composite() Composite {
	return key
}

// ================================================================================
// Defines MAL EntityKey type as a MAL Element

const MAL_ENTITY_KEY_TYPE_SHORT_FORM Integer = 0x19
const MAL_ENTITY_KEY_SHORT_FORM Long = 0x1000001000019

// Returns the absolute short form of the element type.
func (*EntityKey) GetShortForm() Long {
	return MAL_ENTITY_KEY_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*EntityKey) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*EntityKey) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*EntityKey) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*EntityKey) GetTypeShortForm() Integer {
	return MAL_ENTITY_KEY_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (key *EntityKey) Encode(encoder Encoder) error {
	err := encoder.EncodeNullableIdentifier(key.FirstSubKey)
	if err != nil {
		return err
	}
	err = encoder.EncodeNullableLong(key.SecondSubKey)
	if err != nil {
		return err
	}
	err = encoder.EncodeNullableLong(key.ThirdSubKey)
	if err != nil {
		return err
	}
	return encoder.EncodeNullableLong(key.FourthSubKey)
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (key *EntityKey) Decode(decoder Decoder) (Element, error) {
	return DecodeEntityKey(decoder)
}

// Decodes an instance of EntityKey using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded EntityKey instance.
func DecodeEntityKey(decoder Decoder) (*EntityKey, error) {
	firstSubKey, err := decoder.DecodeNullableIdentifier()
	if err != nil {
		return nil, err
	}
	secondSubKey, err := decoder.DecodeNullableLong()
	if err != nil {
		return nil, err
	}
	thirdSubKey, err := decoder.DecodeNullableLong()
	if err != nil {
		return nil, err
	}
	fourthSubKey, err := decoder.DecodeNullableLong()
	if err != nil {
		return nil, err
	}
	var key = EntityKey{
		FirstSubKey:  firstSubKey,
		SecondSubKey: secondSubKey,
		ThirdSubKey:  thirdSubKey,
		FourthSubKey: fourthSubKey,
	}
	return &key, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (key *EntityKey) CreateElement() Element {
	return new(EntityKey)
}

func (key *EntityKey) IsNull() bool {
	return key == nil
}

func (*EntityKey) Null() Element {
	return NullEntityKey
}
