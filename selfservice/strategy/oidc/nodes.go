package oidc

import (
	"kratos/text"
	"kratos/ui/node"
)

func NewLinkNode(provider string) *node.Node {
	return node.NewInputField("link", provider, node.OpenIDConnectGroup, node.InputAttributeTypeSubmit).WithMetaLabel(text.NewInfoSelfServiceSettingsUpdateLinkOIDC(provider))
}

func NewUnlinkNode(provider string) *node.Node {
	return node.NewInputField("unlink", provider, node.OpenIDConnectGroup, node.InputAttributeTypeSubmit).WithMetaLabel(text.NewInfoSelfServiceSettingsUpdateUnlinkOIDC(provider))
}
