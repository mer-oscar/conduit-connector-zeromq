package zeromq

import (
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/mer-oscar/conduit-connector-zeromq/destination"
	"github.com/mer-oscar/conduit-connector-zeromq/source"
)

// Connector combines all constructors for each plugin in one struct.
var Connector = sdk.Connector{
	NewSpecification: Specification,
	NewSource:        source.New,
	NewDestination:   destination.New,
}
