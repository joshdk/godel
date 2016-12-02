// Copyright 2016 Palantir Technologies, Inc.
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

package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/palantir/pkg/matcher"
	"github.com/palantir/pkg/specdir"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/palantir/godel/layout"
)

const (
	ExcludeYML = "exclude.yml"
)

func GetExcludeCfgFromYML(cfgPath string) (matcher.NamesPathsCfg, error) {
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return matcher.NamesPathsCfg{}, nil
	}

	fileBytes, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return matcher.NamesPathsCfg{}, errors.Wrapf(err, "Failed to read file %s", cfgPath)
	}

	excludeCfg := matcher.NamesPathsCfg{}
	if err := yaml.Unmarshal(fileBytes, &excludeCfg); err != nil {
		return matcher.NamesPathsCfg{}, errors.Wrapf(err, "Failed to unmarshal file %s", cfgPath)
	}
	return excludeCfg, nil
}

type excludeCfgStruct struct {
	Exclude matcher.NamesPathsCfg `json:"exclude"`
}

func ReadExcludeJSONFromYML(cfgPath string) ([]byte, error) {
	excludeCfg, err := GetExcludeCfgFromYML(cfgPath)
	if err != nil {
		return nil, err
	}
	return json.Marshal(excludeCfgStruct{
		Exclude: excludeCfg,
	})
}

func GetCfgDirPath(cfgDirPath, wrapperScriptPath string) (string, error) {
	if cfgDirPath != "" {
		return cfgDirPath, nil
	}
	if wrapperScriptPath == "" {
		return "", nil
	}
	wrapper, err := specdir.New(path.Dir(wrapperScriptPath), layout.WrapperSpec(), nil, specdir.Validate)
	if err != nil {
		return "", err
	}
	return wrapper.Path(layout.WrapperConfigDir), nil
}
