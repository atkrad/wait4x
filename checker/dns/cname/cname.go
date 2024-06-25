// Copyright 2023 The Wait4X Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cname

import (
	"context"
	"net"
	"regexp"

	"wait4x.dev/v2/checker"
)

// Option configures an DNS CNAME record
type Option func(d *CNAME)

// CNAME represents DNS CNAME data structure
type CNAME struct {
	nameserver      string
	address         string
	expectedDomains []string
	resolver        *net.Resolver
}

// New creates the DNS CNAME checker
func New(address string, opts ...Option) checker.Checker {
	d := &CNAME{
		address:  address,
		resolver: net.DefaultResolver,
	}

	// apply the list of options to CNAME
	for _, opt := range opts {
		opt(d)
	}

	// Nameserver settings.
	if d.nameserver != "" {
		d.resolver = &net.Resolver{
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				dialer := net.Dialer{}
				return dialer.DialContext(ctx, network, d.nameserver)
			},
		}
	}

	return d
}

// WithNameServer overrides the default nameserver
func WithNameServer(nameserver string) Option {
	return func(d *CNAME) {
		d.nameserver = nameserver
	}
}

// WithExpectedDomains sets expected domains
func WithExpectedDomains(doamins []string) Option {
	return func(d *CNAME) {
		d.expectedDomains = doamins
	}
}

// Identity returns the identity of the checker
func (d *CNAME) Identity() (string, error) {
	return d.address, nil
}

// Check checks DNS TXT records
func (d *CNAME) Check(ctx context.Context) (err error) {
	value, err := d.resolver.LookupCNAME(ctx, d.address)
	if err != nil {
		return err
	}

	if len(value) != 0 && len(d.expectedDomains) == 0 {
		return nil
	}
	for _, expectedDomain := range d.expectedDomains {
		matched, _ := regexp.MatchString(expectedDomain, value)
		if matched {
			return nil
		}
	}

	return checker.NewExpectedError(
		"the CNAME record value doesn't expect", nil,
		"actual", value, "expect", d.expectedDomains,
	)
}
