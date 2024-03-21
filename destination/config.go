package destination

import "github.com/mer-oscar/conduit-connector-zeromq/common"

//go:generate paramgen -output=config_paramgen.go Config
type Config struct {
	common.Config
	// RouterEndpoints is a comma separated list of socket endpoints that we wish to deal messages to
	RouterEndpoints string `json:"routerEndpoints" validate:"required"`
}
