// Public Domain (-) 2011 The Ampify Authors.
// See the Ampify UNLICENSE file for details.

package nodule

import (
	"amp/runtime"
	"bytes"
	"exp/template"
	"os"
	"strings"
)

type Config struct {
	Type    string
	Build   []string
	Run     []string
	Test    []string
	Depends []string
	Ignore  []string
}

type ConfigEnv struct {
	Profile  string
	Platform string
	Darwin   bool
	Linux    bool
	FreeBSD  bool
}

var defaultEnv *ConfigEnv

func GetConfigEnv() *ConfigEnv {
	if defaultEnv != nil {
		return defaultEnv
	}
	env := &ConfigEnv{
		Profile:  runtime.Profile,
		Platform: runtime.Platform,
	}
	switch runtime.Platform {
	case "linux":
		env.Linux = true
	case "freebsd":
		env.FreeBSD = true
	case "darwin":
		env.Darwin = true
	}
	defaultEnv = env
	return env
}

func EvalStrings(name string, list []string, data interface{}) ([]string, os.Error) {
	result := make([]string, len(list))
	for idx, value := range list {
		if strings.IndexRune(value, '{') == -1 {
			result[idx] = value
		} else {
			tpl := template.New(name)
			buf := &bytes.Buffer{}
			err := tpl.Parse(value)
			if err != nil {
				return nil, err
			}
			err = tpl.Execute(buf, data)
			if err != nil {
				return nil, err
			}
			result[idx] = buf.String()
		}
	}
	return result, nil
}
