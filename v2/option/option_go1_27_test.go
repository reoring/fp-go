//go:build go1.27

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

package option

import (
	"strconv"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/stretchr/testify/assert"
)

func TestOptionMethods_MapChain(t *testing.T) {
	parse := func(s string) Option[int] {
		n, err := strconv.Atoi(s)
		if err != nil {
			return None[int]()
		}
		return Some(n)
	}

	t.Run("Some flows through Map and Chain", func(t *testing.T) {
		assert.Equal(t, Some(84), Some("42").Chain(parse).Map(N.Mul(2)))
	})
	t.Run("Chain returning None short-circuits Map", func(t *testing.T) {
		called := false
		result := Some("x").Chain(parse).Map(func(n int) int { called = true; return n })
		assert.Equal(t, None[int](), result)
		assert.False(t, called)
	})
	t.Run("None skips Chain", func(t *testing.T) {
		called := false
		result := None[string]().Chain(func(s string) Option[int] { called = true; return parse(s) })
		assert.Equal(t, None[int](), result)
		assert.False(t, called)
	})
	t.Run("MapTo replaces Some and keeps None", func(t *testing.T) {
		assert.Equal(t, Some("v"), Some(1).MapTo("v"))
		assert.Equal(t, None[string](), None[int]().MapTo("v"))
	})
	t.Run("ChainTo returns the argument only for Some", func(t *testing.T) {
		assert.Equal(t, Some("v"), Some(1).ChainTo(Some("v")))
		assert.Equal(t, None[string](), None[int]().ChainTo(Some("v")))
		assert.Equal(t, None[string](), Some(1).ChainTo(None[string]()))
	})
}

func TestOptionMethods_ChainFirst(t *testing.T) {
	t.Run("keeps the original value when f succeeds", func(t *testing.T) {
		assert.Equal(t, Some(5), Some(5).ChainFirst(func(n int) Option[string] { return Some(strconv.Itoa(n)) }))
	})
	t.Run("returns None when f fails", func(t *testing.T) {
		assert.Equal(t, None[int](), Some(5).ChainFirst(func(int) Option[string] { return None[string]() }))
	})
	t.Run("does not call f for None", func(t *testing.T) {
		called := false
		result := None[int]().ChainFirst(func(int) Option[string] { called = true; return Some("") })
		assert.Equal(t, None[int](), result)
		assert.False(t, called)
	})
}

func TestOptionMethods_Elimination(t *testing.T) {
	t.Run("Fold runs only the Some branch for Some", func(t *testing.T) {
		noneCalled := false
		result := Some(3).Fold(func() int { noneCalled = true; return 0 }, N.Mul(2))
		assert.Equal(t, 6, result)
		assert.False(t, noneCalled)
	})
	t.Run("Fold runs only the None branch for None", func(t *testing.T) {
		someCalled := false
		result := None[int]().Fold(F.Constant(-1), func(n int) int { someCalled = true; return n })
		assert.Equal(t, -1, result)
		assert.False(t, someCalled)
	})
	t.Run("GetOrElse does not evaluate the fallback for Some", func(t *testing.T) {
		called := false
		assert.Equal(t, 7, Some(7).GetOrElse(func() int { called = true; return 0 }))
		assert.False(t, called)
		assert.Equal(t, 0, None[int]().GetOrElse(F.Constant(0)))
	})
	t.Run("Alt does not evaluate the alternative for Some", func(t *testing.T) {
		called := false
		assert.Equal(t, Some(1), Some(1).Alt(func() Option[int] { called = true; return Some(2) }))
		assert.False(t, called)
		assert.Equal(t, Some(2), None[int]().Alt(F.Constant(Some(2))))
		assert.Equal(t, None[int](), None[int]().Alt(None[int]))
	})
}

func TestOptionMethods_Filter(t *testing.T) {
	positive := N.MoreThan(0)
	assert.Equal(t, Some(1), Some(1).Filter(positive))
	assert.Equal(t, None[int](), Some(-1).Filter(positive))
	assert.Equal(t, None[int](), None[int]().Filter(positive))
}

// The method benchmarks mirror the free-function ones below to show that
// method chaining adds no allocations over direct MonadX calls.

func BenchmarkOptionMethods_MapChainFold(b *testing.B) {
	opt := Some(1)
	b.ReportAllocs()
	for b.Loop() {
		_ = opt.
			Map(func(x int) int { return x + 1 }).
			Chain(func(x int) Option[int] { return Some(x * 2) }).
			Fold(func() int { return 0 }, func(x int) int { return x - 5 })
	}
}

func BenchmarkOptionMonad_MapChainFold(b *testing.B) {
	opt := Some(1)
	b.ReportAllocs()
	for b.Loop() {
		_ = MonadFold(
			MonadChain(
				MonadMap(opt, func(x int) int { return x + 1 }),
				func(x int) Option[int] { return Some(x * 2) },
			),
			func() int { return 0 },
			func(x int) int { return x - 5 },
		)
	}
}

func BenchmarkOptionPipe_MapChainFold(b *testing.B) {
	opt := Some(1)
	b.ReportAllocs()
	for b.Loop() {
		_ = F.Pipe3(
			opt,
			Map(func(x int) int { return x + 1 }),
			Chain(func(x int) Option[int] { return Some(x * 2) }),
			Fold(func() int { return 0 }, func(x int) int { return x - 5 }),
		)
	}
}

func BenchmarkOptionMethods_Chain5(b *testing.B) {
	opt := Some(1)
	b.ReportAllocs()
	for b.Loop() {
		_ = opt.
			Chain(func(x int) Option[int] { return Some(x + 1) }).
			Chain(func(x int) Option[int] { return Some(x * 2) }).
			Chain(func(x int) Option[int] { return Some(x - 5) }).
			Chain(func(x int) Option[int] { return Some(x * 10) }).
			Chain(func(x int) Option[int] { return Some(x + 100) })
	}
}

func BenchmarkOptionMonad_Chain5(b *testing.B) {
	opt := Some(1)
	b.ReportAllocs()
	for b.Loop() {
		_ = MonadChain(
			MonadChain(
				MonadChain(
					MonadChain(
						MonadChain(opt, func(x int) Option[int] { return Some(x + 1) }),
						func(x int) Option[int] { return Some(x * 2) },
					),
					func(x int) Option[int] { return Some(x - 5) },
				),
				func(x int) Option[int] { return Some(x * 10) },
			),
			func(x int) Option[int] { return Some(x + 100) },
		)
	}
}
