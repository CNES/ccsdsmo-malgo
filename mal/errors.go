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

const (
	MAL_ERROR_DELIVERY_FAILED       UInteger = 65536
	MAL_ERROR_DELIVERY_TIMEDOUT     UInteger = 65537
	MAL_ERROR_DELIVERY_DELAYED      UInteger = 65538
	MAL_ERROR_DESTINATION_UNKNOWN   UInteger = 65539
	MAL_ERROR_DESTINATION_TRANSIENT UInteger = 65540
	MAL_ERROR_DESTINATION_LOST      UInteger = 65541
	MAL_ERROR_AUTHENTICATION_FAIL   UInteger = 65542
	MAL_ERROR_AUTHORISATION_FAIL    UInteger = 65543
	MAL_ERROR_ENCRYPTION_FAIL       UInteger = 65544
	MAL_ERROR_UNSUPPORTED_AREA      UInteger = 65545
	MAL_ERROR_UNSUPPORTED_OPERATION UInteger = 65546
	MAL_ERROR_UNSUPPORTED_VERSION   UInteger = 65547
	MAL_ERROR_BAD_ENCODING          UInteger = 65548
	MAL_ERROR_INTERNAL              UInteger = 65549
	MAL_ERROR_UNKNOWN               UInteger = 65550
	MAL_ERROR_INCORRECT_STATE       UInteger = 65551
	MAL_ERROR_TOO_MANY              UInteger = 65552
	MAL_ERROR_SHUTDOWN              UInteger = 65553

	MAL_ERROR_DELIVERY_FAILED_MESSAGE       String = "Confirmed communication error."
	MAL_ERROR_DELIVERY_TIMEDOUT_MESSAGE     String = "Unconfirmed communication error."
	MAL_ERROR_DELIVERY_DELAYED_MESSAGE      String = "Message queued somewhere awaiting contact."
	MAL_ERROR_DESTINATION_UNKNOWN_MESSAGE   String = "Destination cannot be contacted."
	MAL_ERROR_DESTINATION_TRANSIENT_MESSAGE String = "Destination middleware reports destination application does not exist."
	MAL_ERROR_DESTINATION_LOST_MESSAGE      String = "Destination lost halfway through conversation."
	MAL_ERROR_AUTHENTICATION_FAIL_MESSAGE   String = "A failure to authenticate the message correctly."
	MAL_ERROR_AUTHORISATION_FAIL_MESSAGE    String = "A failure in the MAL to authorise the message."
	MAL_ERROR_ENCRYPTION_FAIL_MESSAGE       String = "A failure in the MAL to encrypt/decrypt the message."
	MAL_ERROR_UNSUPPORTED_AREA_MESSAGE      String = "The destination does not support the service area."
	MAL_ERROR_UNSUPPORTED_OPERATION_MESSAGE String = "The destination does not support the operation."
	MAL_ERROR_UNSUPPORTED_VERSION_MESSAGE   String = "The destination does not support the area version."
	MAL_ERROR_BAD_ENCODING_MESSAGE          String = "The destination was unable to decode the message."
	MAL_ERROR_INTERNAL_MESSAGE              String = "An internal error has occurred."
	MAL_ERROR_UNKNOWN_MESSAGE               String = "Operation specific."
	MAL_ERROR_INCORRECT_STATE_MESSAGE       String = "The destination was not in the correct state for the received message."
	MAL_ERROR_TOO_MANY_MESSAGE              String = "Maximum number of subscriptions or providers of a broker has been exceeded."
	MAL_ERROR_SHUTDOWN_MESSAGE              String = "The component is being shutdown."
)

// TODO (AF): Defines a map allowing to get message from error code.
