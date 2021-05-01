package link_test

import (
	"testing"

	"kratos/driver/config"
	"kratos/identity"
	"kratos/selfservice/flow/recovery"
)

func initViper(t *testing.T, c *config.Config) {
	c.MustSet(config.ViperKeyDefaultIdentitySchemaURL, "file://./stub/default.schema.json")
	c.MustSet(config.ViperKeySelfServiceBrowserDefaultReturnTo, "https://www.ory.sh")
	c.MustSet(config.ViperKeySelfServiceStrategyConfig+"."+identity.CredentialsTypePassword.String()+".enabled", true)
	c.MustSet(config.ViperKeySelfServiceStrategyConfig+"."+recovery.StrategyRecoveryLinkName+".enabled", true)
	c.MustSet(config.ViperKeySelfServiceRecoveryEnabled, true)
	c.MustSet(config.ViperKeySelfServiceVerificationEnabled, true)
}
