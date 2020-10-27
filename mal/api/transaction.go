/**
 * MIT License
 *
 * Copyright (c) 2017 - 2020 CNES
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
	. "github.com/CNES/ccsdsmo-malgo/mal"
)

// Defines a generic root Transaction interface (context of an incoming interaction)
type Transaction interface {
	init(msg *Message)
	// Get current TransactionId
	getTid() ULong
	// Returns a new Body ready to encode
	NewBody() Body
}

// Defines a generic root Transaction structure
type TransactionX struct {
	ctx     *Context
	uri     *URI
	urifrom *URI

	AuthenticationId Blob
	EncodingId       UOctet
	QoSLevel         QoSLevel
	Priority         UInteger
	Domain           IdentifierList
	NetworkZone      Identifier
	Session          SessionType
	SessionName      Identifier

	tid         ULong
	area        UShort
	areaVersion UOctet
	service     UShort
	operation   UShort
}

// Fix additionnal parameters from incoming message
func (tx *TransactionX) init(msg *Message) {
	tx.AuthenticationId = msg.AuthenticationId
	tx.EncodingId = msg.EncodingId
	tx.QoSLevel = msg.QoSLevel
	tx.Priority = msg.Priority
	tx.Domain = msg.Domain
	tx.NetworkZone = msg.NetworkZone
	tx.Session = msg.Session
	tx.SessionName = msg.SessionName
	tx.tid = msg.TransactionId
	tx.area = msg.ServiceArea
	tx.areaVersion = msg.AreaVersion
	tx.service = msg.Service
	tx.operation = msg.Operation
}

func (tx *TransactionX) getTid() ULong {
	return tx.tid
}

func (tx *TransactionX) NewBody() Body {
	return tx.ctx.NewBody()
}

// ================================================================================
// MAL Send interaction

type SendTransaction interface {
	Transaction
}

type SendTransactionX struct {
	TransactionX
}

// ================================================================================
// MAL Submit interaction

type SubmitTransaction interface {
	Transaction
	Ack(body Body, isError bool) error
}

type SubmitTransactionX struct {
	TransactionX
}

func (tx *SubmitTransactionX) Ack(body Body, isError bool) error {
	msg := &Message{
		UriFrom:          tx.uri,
		UriTo:            tx.urifrom,
		AuthenticationId: tx.AuthenticationId,
		EncodingId:       tx.EncodingId,
		QoSLevel:         tx.QoSLevel,
		Priority:         tx.Priority,
		Domain:           tx.Domain,
		NetworkZone:      tx.NetworkZone,
		Session:          tx.Session,
		SessionName:      tx.SessionName,
		InteractionType:  MAL_INTERACTIONTYPE_SUBMIT,
		InteractionStage: MAL_IP_STAGE_SUBMIT_ACK,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		IsErrorMessage:   Boolean(isError),
		Body:             body,
	}
	return tx.ctx.Send(msg)
}

// ================================================================================
// MAL Request interaction

type RequestTransaction interface {
	Transaction
	Reply(body Body, isError bool) error
}

type RequestTransactionX struct {
	TransactionX
}

func (tx *RequestTransactionX) Reply(body Body, isError bool) error {
	msg := &Message{
		UriFrom:          tx.uri,
		UriTo:            tx.urifrom,
		AuthenticationId: tx.AuthenticationId,
		EncodingId:       tx.EncodingId,
		QoSLevel:         tx.QoSLevel,
		Priority:         tx.Priority,
		Domain:           tx.Domain,
		NetworkZone:      tx.NetworkZone,
		Session:          tx.Session,
		SessionName:      tx.SessionName,
		InteractionType:  MAL_INTERACTIONTYPE_REQUEST,
		InteractionStage: MAL_IP_STAGE_REQUEST_RESPONSE,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		Service:          tx.service,
		Operation:        tx.operation,
		AreaVersion:      tx.areaVersion,
		IsErrorMessage:   Boolean(isError),
		Body:             body,
	}
	return tx.ctx.Send(msg)
}

// ================================================================================
// MAL Invoke interaction

type InvokeTransaction interface {
	Transaction
	Ack(body Body, isError bool) error
	Reply(body Body, isError bool) error
}

type InvokeTransactionX struct {
	TransactionX
}

func (tx *InvokeTransactionX) Ack(body Body, isError bool) error {
	msg := &Message{
		UriFrom:          tx.uri,
		UriTo:            tx.urifrom,
		AuthenticationId: tx.AuthenticationId,
		EncodingId:       tx.EncodingId,
		QoSLevel:         tx.QoSLevel,
		Priority:         tx.Priority,
		Domain:           tx.Domain,
		NetworkZone:      tx.NetworkZone,
		Session:          tx.Session,
		SessionName:      tx.SessionName,
		InteractionType:  MAL_INTERACTIONTYPE_INVOKE,
		InteractionStage: MAL_IP_STAGE_INVOKE_ACK,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		IsErrorMessage:   Boolean(isError),
		Body:             body,
	}
	return tx.ctx.Send(msg)
}

func (tx *InvokeTransactionX) Reply(body Body, isError bool) error {
	msg := &Message{
		UriFrom:          tx.uri,
		UriTo:            tx.urifrom,
		AuthenticationId: tx.AuthenticationId,
		EncodingId:       tx.EncodingId,
		QoSLevel:         tx.QoSLevel,
		Priority:         tx.Priority,
		Domain:           tx.Domain,
		NetworkZone:      tx.NetworkZone,
		Session:          tx.Session,
		SessionName:      tx.SessionName,
		InteractionType:  MAL_INTERACTIONTYPE_INVOKE,
		InteractionStage: MAL_IP_STAGE_INVOKE_RESPONSE,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		IsErrorMessage:   Boolean(isError),
		Body:             body,
	}
	return tx.ctx.Send(msg)
}

// ================================================================================
// MAL Progress interaction

type ProgressTransaction interface {
	Transaction
	Ack(body Body, isError bool) error
	Update(body Body, isError bool) error
	Reply(body Body, isError bool) error
}

type ProgressTransactionX struct {
	TransactionX
}

func (tx *ProgressTransactionX) Ack(body Body, isError bool) error {
	msg := &Message{
		UriFrom:          tx.uri,
		UriTo:            tx.urifrom,
		AuthenticationId: tx.AuthenticationId,
		EncodingId:       tx.EncodingId,
		QoSLevel:         tx.QoSLevel,
		Priority:         tx.Priority,
		Domain:           tx.Domain,
		NetworkZone:      tx.NetworkZone,
		Session:          tx.Session,
		SessionName:      tx.SessionName,
		InteractionType:  MAL_INTERACTIONTYPE_PROGRESS,
		InteractionStage: MAL_IP_STAGE_PROGRESS_ACK,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		IsErrorMessage:   Boolean(isError),
		Body:             body,
	}
	return tx.ctx.Send(msg)
}

func (tx *ProgressTransactionX) Update(body Body, isError bool) error {
	msg := &Message{
		UriFrom:          tx.uri,
		UriTo:            tx.urifrom,
		AuthenticationId: tx.AuthenticationId,
		EncodingId:       tx.EncodingId,
		QoSLevel:         tx.QoSLevel,
		Priority:         tx.Priority,
		Domain:           tx.Domain,
		NetworkZone:      tx.NetworkZone,
		Session:          tx.Session,
		SessionName:      tx.SessionName,
		InteractionType:  MAL_INTERACTIONTYPE_PROGRESS,
		InteractionStage: MAL_IP_STAGE_PROGRESS_UPDATE,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		IsErrorMessage:   Boolean(isError),
		Body:             body,
	}
	return tx.ctx.Send(msg)
}

func (tx *ProgressTransactionX) Reply(body Body, isError bool) error {
	msg := &Message{
		UriFrom:          tx.uri,
		UriTo:            tx.urifrom,
		AuthenticationId: tx.AuthenticationId,
		EncodingId:       tx.EncodingId,
		QoSLevel:         tx.QoSLevel,
		Priority:         tx.Priority,
		Domain:           tx.Domain,
		NetworkZone:      tx.NetworkZone,
		Session:          tx.Session,
		SessionName:      tx.SessionName,
		InteractionType:  MAL_INTERACTIONTYPE_PROGRESS,
		InteractionStage: MAL_IP_STAGE_PROGRESS_RESPONSE,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		IsErrorMessage:   Boolean(isError),
		Body:             body,
	}
	return tx.ctx.Send(msg)
}

/// ================================================================================
// MAL PubSub interaction

// There is only one handler but 2 transactions type depending of the incoming interaction.

type BrokerTransaction interface {
	Transaction
	AckRegister(body Body, isError bool) error
	AckDeregister(body Body, isError bool) error
}

// SubscriberTransaction

type SubscriberTransaction interface {
	BrokerTransaction
	Notify(body Body, isError bool) error
}

type SubscriberTransactionX struct {
	TransactionX
}

func (tx *SubscriberTransactionX) AckRegister(body Body, isError bool) error {
	msg := &Message{
		UriFrom:          tx.uri,
		UriTo:            tx.urifrom,
		AuthenticationId: tx.AuthenticationId,
		EncodingId:       tx.EncodingId,
		QoSLevel:         tx.QoSLevel,
		Priority:         tx.Priority,
		Domain:           tx.Domain,
		NetworkZone:      tx.NetworkZone,
		Session:          tx.Session,
		SessionName:      tx.SessionName,
		InteractionType:  MAL_INTERACTIONTYPE_PUBSUB,
		InteractionStage: MAL_IP_STAGE_PUBSUB_REGISTER_ACK,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		IsErrorMessage:   Boolean(isError),
		Body:             body,
	}
	return tx.ctx.Send(msg)
}

func (tx *SubscriberTransactionX) Notify(body Body, isError bool) error {
	msg := &Message{
		UriFrom:          tx.uri,
		UriTo:            tx.urifrom,
		AuthenticationId: tx.AuthenticationId,
		EncodingId:       tx.EncodingId,
		QoSLevel:         tx.QoSLevel,
		Priority:         tx.Priority,
		Domain:           tx.Domain,
		NetworkZone:      tx.NetworkZone,
		Session:          tx.Session,
		SessionName:      tx.SessionName,
		InteractionType:  MAL_INTERACTIONTYPE_PUBSUB,
		InteractionStage: MAL_IP_STAGE_PUBSUB_NOTIFY,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		IsErrorMessage:   Boolean(isError),
		Body:             body,
	}
	return tx.ctx.Send(msg)
}

func (tx *SubscriberTransactionX) AckDeregister(body Body, isError bool) error {
	msg := &Message{
		UriFrom:          tx.uri,
		UriTo:            tx.urifrom,
		AuthenticationId: tx.AuthenticationId,
		EncodingId:       tx.EncodingId,
		QoSLevel:         tx.QoSLevel,
		Priority:         tx.Priority,
		Domain:           tx.Domain,
		NetworkZone:      tx.NetworkZone,
		Session:          tx.Session,
		SessionName:      tx.SessionName,
		InteractionType:  MAL_INTERACTIONTYPE_PUBSUB,
		InteractionStage: MAL_IP_STAGE_PUBSUB_DEREGISTER_ACK,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		IsErrorMessage:   Boolean(isError),
		Body:             body,
	}
	return tx.ctx.Send(msg)
}

// PublisherTransaction

type PublisherTransaction interface {
	BrokerTransaction
	// PublishError to return error on Publish.
	PublishError(body Body) error
}

type PublisherTransactionX struct {
	TransactionX
}

func (tx *PublisherTransactionX) AckRegister(body Body, isError bool) error {
	msg := &Message{
		UriFrom:          tx.uri,
		UriTo:            tx.urifrom,
		AuthenticationId: tx.AuthenticationId,
		EncodingId:       tx.EncodingId,
		QoSLevel:         tx.QoSLevel,
		Priority:         tx.Priority,
		Domain:           tx.Domain,
		NetworkZone:      tx.NetworkZone,
		Session:          tx.Session,
		SessionName:      tx.SessionName,
		InteractionType:  MAL_INTERACTIONTYPE_PUBSUB,
		InteractionStage: MAL_IP_STAGE_PUBSUB_PUBLISH_REGISTER_ACK,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		IsErrorMessage:   Boolean(isError),
		Body:             body,
	}
	return tx.ctx.Send(msg)
}

func (tx *PublisherTransactionX) AckDeregister(body Body, isError bool) error {
	msg := &Message{
		UriFrom:          tx.uri,
		UriTo:            tx.urifrom,
		AuthenticationId: tx.AuthenticationId,
		EncodingId:       tx.EncodingId,
		QoSLevel:         tx.QoSLevel,
		Priority:         tx.Priority,
		Domain:           tx.Domain,
		NetworkZone:      tx.NetworkZone,
		Session:          tx.Session,
		SessionName:      tx.SessionName,
		InteractionType:  MAL_INTERACTIONTYPE_PUBSUB,
		InteractionStage: MAL_IP_STAGE_PUBSUB_PUBLISH_DEREGISTER_ACK,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		IsErrorMessage:   Boolean(isError),
		Body:             body,
	}
	return tx.ctx.Send(msg)
}

// TODO (AF): Get fields from original Register
func (tx *PublisherTransactionX) PublishError(body Body) error {
	msg := &Message{
		UriFrom:          tx.uri,
		UriTo:            tx.urifrom,
		AuthenticationId: tx.AuthenticationId,
		EncodingId:       tx.EncodingId,
		QoSLevel:         tx.QoSLevel,
		Priority:         tx.Priority,
		Domain:           tx.Domain,
		NetworkZone:      tx.NetworkZone,
		Session:          tx.Session,
		SessionName:      tx.SessionName,
		InteractionType:  MAL_INTERACTIONTYPE_PUBSUB,
		InteractionStage: MAL_IP_STAGE_PUBSUB_PUBLISH_REGISTER_ACK,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		IsErrorMessage:   Boolean(true),
		Body:             body,
	}
	return tx.ctx.Send(msg)
}
