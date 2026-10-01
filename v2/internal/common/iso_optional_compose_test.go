// Copyright (c) 2023 - 2026 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	"maps"
	"testing"

	"github.com/stretchr/testify/assert"
)

type composeConfig struct {
	settings map[string]int
}

type composePort struct {
	n int
}

// configSettings is a total Iso between composeConfig and its settings map.
var configSettings = MakeIso(
	func(c composeConfig) map[string]int { return c.settings },
	func(m map[string]int) composeConfig { return composeConfig{settings: m} },
)

// portEntry focuses on the "port" key; Set is a no-op when the key is absent.
var portEntry = MakeOptional(
	func(m map[string]int) Option[int] {
		if v, ok := m["port"]; ok {
			return OptionSome(v)
		}
		return OptionNone[int]()
	},
	func(m map[string]int, v int) map[string]int {
		if _, ok := m["port"]; !ok {
			return m
		}
		next := maps.Clone(m)
		next["port"] = v
		return next
	},
)

var intPort = MakeIso(
	func(n int) composePort { return composePort{n} },
	func(p composePort) int { return p.n },
)

func TestIsoComposeOptional(t *testing.T) {
	configPort := IsoComposeOptional[composeConfig](portEntry)(configSettings)
	present := composeConfig{settings: map[string]int{"port": 8080, "debug": 1}}
	absent := composeConfig{settings: map[string]int{"debug": 1}}

	t.Run("GetOption focuses through the Iso when the key exists", func(t *testing.T) {
		assert.Equal(t, OptionSome(8080), configPort.GetOption(present))
	})
	t.Run("GetOption is None when the key is missing", func(t *testing.T) {
		assert.Equal(t, OptionNone[int](), configPort.GetOption(absent))
	})
	t.Run("Set updates the focus and keeps the rest", func(t *testing.T) {
		updated := configPort.Set(9090)(present)
		assert.Equal(t, map[string]int{"port": 9090, "debug": 1}, updated.settings)
		assert.Equal(t, 8080, present.settings["port"], "source must not be mutated")
	})
	t.Run("Set leaves the source unchanged when the focus is missing", func(t *testing.T) {
		assert.Equal(t, absent, configPort.Set(9090)(absent))
	})
	t.Run("GetOption after Set returns the new value", func(t *testing.T) {
		assert.Equal(t, OptionSome(1), configPort.GetOption(configPort.Set(1)(present)))
	})
}

func TestOptionalComposeIso(t *testing.T) {
	typedPort := OptionalComposeIso[map[string]int](intPort)(portEntry)
	present := map[string]int{"port": 8080}
	absent := map[string]int{}

	t.Run("GetOption converts the focus through the Iso", func(t *testing.T) {
		assert.Equal(t, OptionSome(composePort{8080}), typedPort.GetOption(present))
	})
	t.Run("GetOption is None exactly when the Optional is", func(t *testing.T) {
		assert.Equal(t, OptionNone[composePort](), typedPort.GetOption(absent))
	})
	t.Run("Set converts back through the Iso", func(t *testing.T) {
		assert.Equal(t, map[string]int{"port": 1}, typedPort.Set(composePort{1})(present))
	})
	t.Run("Set leaves the source unchanged when the focus is missing", func(t *testing.T) {
		assert.Equal(t, absent, typedPort.Set(composePort{1})(absent))
	})
}
