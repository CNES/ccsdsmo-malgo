package event

import (
	"github.com/CNES/ccsdsmo-malgo/mal"
)

const (
	// standard service identifiers
	SERVICE_NUMBER mal.UShort = 1
	SERVICE_NAME              = mal.Identifier("Event")

	// standard operation identifiers
	MONITOREVENT_OPERATION_NUMBER mal.UShort = 1
)
