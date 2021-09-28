/*******************************************************************************
 * Copyright 2021 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package config

import (
	"encoding/json"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	logging "github.com/project-alvarium/provider-logging/pkg/config"
)

type MongoConfig struct {
	Host       string `json:"host,omitempty"`
	Port       int    `json:"port,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	Collection string `json:"collection,omitempty"`
	DbName     string `json:"dbName,omitempty"`
}

type ApplicationConfig struct {
	Mongo   MongoConfig         `json:"mongo,omitempty"`
	Sdk     config.SdkInfo      `json:"sdk,omitempty"`
	Logging logging.LoggingInfo `json:"logging,omitempty"`
}

func (a ApplicationConfig) AsString() string {
	b, _ := json.Marshal(a)
	return string(b)
}
