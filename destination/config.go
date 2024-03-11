package destination

import "github.com/mer-oscar/conduit-connector-zeromq/common"

//go:generate paramgen -output=config_paramgen.go Config
type Config struct {
	common.Config
}
