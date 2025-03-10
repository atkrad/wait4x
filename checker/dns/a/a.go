// Copyright 2019-2025 The Wait4X Authors
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

// Package a provides functionality for checking the A records of a domain.
package a

import (
	"context"
	"net"

	"wait4x.dev/v3/checker"
)

// Option configures an DNS A records
type Option func(d *A)

// A represents DNS A data structure
type A struct {
	nameserver  string
	address     string
	expectedIPs []string
	resolver    *net.Resolver
}

// New creates the DNS A checker
func New(address string, opts ...Option) checker.Checker {
	d := &A{
		address:  address,
		resolver: net.DefaultResolver,
	}

	// apply the list of options to A
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
	return func(d *A) {
		d.nameserver = nameserver
	}
}

// WithExpectedIPV4s sets expected IPv4s
func WithExpectedIPV4s(ips []string) Option {
	return func(d *A) {
		d.expectedIPs = ips
	}
}

// Identity returns the identity of the checker
func (d *A) Identity() (string, error) {
	return d.address, nil
}

// Check checks A DNS records
func (d *A) Check(ctx context.Context) (err error) {
	ips, err := d.resolver.LookupIP(ctx, "ip4", d.address)
	if err != nil {
		return err
	}

	for _, ip := range ips {
		if len(d.expectedIPs) == 0 {
			return nil
		}
		for _, expectedIP := range d.expectedIPs {
			if expectedIP == ip.String() {
				return nil
			}
		}
	}

	return checker.NewExpectedError(
		"the A record value doesn't expect", nil,
		"actual", ips, "expect", d.expectedIPs,
	)
}
