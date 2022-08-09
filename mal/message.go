/**
 * MIT License
 *
 * Copyright (c) 2017 - 2019 ccsdsmo
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

// ================================================================================
// Defines InteractionType enumeration

type InteractionStage UOctet

const (
	MAL_IP_STAGE_INIT                          InteractionStage = 1
	MAL_IP_STAGE_SEND                          InteractionStage = MAL_IP_STAGE_INIT
	MAL_IP_STAGE_SUBMIT                        InteractionStage = MAL_IP_STAGE_INIT
	MAL_IP_STAGE_SUBMIT_ACK                    InteractionStage = 2
	MAL_IP_STAGE_REQUEST                       InteractionStage = MAL_IP_STAGE_INIT
	MAL_IP_STAGE_REQUEST_RESPONSE              InteractionStage = 2
	MAL_IP_STAGE_INVOKE                        InteractionStage = MAL_IP_STAGE_INIT
	MAL_IP_STAGE_INVOKE_ACK                    InteractionStage = 2
	MAL_IP_STAGE_INVOKE_RESPONSE               InteractionStage = 3
	MAL_IP_STAGE_PROGRESS                      InteractionStage = MAL_IP_STAGE_INIT
	MAL_IP_STAGE_PROGRESS_ACK                  InteractionStage = 2
	MAL_IP_STAGE_PROGRESS_UPDATE               InteractionStage = 3
	MAL_IP_STAGE_PROGRESS_RESPONSE             InteractionStage = 4
	MAL_IP_STAGE_PUBSUB_REGISTER               InteractionStage = MAL_IP_STAGE_INIT
	MAL_IP_STAGE_PUBSUB_REGISTER_ACK           InteractionStage = 2
	MAL_IP_STAGE_PUBSUB_PUBLISH_REGISTER       InteractionStage = 3
	MAL_IP_STAGE_PUBSUB_PUBLISH_REGISTER_ACK   InteractionStage = 4
	MAL_IP_STAGE_PUBSUB_PUBLISH                InteractionStage = 5
	MAL_IP_STAGE_PUBSUB_NOTIFY                 InteractionStage = 6
	MAL_IP_STAGE_PUBSUB_DEREGISTER             InteractionStage = 7
	MAL_IP_STAGE_PUBSUB_DEREGISTER_ACK         InteractionStage = 8
	MAL_IP_STAGE_PUBSUB_PUBLISH_DEREGISTER     InteractionStage = 9
	MAL_IP_STAGE_PUBSUB_PUBLISH_DEREGISTER_ACK InteractionStage = 10
)

type InteractionType UOctet

const (
	MAL_INTERACTIONTYPE_SEND InteractionType = iota + 1
	MAL_INTERACTIONTYPE_SUBMIT
	MAL_INTERACTIONTYPE_REQUEST
	MAL_INTERACTIONTYPE_INVOKE
	MAL_INTERACTIONTYPE_PROGRESS
	MAL_INTERACTIONTYPE_PUBSUB
)

const MAL_INTERACTIONTYPE_TYPE_SHORT_FORM Integer = 0x13
const MAL_INTERACTIONTYPE_SHORT_FORM Long = 0x1000001000013

// ================================================================================
// Defines SessionType enumeration

type SessionType UOctet

const (
	MAL_SESSIONTYPE_LIVE SessionType = iota + 1
	MAL_SESSIONTYPE_SIMULATION
	MAL_SESSIONTYPE_REPLAY
)

const MAL_SESSIONTYPE_TYPE_SHORT_FORM Integer = 0x14
const MAL_SESSIONTYPE_SHORT_FORM Long = 0x1000001000014

// ================================================================================
// Defines QoSLevel enumeration

type QoSLevel UOctet

const (
	MAL_QOSLEVEL_BESTEFFORT QoSLevel = iota + 1
	MAL_QOSLEVEL_ASSURED
	MAL_QOSLEVEL_QUEUED
	MAL_QOSLEVEL_TIMELY
)

const MAL_QOSLEVEL_TYPE_SHORT_FORM Integer = 0x15
const MAL_QOSLEVEL_SHORT_FORM Long = 0x1000001000015

// ================================================================================
// Defines UpdateType enumeration

type UpdateType uint8

const (
	MAL_UPDATETYPE_CREATION UpdateType = iota + 1
	MAL_UPDATETYPE_UPDATE
	MAL_UPDATETYPE_MODIFICATION
	MAL_UPDATETYPE_DELETION
)

const MAL_UPDATETYPE_TYPE_SHORT_FORM Integer = 0x16
const MAL_UPDATETYPE_SHORT_FORM Long = 0x1000001000016

// ================================================================================
// Defines MAL Message structure and interface

// TODO (AF): Handles the polymorphism of message, use an interface implemented by
// MAL transport.

type Body interface {
	SetEncodingFactory(factory EncodingFactory)
	// Initialize the body to write (resp. read) parameter using the specified
	// encoding. This method shall be called after setting content and factory.
	Reset(writeable bool)
	DecodeParameter(element Element) (Element, error)
	DecodeLastParameter(element Element, abstract bool) (Element, error)
	EncodeParameter(element Element) error
	EncodeLastParameter(element Element, abstract bool) error
}

type Message struct {
	UriFrom          *URI
	UriTo            *URI
	AuthenticationId Blob
	EncodingId       UOctet
	Timestamp        Time
	QoSLevel         QoSLevel
	Priority         UInteger
	Domain           IdentifierList
	NetworkZone      Identifier
	Session          SessionType
	SessionName      Identifier
	InteractionType  InteractionType
	InteractionStage InteractionStage
	TransactionId    ULong
	ServiceArea      UShort
	Service          UShort
	Operation        UShort
	AreaVersion      UOctet
	IsErrorMessage   Boolean
	Body             Body
}

func (msg *Message) SetEncodingFactory(factory EncodingFactory) {
	msg.Body.SetEncodingFactory(factory)
}

func (msg *Message) DecodeParameter(element Element) (Element, error) {
	return msg.Body.DecodeParameter(element)
}

func (msg *Message) DecodeLastParameter(element Element, abstract bool) (Element, error) {
	return msg.Body.DecodeLastParameter(element, abstract)
}

func (msg *Message) EncodeParameter(element Element) error {
	return msg.Body.EncodeParameter(element)
}

func (msg *Message) EncodeLastParameter(element Element, abstract bool) error {
	return msg.Body.EncodeLastParameter(element, abstract)
}
