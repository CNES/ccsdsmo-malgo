package com

import (
	"github.com/CNES/ccsdsmo-malgo/mal"
)

const (
	AREA_NUMBER     mal.UShort   = 2
	AREA_VERSION    mal.UOctet   = 1
	AREA_NAME                    = mal.Identifier("COM")
	ERROR_INVALID   mal.UInteger = 70000
	ERROR_DUPLICATE mal.UInteger = 70001
)
