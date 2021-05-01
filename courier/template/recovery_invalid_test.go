package template_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kratos/courier/template"
	"kratos/internal"
)

func TestRecoverInvalid(t *testing.T) {
	conf, _ := internal.NewFastRegistryWithMocks(t)
	tpl := template.NewRecoveryInvalid(conf, &template.RecoveryInvalidModel{})

	rendered, err := tpl.EmailBody()
	require.NoError(t, err)
	assert.NotEmpty(t, rendered)

	rendered, err = tpl.EmailSubject()
	require.NoError(t, err)
	assert.NotEmpty(t, rendered)
}
