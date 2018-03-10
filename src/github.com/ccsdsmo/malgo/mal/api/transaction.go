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
package api

import (
	. "github.com/ccsdsmo/malgo/mal"
)

// Defines a generic root Transaction interface (context of an incoming interaction)
type Transaction interface {
	getTid() ULong
}

// Defines a generic root Transaction structure
type TransactionX struct {
	ctx         *Context
	uri         *URI
	urifrom     *URI
	tid         ULong
	area        UShort
	areaVersion UOctet
	service     UShort
	operation   UShort
}

func (tx *TransactionX) getTid() ULong {
	return tx.tid
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
	Ack(body []byte, isError bool) error
}

type SubmitTransactionX struct {
	TransactionX
}

func (tx *SubmitTransactionX) Ack(body []byte, isError bool) error {
	msg := &Message{
		InteractionType:  MAL_INTERACTIONTYPE_SUBMIT,
		InteractionStage: MAL_IP_STAGE_SUBMIT_ACK,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		UriFrom:          tx.uri,
		UriTo:            tx.urifrom,
		IsErrorMessage:   Boolean(isError),
		Body:             body,
	}
	return tx.ctx.Send(msg)
}

// ================================================================================
// MAL Request interaction

type RequestTransaction interface {
	Transaction
	Reply(body []byte, isError bool) error
}

type RequestTransactionX struct {
	TransactionX
}

func (tx *RequestTransactionX) Reply(body []byte, isError bool) error {
	msg := &Message{
		InteractionType:  MAL_INTERACTIONTYPE_REQUEST,
		InteractionStage: MAL_IP_STAGE_REQUEST_RESPONSE,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		UriFrom:          tx.uri,
		UriTo:            tx.urifrom,
		IsErrorMessage:   Boolean(isError),
		Body:             body,
	}
	return tx.ctx.Send(msg)
}

// ================================================================================
// MAL Invoke interaction

type InvokeTransaction interface {
	Transaction
	Ack(body []byte, isError bool) error
	Reply(body []byte, isError bool) error
}

type InvokeTransactionX struct {
	TransactionX
}

func (tx *InvokeTransactionX) Ack(body []byte, isError bool) error {
	msg := &Message{
		InteractionType:  MAL_INTERACTIONTYPE_INVOKE,
		InteractionStage: MAL_IP_STAGE_INVOKE_ACK,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		UriFrom:          tx.uri,
		UriTo:            tx.urifrom,
		IsErrorMessage:   Boolean(isError),
		Body:             body,
	}
	return tx.ctx.Send(msg)
}

func (tx *InvokeTransactionX) Reply(body []byte, isError bool) error {
	msg := &Message{
		InteractionType:  MAL_INTERACTIONTYPE_INVOKE,
		InteractionStage: MAL_IP_STAGE_INVOKE_RESPONSE,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		UriFrom:          tx.uri,
		UriTo:            tx.urifrom,
		IsErrorMessage:   Boolean(isError),
		Body:             body,
	}
	return tx.ctx.Send(msg)
}

// ================================================================================
// MAL Progress interaction

type ProgressTransaction interface {
	Transaction
	Ack(body []byte, isError bool) error
	Update(body []byte, isError bool) error
	Reply(body []byte, isError bool) error
}

type ProgressTransactionX struct {
	TransactionX
}

func (tx *ProgressTransactionX) Ack(body []byte, isError bool) error {
	msg := &Message{
		InteractionType:  MAL_INTERACTIONTYPE_PROGRESS,
		InteractionStage: MAL_IP_STAGE_PROGRESS_ACK,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		UriFrom:          tx.uri,
		UriTo:            tx.urifrom,
		IsErrorMessage:   Boolean(isError),
		Body:             body,
	}
	return tx.ctx.Send(msg)
}

func (tx *ProgressTransactionX) Update(body []byte, isError bool) error {
	msg := &Message{
		InteractionType:  MAL_INTERACTIONTYPE_PROGRESS,
		InteractionStage: MAL_IP_STAGE_PROGRESS_UPDATE,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		UriFrom:          tx.uri,
		UriTo:            tx.urifrom,
		IsErrorMessage:   Boolean(isError),
		Body:             body,
	}
	return tx.ctx.Send(msg)
}

func (tx *ProgressTransactionX) Reply(body []byte, isError bool) error {
	msg := &Message{
		InteractionType:  MAL_INTERACTIONTYPE_PROGRESS,
		InteractionStage: MAL_IP_STAGE_PROGRESS_RESPONSE,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		UriFrom:          tx.uri,
		UriTo:            tx.urifrom,
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
	AckRegister(body []byte, isError bool) error
	AckDeregister(body []byte, isError bool) error
}

// SubscriberTransaction

type SubscriberTransaction interface {
	BrokerTransaction
	Notify(body []byte, isError bool) error
}

type SubscriberTransactionX struct {
	TransactionX
}

func (tx *SubscriberTransactionX) AckRegister(body []byte, isError bool) error {
	msg := &Message{
		InteractionType:  MAL_INTERACTIONTYPE_PUBSUB,
		InteractionStage: MAL_IP_STAGE_PUBSUB_REGISTER_ACK,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		UriFrom:          tx.uri,
		UriTo:            tx.urifrom,
		IsErrorMessage:   Boolean(isError),
		Body:             body,
	}
	return tx.ctx.Send(msg)
}

func (tx *SubscriberTransactionX) Notify(body []byte, isError bool) error {
	msg := &Message{
		InteractionType:  MAL_INTERACTIONTYPE_PUBSUB,
		InteractionStage: MAL_IP_STAGE_PUBSUB_NOTIFY,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		UriFrom:          tx.uri,
		UriTo:            tx.urifrom,
		IsErrorMessage:   Boolean(isError),
		Body:             body,
	}
	return tx.ctx.Send(msg)
}

func (tx *SubscriberTransactionX) AckDeregister(body []byte, isError bool) error {
	msg := &Message{
		InteractionType:  MAL_INTERACTIONTYPE_PUBSUB,
		InteractionStage: MAL_IP_STAGE_PUBSUB_DEREGISTER_ACK,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		UriFrom:          tx.uri,
		UriTo:            tx.urifrom,
		IsErrorMessage:   Boolean(isError),
		Body:             body,
	}
	return tx.ctx.Send(msg)
}

// PublisherTransaction

type PublisherTransaction interface {
	BrokerTransaction
}

type PublisherTransactionX struct {
	TransactionX
}

func (tx *PublisherTransactionX) AckRegister(body []byte, isError bool) error {
	msg := &Message{
		InteractionType:  MAL_INTERACTIONTYPE_PUBSUB,
		InteractionStage: MAL_IP_STAGE_PUBSUB_PUBLISH_REGISTER_ACK,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		UriFrom:          tx.uri,
		UriTo:            tx.urifrom,
		IsErrorMessage:   Boolean(isError),
		Body:             body,
	}
	return tx.ctx.Send(msg)
}

func (tx *PublisherTransactionX) AckDeregister(body []byte, isError bool) error {
	msg := &Message{
		InteractionType:  MAL_INTERACTIONTYPE_PUBSUB,
		InteractionStage: MAL_IP_STAGE_PUBSUB_PUBLISH_DEREGISTER_ACK,
		TransactionId:    tx.tid,
		ServiceArea:      tx.area,
		AreaVersion:      tx.areaVersion,
		Service:          tx.service,
		Operation:        tx.operation,
		UriFrom:          tx.uri,
		UriTo:            tx.urifrom,
		IsErrorMessage:   Boolean(isError),
		Body:             body,
	}
	return tx.ctx.Send(msg)
}
