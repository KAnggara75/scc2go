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
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	auth := os.Getenv("AUTH")
	scc2go.GetEnv(auth)
	assert.NotNil(t, viper.Get("db.pakaiwa.host"))
}
