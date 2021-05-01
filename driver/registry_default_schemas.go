package driver

import (
	"context"
	"net/url"

	"kratos/schema"
)

func (m *RegistryDefault) IdentityTraitsSchemas(ctx context.Context) schema.Schemas {
	ms := m.Config(ctx).IdentityTraitsSchemas()
	var ss schema.Schemas

	for _, s := range ms {
		surl, err := url.Parse(s.URL)
		if err != nil {
			m.l.Fatalf("Could not parse url %s for schema %s", s.URL, s.ID)
		}

		ss = append(ss, schema.Schema{
			ID:     s.ID,
			URL:    surl,
			RawURL: s.URL,
		})
	}

	return ss
}
