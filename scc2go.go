// Copyright (c) 2025 KAnggara75
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// @author KAnggara75 on Sun 27/04/25 00.20
// @project scc2go scc2go
//
// https://github.com/KAnggara75/scc2go
//

package scc2go

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"resty.dev/v3"
)

type springCloudConfig struct {
	Name            string           `json:"name"`
	Profiles        []string         `json:"profiles"`
	Label           string           `json:"label"`
	Version         string           `json:"version"`
	State           string           `json:"state"`
	PropertySources []propertySource `json:"propertySources"`
}

type propertySource struct {
	Name   string         `json:"name"`
	Source map[string]any `json:"source"`
}

func GetEnv(sccUrl, auth string, disableTlsOpt ...bool) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	disableTls := false
	if len(disableTlsOpt) > 0 {
		disableTls = disableTlsOpt[0]
	}

	if sccUrl == "" || sccUrl == "local" {
		log.Info().Msg("SCC URL is local or empty, loading from environment variables")
		loadFromEnv()
		return
	}

	log.Info().Msgf("Using SCC URL: %s", sccUrl)
	resBody, err := getSCC(sccUrl, auth, disableTls)
	if err != nil {
		log.Error().Msgf("Error when get scc: %v", err)
		return
	}

	scc := new(springCloudConfig)
	err = json.Unmarshal(resBody, scc)
	if err != nil {
		log.Error().Msgf("Spring Cloud Config: %v", err)
		return
	}

	for i := len(scc.PropertySources) - 1; i >= 0; i-- {
		for key, value := range scc.PropertySources[i].Source {
			log.Trace().Msgf("Retrive %s", key)
			setIfNotExists(key, value)
		}
	}
}

// loadFromEnv reads all OS environment variables and stores them in viper.
// Each env var key is transformed: underscores (_) become dots (.)
// and the key is lowercased so that EXAMPLE_VAR becomes example.var.
func loadFromEnv() {
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		rawKey, value := parts[0], parts[1]
		viperKey := strings.ToLower(strings.ReplaceAll(rawKey, "_", "."))
		log.Trace().Msgf("Loading env var %s as %s", rawKey, viperKey)
		setIfNotExists(viperKey, value)
	}
}

func setIfNotExists(k string, v any) {
	if viper.IsSet(k) {
		return
	}
	viper.Set(k, v)
}

func getSCC(url, authHeader string, disableTls bool) ([]byte, error) {
	tlsConfig := &tls.Config{}
	if disableTls {
		tlsConfig = &tls.Config{InsecureSkipVerify: true} // #nosec G402 -- caller explicitly opted in
	}

	client := resty.New().
		SetTimeout(5 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(time.Second).
		SetTLSClientConfig(tlsConfig)
	defer func(client *resty.Client) {
		err := client.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(client)

	res, err := client.R().
		SetHeader("Authorization", authHeader).
		Get(url)
	if err != nil {
		return nil, fmt.Errorf("fail get config from %s with error: %v", url, err)
	}

	if res.IsError() {
		return nil, fmt.Errorf("fail get config from %s with error: %s", url, res.Status())
	}

	return res.Bytes(), nil
}
