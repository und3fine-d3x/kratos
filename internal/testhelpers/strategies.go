package testhelpers

import (
	"fmt"
	"testing"

	"kratos/driver/config"
)

func StrategyEnable(t *testing.T, c *config.Config, strategy string, enable bool) {
	c.MustSet(fmt.Sprintf("%s.%s.enabled", config.ViperKeySelfServiceStrategyConfig, strategy), enable)
}
