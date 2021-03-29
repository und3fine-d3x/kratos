package link_test

import (
	"github.com/ory/viper"

	"kratos/driver/configuration"
	"kratos/identity"
	"kratos/selfservice/flow/recovery"
)

func initViper() {
	viper.Set(configuration.ViperKeyDefaultIdentitySchemaURL, "file://./stub/default.schema.json")
	viper.Set(configuration.ViperKeySelfServiceBrowserDefaultReturnTo, "https://www.ory.sh")
	viper.Set(configuration.ViperKeySelfServiceStrategyConfig+"."+identity.CredentialsTypePassword.String()+".enabled", true)
	viper.Set(configuration.ViperKeySelfServiceStrategyConfig+"."+recovery.StrategyRecoveryLinkName+".enabled", true)
	viper.Set(configuration.ViperKeySelfServiceRecoveryEnabled, true)
	viper.Set(configuration.ViperKeySelfServiceVerificationEnabled, true)
}
