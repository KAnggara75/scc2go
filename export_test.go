// Copyright (c) 2025 KAnggara75
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package scc2go

import "resty.dev/v3"

// SetOsEnvironForTest allows test suites to mock os.Environ.
func SetOsEnvironForTest(fn func() []string) func() {
	orig := osEnviron
	osEnviron = fn
	return func() {
		osEnviron = orig
	}
}

// SetCloseClientForTest allows test suites to mock client.Close behavior.
func SetCloseClientForTest(fn func(client *resty.Client) error) func() {
	orig := closeClient
	closeClient = fn
	return func() {
		closeClient = orig
	}
}
