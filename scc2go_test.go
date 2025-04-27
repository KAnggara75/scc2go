// Copyright (c) 2025 KAnggara75
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// @author KAnggara75 on Sun 27/04/25 00.23
// @project scc2go scc2go
//

package scc2go_test

import (
	"github.com/KAnggara75/scc2go"
	"testing"
)

func TestGetEnvNoUrl(t *testing.T) {
	scc2go.GetEnv("", "")
}

func TestGetEnvInvalidUrl(t *testing.T) {
	scc2go.GetEnv("coba", "")
}

func TestGetEnvInvalidRes(t *testing.T) {
	scc2go.GetEnv("https://google.com", "")
}

func TestGetEnvInvalidAuth(t *testing.T) {
	scc2go.GetEnv("https://httpbin.org/bearer", "")
}

func TestGetEnvSuccess(t *testing.T) {
	scc2go.GetEnv("https://mdlwg.wiremockapi.cloud/test/prd", "Bearer coba")
}
