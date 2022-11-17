// Copyright 2022 Mohammad Abdolirad
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
	"context"
	"fmt"
	"net"
	"regexp"

	"github.com/atkrad/wait4x/v2/pkg/checker"
	"github.com/atkrad/wait4x/v2/pkg/checker/errors"
)

// Option configures an DNS.
type Option func(d *DNS)

type RecordType uint8

const (
	A RecordType = iota
	AAAA
	CNAME
	MX
	TXT
	NS
)

// DNS data structure.
type DNS struct {
	recordType    RecordType
	nameserver    string
	address       string
	expectedValue string
	resolver      *net.Resolver
}

func (rt RecordType) String() string {
	switch rt {
	case A:
		return "A"
	case AAAA:
		return "AAAA"
	case CNAME:
		return "CNAME"
	case MX:
		return "MX"
	case TXT:
		return "TXT"
	case NS:
		return "NS"
	}
	return ""
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

	// Nameserver settings.
	d.resolver = net.DefaultResolver
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

func WithNameServer(nameserver string) Option {
	return func(d *DNS) {
		d.nameserver = nameserver
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
		return errors.Wrap(err, errors.DebugLevel)
	}

	for _, ip := range ips {
		if expectedValue == ip.String() {
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
		return errors.Wrap(err, errors.DebugLevel)
	}

	for _, ip := range values {
		if expectedValue == ip.String() {
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
		return errors.Wrap(err, errors.DebugLevel)
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
		return errors.Wrap(err, errors.DebugLevel)
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
		return errors.Wrap(err, errors.DebugLevel)
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
		return errors.Wrap(err, errors.DebugLevel)
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
	return fmt.Sprintf("%s %s %s", d.recordType, d.address, d.expectedValue), nil
}

// Check checks DNS records
func (d *DNS) Check(ctx context.Context) (err error) {
	switch d.recordType {
	case A:
		return d.CheckARecords(ctx, d.address, d.expectedValue)
	case AAAA:
		return d.CheckAAAARecords(ctx, d.address, d.expectedValue)
	case CNAME:
		return d.CheckCNAMERecord(ctx, d.address, d.expectedValue)
	case MX:
		return d.CheckMXRecords(ctx, d.address, d.expectedValue)
	case TXT:
		return d.CheckTXTRecords(ctx, d.address, d.expectedValue)
	case NS:
		return d.CheckNSRecords(ctx, d.address, d.expectedValue)
	}
	return nil
}
