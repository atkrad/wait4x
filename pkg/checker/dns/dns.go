// Copyright 2020 Mohammad Abdolirad
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

package dns

import (
	"github.com/atkrad/wait4x/v2/pkg/checker"
	"github.com/atkrad/wait4x/v2/pkg/checker/errors"
	"net"
	"regexp"
)
import "context"

// Option configures an DNS.
type Option func(d *DNS)

type RecordType = uint8

const (
	A     RecordType = 1
	AAAA             = 2
	CNAME            = 3
	MX               = 4
	TXT              = 5
	NS               = 6
)

// DNS data structure.
type DNS struct {
	recordType    RecordType
	nameserver    string
	address       string
	expectedValue string
	resolver      *net.Resolver
}

// New creates the DNS checker
func New(recordType RecordType, address string, opts ...Option) checker.Checker {
	d := &DNS{
		recordType: recordType,
		address:    address,
	}

	// apply the list of options to HTTP
	for _, opt := range opts {
		opt(d)
	}

	return d
}

func WithNameServer(nameserver string) Option {
	return func(d *DNS) {
		d.nameserver = nameserver

		// Nameserver settings.
		if d.nameserver != "" {
			d.resolver = &net.Resolver{
				Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
					dialer := net.Dialer{}
					return dialer.DialContext(ctx, network, d.nameserver)
				},
			}
		} else {
			d.resolver = &net.Resolver{
				PreferGo: true,
			}
		}
	}
}

func WithExpectedValue(value string) Option {
	return func(d *DNS) {
		d.expectedValue = value
	}
}

func (d *DNS) CheckARecords(ctx context.Context, address string, expectedValue string) (err error) {
	ips, err := d.resolver.LookupIP(ctx, "ip4", address)
	if err != nil {
		return errors.New(
			"cannot get A records",
			errors.InfoLevel,
			errors.WithFields("actual", err.Error(), "expect", expectedValue),
		)
	}
	for _, ip := range ips {
		matched, _ := regexp.MatchString(expectedValue, ip.String())
		if matched {
			return nil
		}
	}
	return errors.New(
		"the A record value doesn't expect",
		errors.InfoLevel,
		errors.WithFields("actual", ips, "expect", expectedValue),
	)
}

func (d *DNS) CheckAAAARecords(ctx context.Context, address string, expectedValue string) (err error) {
	values, err := d.resolver.LookupIP(ctx, "ip6", address)
	if err != nil {
		return errors.New(
			"cannot get AAAA records",
			errors.InfoLevel,
			errors.WithFields("actual", err.Error(), "expect", expectedValue),
		)
	}
	for _, ip := range values {
		matched, _ := regexp.MatchString(expectedValue, ip.String())
		if matched {
			return nil
		}
	}
	return errors.New(
		"the AAAA record value doesn't expect",
		errors.InfoLevel,
		errors.WithFields("actual", values, "expect", expectedValue),
	)
}

func (d *DNS) CheckCNAMERecord(ctx context.Context, address string, expectedValue string) (err error) {
	value, err := d.resolver.LookupCNAME(ctx, address)
	if err != nil {
		return errors.New(
			"cannot get CNAME record",
			errors.InfoLevel,
			errors.WithFields("actual", err.Error(), "expect", expectedValue),
		)
	}
	matched, _ := regexp.MatchString(expectedValue, value)
	if matched {
		return nil
	}
	return errors.New(
		"the CNAME record value doesn't expect",
		errors.InfoLevel,
		errors.WithFields("actual", value, "expect", expectedValue),
	)
}

func (d *DNS) CheckMXRecords(ctx context.Context, address string, expectedValue string) (err error) {
	values, err := d.resolver.LookupMX(ctx, address)
	if err != nil {
		return errors.New(
			"cannot get MX record",
			errors.InfoLevel,
			errors.WithFields("actual", err.Error(), "expect", expectedValue),
		)
	}
	for _, mx := range values {
		matched, _ := regexp.MatchString(expectedValue, mx.Host)
		if matched {
			return nil
		}
	}
	return errors.New(
		"the MX record value doesn't expect",
		errors.InfoLevel,
		errors.WithFields("actual", values, "expect", expectedValue),
	)
}

func (d *DNS) CheckTXTRecords(ctx context.Context, address string, expectedValue string) (err error) {
	values, err := d.resolver.LookupTXT(ctx, address)
	if err != nil {
		return errors.New(
			"cannot get TXT record",
			errors.InfoLevel,
			errors.WithFields("actual", err.Error(), "expect", expectedValue),
		)
	}
	for _, txt := range values {
		matched, _ := regexp.MatchString(expectedValue, txt)
		if matched {
			return nil
		}
	}
	return errors.New(
		"the TXT record value doesn't expect",
		errors.InfoLevel,
		errors.WithFields("actual", values, "expect", expectedValue),
	)
}

func (d *DNS) CheckNSRecords(ctx context.Context, address string, expectedValue string) (err error) {
	values, err := d.resolver.LookupNS(ctx, address)
	if err != nil {
		return errors.New(
			"cannot get NS record",
			errors.InfoLevel,
			errors.WithFields("actual", err.Error(), "expect", expectedValue),
		)
	}
	for _, ns := range values {
		matched, _ := regexp.MatchString(expectedValue, ns.Host)
		if matched {
			return nil
		}
	}
	return errors.New(
		"the NS record value doesn't expect",
		errors.InfoLevel,
		errors.WithFields("actual", values, "expect", expectedValue),
	)
}

// Identity returns the identity of the checker
func (d *DNS) Identity() (string, error) {
	return d.nameserver, nil
}

// Check checks DNS records
func (d *DNS) Check(ctx context.Context) (err error) {
	if d.recordType == A {
		return d.CheckARecords(ctx, d.address, d.expectedValue)
	} else if d.recordType == AAAA {
		return d.CheckAAAARecords(ctx, d.address, d.expectedValue)
	} else if d.recordType == CNAME {
		return d.CheckCNAMERecord(ctx, d.address, d.expectedValue)
	} else if d.recordType == MX {
		return d.CheckMXRecords(ctx, d.address, d.expectedValue)
	} else if d.recordType == TXT {
		return d.CheckTXTRecords(ctx, d.address, d.expectedValue)
	} else if d.recordType == NS {
		return d.CheckNSRecords(ctx, d.address, d.expectedValue)
	}
	return nil
}
