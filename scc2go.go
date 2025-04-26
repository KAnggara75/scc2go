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
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"resty.dev/v3"
	"time"
)

var ErrorLoadConfig error

type springCloudConfig struct {
	Name            string           `json:"name"`
	Profiles        []string         `json:"profiles"`
	Label           interface{}      `json:"label"`
	Version         string           `json:"version"`
	State           interface{}      `json:"state"`
	PropertySources []propertySource `json:"propertySources"`
}

type propertySource struct {
	Name   string                 `json:"name"`
	Source map[string]interface{} `json:"source"`
}

func GetEnv(auth string) {
	sccUrl := os.Getenv("SCC_URL")

	if sccUrl != "" {
		logrus.Info("Using SCC URL: ", sccUrl)
		resBody, err := getSCC(sccUrl, auth)
		if err != nil {
			logrus.Errorf("Spring Cloud Config: %s\n", err)
		}

		scc := new(springCloudConfig)
		err = json.Unmarshal(resBody, scc)
		if err != nil {
			logrus.Errorf("Spring Cloud Config: %s\n", err)
		}

		for i := len(scc.PropertySources) - 1; i >= 0; i-- {
			for key, value := range scc.PropertySources[i].Source {
				setIfNotExists(key, value)
			}
		}

	} else {
		logrus.Warn("cloud config URL is not defined, please set cloud config url in env variable: SCC_URL")
	}
}

func setIfNotExists(k string, v interface{}) {
	if viper.Get(k) != nil {
		return
	}
	viper.Set(k, v)
	return
}

func getSCC(url, authHeader string) ([]byte, error) {
	client := resty.New().
		SetTimeout(5 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(time.Second)
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
		logrus.Error(err)
		return nil, fmt.Errorf("fail get config from %s with error: %s", url, err)
	}

	return res.Bytes(), nil
}
