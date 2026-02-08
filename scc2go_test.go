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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KAnggara75/scc2go"
	"github.com/spf13/viper"
)

func TestGetEnvNoUrl(t *testing.T) {
	// Test with empty URL - should log error and return early
	scc2go.GetEnv("", "")
	scc2go.GetEnv("", "")
}

func TestGetEnvInvalidUrl(t *testing.T) {
	// Test with invalid URL format
	scc2go.GetEnv("coba", "")
}

func TestGetEnvInvalidRes(t *testing.T) {
	// Test with valid URL but invalid response format
	scc2go.GetEnv("https://google.com", "")
}

func TestGetEnvInvalidAuth(t *testing.T) {
	// Test with URL requiring auth but no auth provided
	scc2go.GetEnv("https://httpbin.org/bearer", "")
}

func TestGetEnvSuccess(t *testing.T) {
	// Test with valid Spring Cloud Config response
	scc2go.GetEnv("https://mdlwg.wiremockapi.cloud/test/prd", "Bearer coba")
}

func TestGetEnvWithDisableTLS(t *testing.T) {
	// Test with TLS disabled
	scc2go.GetEnv("https://mdlwg.wiremockapi.cloud/test/prd", "Bearer coba", true)
}

func TestGetEnvWithMockServer(t *testing.T) {
	tests := []struct {
		name          string
		setupServer   func() *httptest.Server
		auth          string
		disableTLS    bool
		expectError   bool
		setupViper    func()
		validateViper func(t *testing.T)
	}{
		{
			name: "successful config retrieval",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					response := map[string]any{
						"name":     "test-app",
						"profiles": []string{"default"},
						"label":    "main",
						"version":  "1.0.0",
						"state":    "active",
						"propertySources": []map[string]any{
							{
								"name": "test-source",
								"source": map[string]any{
									"app.name":    "test-application",
									"app.version": "1.0.0",
									"app.port":    8080,
								},
							},
						},
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(response)
				}))
			},
			auth:        "Bearer test-token",
			disableTLS:  false,
			expectError: false,
			setupViper: func() {
				viper.Reset()
			},
			validateViper: func(t *testing.T) {
				if !viper.IsSet("app.name") {
					t.Error("Expected app.name to be set in viper")
				}
				if viper.GetString("app.name") != "test-application" {
					t.Errorf("Expected app.name to be 'test-application', got '%s'", viper.GetString("app.name"))
				}
			},
		},
		{
			name: "multiple property sources with override",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					response := map[string]any{
						"name":     "test-app",
						"profiles": []string{"default", "prod"},
						"label":    "main",
						"version":  "1.0.0",
						"state":    "active",
						"propertySources": []map[string]any{
							{
								"name": "default-source",
								"source": map[string]any{
									"app.name": "default-app",
									"app.env":  "default",
								},
							},
							{
								"name": "prod-source",
								"source": map[string]any{
									"app.name": "prod-app",
									"app.env":  "production",
								},
							},
						},
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(response)
				}))
			},
			auth:        "Bearer test-token",
			disableTLS:  false,
			expectError: false,
			setupViper: func() {
				viper.Reset()
			},
			validateViper: func(t *testing.T) {
				// Property sources are processed in reverse order, so prod-source should be applied first
				// Then default-source, but it won't override existing values
				if viper.GetString("app.name") != "prod-app" {
					t.Errorf("Expected app.name to be 'prod-app', got '%s'", viper.GetString("app.name"))
				}
			},
		},
		{
			name: "viper already has value - should not override",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					response := map[string]any{
						"name":     "test-app",
						"profiles": []string{"default"},
						"label":    "main",
						"version":  "1.0.0",
						"state":    "active",
						"propertySources": []map[string]any{
							{
								"name": "test-source",
								"source": map[string]any{
									"app.name": "from-scc",
								},
							},
						},
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(response)
				}))
			},
			auth:        "Bearer test-token",
			disableTLS:  false,
			expectError: false,
			setupViper: func() {
				viper.Reset()
				viper.Set("app.name", "existing-value")
			},
			validateViper: func(t *testing.T) {
				if viper.GetString("app.name") != "existing-value" {
					t.Errorf("Expected app.name to remain 'existing-value', got '%s'", viper.GetString("app.name"))
				}
			},
		},
		{
			name: "server returns 401 unauthorized",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unauthorized"))
				}))
			},
			auth:        "",
			disableTLS:  false,
			expectError: true,
			setupViper: func() {
				viper.Reset()
			},
			validateViper: func(t *testing.T) {
				// No values should be set on error
			},
		},
		{
			name: "server returns 404 not found",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte("Not Found"))
				}))
			},
			auth:        "Bearer test-token",
			disableTLS:  false,
			expectError: true,
			setupViper: func() {
				viper.Reset()
			},
			validateViper: func(t *testing.T) {
				// No values should be set on error
			},
		},
		{
			name: "server returns invalid JSON",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.Write([]byte("invalid json {"))
				}))
			},
			auth:        "Bearer test-token",
			disableTLS:  false,
			expectError: true,
			setupViper: func() {
				viper.Reset()
			},
			validateViper: func(t *testing.T) {
				// No values should be set on error
			},
		},
		{
			name: "empty property sources",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					response := map[string]any{
						"name":            "test-app",
						"profiles":        []string{"default"},
						"label":           "main",
						"version":         "1.0.0",
						"state":           "active",
						"propertySources": []map[string]any{},
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(response)
				}))
			},
			auth:        "Bearer test-token",
			disableTLS:  false,
			expectError: false,
			setupViper: func() {
				viper.Reset()
			},
			validateViper: func(t *testing.T) {
				// No error, but no values set either
			},
		},
		{
			name: "complex nested values",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					response := map[string]any{
						"name":     "test-app",
						"profiles": []string{"default"},
						"label":    "main",
						"version":  "1.0.0",
						"state":    "active",
						"propertySources": []map[string]any{
							{
								"name": "test-source",
								"source": map[string]any{
									"database.host":    "localhost",
									"database.port":    5432,
									"database.enabled": true,
									"features.list":    []string{"feature1", "feature2"},
									"metadata.tags":    map[string]any{"env": "test", "team": "dev"},
								},
							},
						},
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(response)
				}))
			},
			auth:        "Bearer test-token",
			disableTLS:  false,
			expectError: false,
			setupViper: func() {
				viper.Reset()
			},
			validateViper: func(t *testing.T) {
				if viper.GetString("database.host") != "localhost" {
					t.Errorf("Expected database.host to be 'localhost', got '%s'", viper.GetString("database.host"))
				}
				if viper.GetInt("database.port") != 5432 {
					t.Errorf("Expected database.port to be 5432, got %d", viper.GetInt("database.port"))
				}
				if !viper.GetBool("database.enabled") {
					t.Error("Expected database.enabled to be true")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			tt.setupViper()
			server := tt.setupServer()
			defer server.Close()

			// Execute
			scc2go.GetEnv(server.URL, tt.auth, tt.disableTLS)

			// Validate
			tt.validateViper(t)
		})
	}
}

func TestGetEnvEdgeCases(t *testing.T) {
	t.Run("empty URL", func(t *testing.T) {
		// Should log error and return early
		scc2go.GetEnv("", "Bearer token")
	})

	t.Run("with TLS disabled", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := map[string]any{
				"name":            "test-app",
				"profiles":        []string{"default"},
				"label":           "main",
				"version":         "1.0.0",
				"state":           "active",
				"propertySources": []map[string]any{},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		viper.Reset()
		scc2go.GetEnv(server.URL, "Bearer token", true)
	})

	t.Run("with TLS enabled (default)", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := map[string]any{
				"name":            "test-app",
				"profiles":        []string{"default"},
				"label":           "main",
				"version":         "1.0.0",
				"state":           "active",
				"propertySources": []map[string]any{},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		viper.Reset()
		scc2go.GetEnv(server.URL, "Bearer token", false)
	})
}

func TestGetEnvAuthorizationHeader(t *testing.T) {
	tests := []struct {
		name         string
		authHeader   string
		expectHeader string
	}{
		{
			name:         "Bearer token",
			authHeader:   "Bearer test-token-123",
			expectHeader: "Bearer test-token-123",
		},
		{
			name:         "Basic auth",
			authHeader:   "Basic dXNlcjpwYXNz",
			expectHeader: "Basic dXNlcjpwYXNz",
		},
		{
			name:         "Empty auth",
			authHeader:   "",
			expectHeader: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedAuth := r.Header.Get("Authorization")
				if receivedAuth != tt.expectHeader {
					t.Errorf("Expected Authorization header '%s', got '%s'", tt.expectHeader, receivedAuth)
				}

				response := map[string]any{
					"name":            "test-app",
					"profiles":        []string{"default"},
					"label":           "main",
					"version":         "1.0.0",
					"state":           "active",
					"propertySources": []map[string]any{},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			}))
			defer server.Close()

			viper.Reset()
			scc2go.GetEnv(server.URL, tt.authHeader)
		})
	}
}
