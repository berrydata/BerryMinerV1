package tracker

import (
	"fmt"
	"github.com/berry-data/BerryMiner/apiOracle"
	"github.com/berry-data/BerryMiner/config"
	"github.com/berry-data/BerryMiner/util"
	"os"
	"testing"
)

var configJSON = `{
	"publicAddress":"0000000000000000000000000000000000000000",
	"privateKey":"1111111111111111111111111111111111111111111111111111111111111111",
	"contractAddress":"0x724D1B69a7Ba352F11D73fDBdEB7fF869cB22E19",
	"psrFolder": ".."
}
`

func TestMain(m *testing.M) {
	err := config.ParseConfigBytes([]byte(configJSON))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse mock config: %v\n", err)
		os.Exit(-1)
	}
	util.ParseLoggingConfig("")
	apiOracle.EnsureValueOracle()
	os.Exit(m.Run())
}

