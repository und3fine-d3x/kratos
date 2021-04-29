package node

import "kratos/x"

func NewCSRFNode(token string) *Node {
	return &Node{
		Type:  Input,
		Group: DefaultGroup,
		Attributes: &InputAttributes{
			Name:       x.CSRFTokenName,
			Type:       InputAttributeTypeHidden,
			FieldValue: token,
			Required:   true,
		},
	}
}
