package archive

import (
	"github.com/CNES/ccsdsmo-malgo/mal"
)

// Defines the abstract composite interfaces.
type QueryFilter interface {
	mal.Composite
	QueryFilter() QueryFilter
}

var NullQueryFilter QueryFilter = nil

type QueryFilterList interface {
	mal.ElementList
	QueryFilterList() QueryFilterList
}

var NullQueryFilterList QueryFilterList = nil
