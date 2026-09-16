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
	"context"
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

// ConfigTarget is an abstraction for setting key-values into Viper instances.
type ConfigTarget interface {
	IsSet(key string) bool
	Set(key string, value any)
}

type globalViperTarget struct{}

func (globalViperTarget) IsSet(key string) bool {
	return viper.IsSet(key)
}

func (globalViperTarget) Set(key string, value any) {
	viper.Set(key, value)
}

// Option configures behavior of SCC loading.
type Option func(*loaderConfig)

type loaderConfig struct {
	ctx        context.Context
	target     ConfigTarget
	debug      bool
	disableTls bool
	timeout    time.Duration
}

// WithContext supplies a context for cancellation and timeouts.
func WithContext(ctx context.Context) Option {
	return func(c *loaderConfig) {
		if ctx != nil {
			c.ctx = ctx
		}
	}
}

// WithViper specifies a custom *viper.Viper instance to populate instead of the global singleton.
func WithViper(v *viper.Viper) Option {
	return func(c *loaderConfig) {
		if v != nil {
			c.target = v
		}
	}
}

// WithDebug enables or disables trace-level logging.
func WithDebug(debug bool) Option {
	return func(c *loaderConfig) {
		c.debug = debug
	}
}

// WithDisableTLS skips TLS certificate verification.
func WithDisableTLS(disable bool) Option {
	return func(c *loaderConfig) {
		c.disableTls = disable
	}
}

// WithTimeout sets a custom HTTP request timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *loaderConfig) {
		if timeout > 0 {
			c.timeout = timeout
		}
	}
}

// Load fetches configurations and stores them into the designated target, returning an error if retrieval or unmarshaling fails.
func Load(sccUrl, auth string, opts ...Option) error {
	cfg := &loaderConfig{
		ctx:        context.Background(),
		target:     globalViperTarget{},
		debug:      false,
		disableTls: false,
		timeout:    5 * time.Second,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(cfg)
		}
	}

	level := zerolog.InfoLevel
	if cfg.debug {
		level = zerolog.TraceLevel
	}

	logger := zerolog.New(os.Stdout).
		Level(level).
		With().
		Timestamp().
		Str("component", "scc_loader").
		Logger()

	if sccUrl == "" || sccUrl == "local" {
		logger.Info().Msg("SCC URL is local or empty, loading from environment variables")
		loadFromEnvToTarget(cfg.target)
		return nil
	}

	logger.Info().
		Str("scc_url", sccUrl).
		Msg("using SCC URL")

	resBody, err := getSCCWithContext(cfg.ctx, sccUrl, auth, cfg.disableTls, cfg.timeout)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("error when get scc")
		return err
	}

	var scc springCloudConfig
	if err := json.Unmarshal(resBody, &scc); err != nil {
		logger.Error().
			Err(err).
			Msg("spring cloud config unmarshal failed")
		return fmt.Errorf("spring cloud config unmarshal failed: %w", err)
	}

	for i := len(scc.PropertySources) - 1; i >= 0; i-- {
		for key, value := range scc.PropertySources[i].Source {
			logger.Trace().
				Str("key", key).
				Msg("retrieve property")
			setIfNotExistsOnTarget(cfg.target, key, value)
		}
	}

	return nil
}

// GetEnv preserves backward compatibility by populating the global viper singleton without returning an error.
func GetEnv(sccUrl, auth string, disableTlsOpt ...bool) {
	GetEnvWithDebug(sccUrl, auth, false, disableTlsOpt...)
}

// GetEnvWithDebug preserves backward compatibility with debug logging flag.
func GetEnvWithDebug(sccUrl, auth string, debug bool, disableTlsOpt ...bool) {
	disableTls := false
	if len(disableTlsOpt) > 0 {
		disableTls = disableTlsOpt[0]
	}

	_ = Load(sccUrl, auth,
		WithDebug(debug),
		WithDisableTLS(disableTls),
	)
}

// loadFromEnv reads all OS environment variables and stores them in the global viper.
func loadFromEnv() {
	loadFromEnvToTarget(globalViperTarget{})
}

func loadFromEnvToTarget(target ConfigTarget) {
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		rawKey, value := parts[0], parts[1]
		viperKey := strings.ToLower(strings.ReplaceAll(rawKey, "_", "."))
		log.Trace().Msgf("Loading env var %s as %s", rawKey, viperKey)
		setIfNotExistsOnTarget(target, viperKey, value)
	}
}

func setIfNotExists(k string, v any) {
	setIfNotExistsOnTarget(globalViperTarget{}, k, v)
}

func setIfNotExistsOnTarget(target ConfigTarget, k string, v any) {
	if target.IsSet(k) {
		return
	}
	target.Set(k, v)
}

func getSCC(url, authHeader string, disableTls bool) ([]byte, error) {
	return getSCCWithContext(context.Background(), url, authHeader, disableTls, 5*time.Second)
}

func getSCCWithContext(ctx context.Context, url, authHeader string, disableTls bool, timeout time.Duration) ([]byte, error) {
	tlsConfig := &tls.Config{}
	if disableTls {
		tlsConfig = &tls.Config{InsecureSkipVerify: true} // #nosec G402 -- caller explicitly opted in
	}

	client := resty.New().
		SetTimeout(timeout).
		SetRetryCount(3).
		SetRetryWaitTime(time.Second).
		SetTLSClientConfig(tlsConfig)
	defer func(client *resty.Client) {
		err := client.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(client)

	req := client.R()
	if ctx != nil {
		req.SetContext(ctx)
	}
	if authHeader != "" {
		req.SetHeader("Authorization", authHeader)
	}

	res, err := req.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fail get config from %s with error: %v", url, err)
	}

	if res.IsStatusFailure() {
		return nil, fmt.Errorf("fail get config from %s with error: %s", url, res.Status())
	}

	return res.Bytes(), nil
}
