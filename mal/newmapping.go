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
package mal

import ()

// This file implements as much as possible of the new go language mapping defined along with the generator

const (
	// area constants
	AREA_NUMBER  UShort = 1
	AREA_VERSION UOctet = 1
	AREA_NAME           = Identifier("MAL")

	// type constants
	BLOB_TYPE_SHORT_FORM                 Integer = 1
	BLOB_SHORT_FORM                      Long    = 0x1000001000001
	BLOB_LIST_TYPE_SHORT_FORM            Integer = -1
	BLOB_LIST_SHORT_FORM                 Long    = 0x1000001FFFFFF
	BOOLEAN_TYPE_SHORT_FORM              Integer = 2
	BOOLEAN_SHORT_FORM                   Long    = 0x1000001000002
	BOOLEAN_LIST_TYPE_SHORT_FORM         Integer = -2
	BOOLEAN_LIST_SHORT_FORM              Long    = 0x1000001FFFFFE
	DURATION_TYPE_SHORT_FORM             Integer = 3
	DURATION_SHORT_FORM                  Long    = 0x1000001000003
	DURATION_LIST_TYPE_SHORT_FORM        Integer = -3
	DURATION_LIST_SHORT_FORM             Long    = 0x1000001FFFFFD
	FLOAT_TYPE_SHORT_FORM                Integer = 4
	FLOAT_SHORT_FORM                     Long    = 0x1000001000004
	FLOAT_LIST_TYPE_SHORT_FORM           Integer = -4
	FLOAT_LIST_SHORT_FORM                Long    = 0x1000001FFFFFC
	DOUBLE_TYPE_SHORT_FORM               Integer = 5
	DOUBLE_SHORT_FORM                    Long    = 0x1000001000005
	DOUBLE_LIST_TYPE_SHORT_FORM          Integer = -5
	DOUBLE_LIST_SHORT_FORM               Long    = 0x1000001FFFFFB
	IDENTIFIER_TYPE_SHORT_FORM           Integer = 6
	IDENTIFIER_SHORT_FORM                Long    = 0x1000001000006
	IDENTIFIER_LIST_TYPE_SHORT_FORM      Integer = -6
	IDENTIFIER_LIST_SHORT_FORM           Long    = 0x1000001FFFFFA
	OCTET_TYPE_SHORT_FORM                Integer = 7
	OCTET_SHORT_FORM                     Long    = 0x1000001000007
	OCTET_LIST_TYPE_SHORT_FORM           Integer = -7
	OCTET_LIST_SHORT_FORM                Long    = 0x1000001FFFFF9
	UOCTET_TYPE_SHORT_FORM               Integer = 8
	UOCTET_SHORT_FORM                    Long    = 0x1000001000008
	UOCTET_LIST_TYPE_SHORT_FORM          Integer = -8
	UOCTET_LIST_SHORT_FORM               Long    = 0x1000001FFFFF8
	SHORT_TYPE_SHORT_FORM                Integer = 9
	SHORT_SHORT_FORM                     Long    = 0x1000001000009
	SHORT_LIST_TYPE_SHORT_FORM           Integer = -9
	SHORT_LIST_SHORT_FORM                Long    = 0x1000001FFFFF7
	USHORT_TYPE_SHORT_FORM               Integer = 10
	USHORT_SHORT_FORM                    Long    = 0x100000100000A
	USHORT_LIST_TYPE_SHORT_FORM          Integer = -10
	USHORT_LIST_SHORT_FORM               Long    = 0x1000001FFFFF6
	INTEGER_TYPE_SHORT_FORM              Integer = 11
	INTEGER_SHORT_FORM                   Long    = 0x100000100000B
	INTEGER_LIST_TYPE_SHORT_FORM         Integer = -11
	INTEGER_LIST_SHORT_FORM              Long    = 0x1000001FFFFF5
	UINTEGER_TYPE_SHORT_FORM             Integer = 12
	UINTEGER_SHORT_FORM                  Long    = 0x100000100000C
	UINTEGER_LIST_TYPE_SHORT_FORM        Integer = -12
	UINTEGER_LIST_SHORT_FORM             Long    = 0x1000001FFFFF4
	LONG_TYPE_SHORT_FORM                 Integer = 13
	LONG_SHORT_FORM                      Long    = 0x100000100000D
	LONG_LIST_TYPE_SHORT_FORM            Integer = -13
	LONG_LIST_SHORT_FORM                 Long    = 0x1000001FFFFF3
	ULONG_TYPE_SHORT_FORM                Integer = 14
	ULONG_SHORT_FORM                     Long    = 0x100000100000E
	ULONG_LIST_TYPE_SHORT_FORM           Integer = -14
	ULONG_LIST_SHORT_FORM                Long    = 0x1000001FFFFF2
	STRING_TYPE_SHORT_FORM               Integer = 15
	STRING_SHORT_FORM                    Long    = 0x100000100000F
	STRING_LIST_TYPE_SHORT_FORM          Integer = -15
	STRING_LIST_SHORT_FORM               Long    = 0x1000001FFFFF1
	TIME_TYPE_SHORT_FORM                 Integer = 16
	TIME_SHORT_FORM                      Long    = 0x1000001000010
	TIME_LIST_TYPE_SHORT_FORM            Integer = -16
	TIME_LIST_SHORT_FORM                 Long    = 0x1000001FFFFF0
	FINETIME_TYPE_SHORT_FORM             Integer = 17
	FINETIME_SHORT_FORM                  Long    = 0x1000001000011
	FINETIME_LIST_TYPE_SHORT_FORM        Integer = -17
	FINETIME_LIST_SHORT_FORM             Long    = 0x1000001FFFFEF
	URI_TYPE_SHORT_FORM                  Integer = 18
	URI_SHORT_FORM                       Long    = 0x1000001000012
	URI_LIST_TYPE_SHORT_FORM             Integer = -18
	URI_LIST_SHORT_FORM                  Long    = 0x1000001FFFFEE
	INTERACTIONTYPE_TYPE_SHORT_FORM      Integer = 19
	INTERACTIONTYPE_SHORT_FORM           Long    = 0x1000001000013
	INTERACTIONTYPE_LIST_TYPE_SHORT_FORM Integer = -19
	INTERACTIONTYPE_LIST_SHORT_FORM      Long    = 0x1000001FFFFED
	SESSIONTYPE_TYPE_SHORT_FORM          Integer = 20
	SESSIONTYPE_SHORT_FORM               Long    = 0x1000001000014
	SESSIONTYPE_LIST_TYPE_SHORT_FORM     Integer = -20
	SESSIONTYPE_LIST_SHORT_FORM          Long    = 0x1000001FFFFEC
	QOSLEVEL_TYPE_SHORT_FORM             Integer = 21
	QOSLEVEL_SHORT_FORM                  Long    = 0x1000001000015
	QOSLEVEL_LIST_TYPE_SHORT_FORM        Integer = -21
	QOSLEVEL_LIST_SHORT_FORM             Long    = 0x1000001FFFFEB
	UPDATETYPE_TYPE_SHORT_FORM           Integer = 22
	UPDATETYPE_SHORT_FORM                Long    = 0x1000001000016
	UPDATETYPE_LIST_TYPE_SHORT_FORM      Integer = -22
	UPDATETYPE_LIST_SHORT_FORM           Long    = 0x1000001FFFFEA
	SUBSCRIPTION_TYPE_SHORT_FORM         Integer = 23
	SUBSCRIPTION_SHORT_FORM              Long    = 0x1000001000017
	SUBSCRIPTION_LIST_TYPE_SHORT_FORM    Integer = -23
	SUBSCRIPTION_LIST_SHORT_FORM         Long    = 0x1000001FFFFE9
	ENTITYREQUEST_TYPE_SHORT_FORM        Integer = 24
	ENTITYREQUEST_SHORT_FORM             Long    = 0x1000001000018
	ENTITYREQUEST_LIST_TYPE_SHORT_FORM   Integer = -24
	ENTITYREQUEST_LIST_SHORT_FORM        Long    = 0x1000001FFFFE8
	ENTITYKEY_TYPE_SHORT_FORM            Integer = 25
	ENTITYKEY_SHORT_FORM                 Long    = 0x1000001000019
	ENTITYKEY_LIST_TYPE_SHORT_FORM       Integer = -25
	ENTITYKEY_LIST_SHORT_FORM            Long    = 0x1000001FFFFE7
	UPDATEHEADER_TYPE_SHORT_FORM         Integer = 26
	UPDATEHEADER_SHORT_FORM              Long    = 0x100000100001A
	UPDATEHEADER_LIST_TYPE_SHORT_FORM    Integer = -26
	UPDATEHEADER_LIST_SHORT_FORM         Long    = 0x1000001FFFFE6
	IDBOOLEANPAIR_TYPE_SHORT_FORM        Integer = 27
	IDBOOLEANPAIR_SHORT_FORM             Long    = 0x100000100001B
	IDBOOLEANPAIR_LIST_TYPE_SHORT_FORM   Integer = -27
	IDBOOLEANPAIR_LIST_SHORT_FORM        Long    = 0x1000001FFFFE5
	PAIR_TYPE_SHORT_FORM                 Integer = 28
	PAIR_SHORT_FORM                      Long    = 0x100000100001C
	PAIR_LIST_TYPE_SHORT_FORM            Integer = -28
	PAIR_LIST_SHORT_FORM                 Long    = 0x1000001FFFFE4
	NAMEDVALUE_TYPE_SHORT_FORM           Integer = 29
	NAMEDVALUE_SHORT_FORM                Long    = 0x100000100001D
	NAMEDVALUE_LIST_TYPE_SHORT_FORM      Integer = -29
	NAMEDVALUE_LIST_SHORT_FORM           Long    = 0x1000001FFFFE3
	FILE_TYPE_SHORT_FORM                 Integer = 30
	FILE_SHORT_FORM                      Long    = 0x100000100001E
	FILE_LIST_TYPE_SHORT_FORM            Integer = -30
	FILE_LIST_SHORT_FORM                 Long    = 0x1000001FFFFE2

	// errors
	ERROR_DELIVERY_FAILED       UInteger = 65536
	ERROR_DELIVERY_TIMEDOUT     UInteger = 65537
	ERROR_DELIVERY_DELAYED      UInteger = 65538
	ERROR_DESTINATION_UNKNOWN   UInteger = 65539
	ERROR_DESTINATION_TRANSIENT UInteger = 65540
	ERROR_DESTINATION_LOST      UInteger = 65541
	ERROR_AUTHENTICATION_FAIL   UInteger = 65542
	ERROR_AUTHORISATION_FAIL    UInteger = 65543
	ERROR_ENCRYPTION_FAIL       UInteger = 65544
	ERROR_UNSUPPORTED_AREA      UInteger = 65545
	ERROR_UNSUPPORTED_OPERATION UInteger = 65546
	ERROR_UNSUPPORTED_VERSION   UInteger = 65547
	ERROR_BAD_ENCODING          UInteger = 65548
	ERROR_INTERNAL              UInteger = 65549
	ERROR_UNKNOWN               UInteger = 65550
	ERROR_INCORRECT_STATE       UInteger = 65551
	ERROR_TOO_MANY              UInteger = 65552
	ERROR_SHUTDOWN              UInteger = 65553

	ERROR_DELIVERY_FAILED_MESSAGE       String = "Confirmed communication error."
	ERROR_DELIVERY_TIMEDOUT_MESSAGE     String = "Unconfirmed communication error."
	ERROR_DELIVERY_DELAYED_MESSAGE      String = "Message queued somewhere awaiting contact."
	ERROR_DESTINATION_UNKNOWN_MESSAGE   String = "Destination cannot be contacted."
	ERROR_DESTINATION_TRANSIENT_MESSAGE String = "Destination middleware reports destination application does not exist."
	ERROR_DESTINATION_LOST_MESSAGE      String = "Destination lost halfway through conversation."
	ERROR_AUTHENTICATION_FAIL_MESSAGE   String = "A failure to authenticate the message correctly."
	ERROR_AUTHORISATION_FAIL_MESSAGE    String = "A failure in the MAL to authorise the message."
	ERROR_ENCRYPTION_FAIL_MESSAGE       String = "A failure in the MAL to encrypt/decrypt the message."
	ERROR_UNSUPPORTED_AREA_MESSAGE      String = "The destination does not support the service area."
	ERROR_UNSUPPORTED_OPERATION_MESSAGE String = "The destination does not support the operation."
	ERROR_UNSUPPORTED_VERSION_MESSAGE   String = "The destination does not support the area version."
	ERROR_BAD_ENCODING_MESSAGE          String = "The destination was unable to decode the message."
	ERROR_INTERNAL_MESSAGE              String = "An internal error has occurred."
	ERROR_UNKNOWN_MESSAGE               String = "Operation specific."
	ERROR_INCORRECT_STATE_MESSAGE       String = "The destination was not in the correct state for the received message."
	ERROR_TOO_MANY_MESSAGE              String = "Maximum number of subscriptions or providers of a broker has been exceeded."
	ERROR_SHUTDOWN_MESSAGE              String = "The component is being shutdown."
)

// Enumeration values
const (
	INTERACTIONTYPE_SEND_OVAL = iota
	INTERACTIONTYPE_SUBMIT_OVAL
	INTERACTIONTYPE_REQUEST_OVAL
	INTERACTIONTYPE_INVOKE_OVAL
	INTERACTIONTYPE_PROGRESS_OVAL
	INTERACTIONTYPE_PUBSUB_OVAL
)

var (
	// InteractionType is currently defined as a UOctet
	INTERACTIONTYPE_SEND     = InteractionType(UOctet(uint8(INTERACTIONTYPE_SEND_OVAL)))
	INTERACTIONTYPE_SUBMIT   = InteractionType(UOctet(uint8(INTERACTIONTYPE_SUBMIT_OVAL)))
	INTERACTIONTYPE_REQUEST  = InteractionType(UOctet(uint8(INTERACTIONTYPE_REQUEST_OVAL)))
	INTERACTIONTYPE_INVOKE   = InteractionType(UOctet(uint8(INTERACTIONTYPE_INVOKE_OVAL)))
	INTERACTIONTYPE_PROGRESS = InteractionType(UOctet(uint8(INTERACTIONTYPE_PROGRESS_OVAL)))
	INTERACTIONTYPE_PUBSUB   = InteractionType(UOctet(uint8(INTERACTIONTYPE_PUBSUB_OVAL)))
)

const (
	SESSIONTYPE_LIVE_OVAL = iota
	SESSIONTYPE_SIMULATION_OVAL
	SESSIONTYPE_REPLAY_OVAL
)

var (
	// SessionType is currently defined as a UOctet
	SESSIONTYPE_LIVE       = SessionType(UOctet(uint8(SESSIONTYPE_LIVE_OVAL)))
	SESSIONTYPE_SIMULATION = SessionType(UOctet(uint8(SESSIONTYPE_SIMULATION_OVAL)))
	SESSIONTYPE_REPLAY     = SessionType(UOctet(uint8(SESSIONTYPE_REPLAY_OVAL)))
)

const (
	QOSLEVEL_BESTEFFORT_OVAL = iota
	QOSLEVEL_ASSURED_OVAL
	QOSLEVEL_QUEUED_OVAL
	QOSLEVEL_TIMELY_OVAL
)
const (
	// QoSLevel is currently defined as a UOctet
	QOSLEVEL_BESTEFFORT = QoSLevel(UOctet(uint8(QOSLEVEL_BESTEFFORT_OVAL)))
	QOSLEVEL_ASSURED    = QoSLevel(UOctet(uint8(QOSLEVEL_ASSURED_OVAL)))
	QOSLEVEL_QUEUED     = QoSLevel(UOctet(uint8(QOSLEVEL_QUEUED_OVAL)))
	QOSLEVEL_TIMELY     = QoSLevel(UOctet(uint8(QOSLEVEL_TIMELY_OVAL)))
)

const (
	UPDATETYPE_CREATION_OVAL = iota
	UPDATETYPE_UPDATE_OVAL
	UPDATETYPE_MODIFICATION_OVAL
	UPDATETYPE_DELETION_OVAL
)
const (
	// UpdateType is currently defined as a uint8
	UPDATETYPE_CREATION     = UpdateType(uint8(UPDATETYPE_CREATION_OVAL))
	UPDATETYPE_UPDATE       = UpdateType(uint8(UPDATETYPE_UPDATE_OVAL))
	UPDATETYPE_MODIFICATION = UpdateType(uint8(UPDATETYPE_MODIFICATION_OVAL))
	UPDATETYPE_DELETION     = UpdateType(uint8(UPDATETYPE_DELETION_OVAL))
)
