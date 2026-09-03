// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"
	"net/url"

	"github.com/pkg/errors"

	"github.com/ory/x/logrusx"
)

func Validate(ctx context.Context, l *logrusx.Logger, p *DefaultProvider) error {
	issuerURLs, err := p.IssuerURLs(ctx)
	if err != nil {
		l.WithError(err).Errorf("Configuration key `%s` contains an invalid URL.", KeyIssuerURL)
		return err
	}

	if p.IssuerURL(ctx).String() == "" && !p.IsDevelopmentMode(ctx) {
		l.Errorf("Configuration key `%s` must be set `dev` is `false`. To find out more, use `hydra help serve`.", KeyIssuerURL)
		return errors.New("issuer URL must be set unless development mode is enabled")
	}

	if !p.IsDevelopmentMode(ctx) {
		for _, issuer := range issuerURLs {
			if issuer.Scheme != "https" {
				l.Errorf("Scheme from configuration key `%s` must be `https` when `dev` is `false`. Got scheme in value `%s` is `%s`. To find out more, use `hydra help serve`.", KeyIssuerURL, issuer.String(), issuer.Scheme)
				return errors.New("issuer URL scheme must be HTTPS unless development mode is enabled")
			}
		}
	}

	return nil
}

func urlRoot(u *url.URL) *url.URL {
	if u.Path == "" {
		u.Path = "/"
	}
	return u
}
