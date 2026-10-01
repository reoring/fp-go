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

package either

import (
	"errors"
	"strconv"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

var errMethod = errors.New("method error")

func TestEitherMethods_RightBias(t *testing.T) {
	parse := func(s string) Either[error, int] {
		n, err := strconv.Atoi(s)
		if err != nil {
			return Left[int](errMethod)
		}
		return Right[error](n)
	}

	t.Run("Right flows through Chain and Map", func(t *testing.T) {
		assert.Equal(t, Right[error](84), Right[error]("42").Chain(parse).Map(N.Mul(2)))
	})
	t.Run("Left from Chain short-circuits Map", func(t *testing.T) {
		called := false
		result := Right[error]("x").Chain(parse).Map(func(n int) int { called = true; return n })
		assert.Equal(t, Left[int](errMethod), result)
		assert.False(t, called)
	})
	t.Run("Left skips Chain and keeps its value", func(t *testing.T) {
		called := false
		result := Left[string](errMethod).Chain(func(s string) Either[error, int] { called = true; return parse(s) })
		assert.Equal(t, Left[int](errMethod), result)
		assert.False(t, called)
	})
	t.Run("MapTo replaces Right and keeps Left", func(t *testing.T) {
		assert.Equal(t, Right[error]("v"), Right[error](1).MapTo("v"))
		assert.Equal(t, Left[string](errMethod), Left[int](errMethod).MapTo("v"))
	})
}

func TestEitherMethods_LeftSide(t *testing.T) {
	t.Run("MapLeft changes the Left type and leaves Right untouched", func(t *testing.T) {
		assert.Equal(t, Left[int]("method error"), Left[int](errMethod).MapLeft(error.Error))
		assert.Equal(t, Right[string](1), Right[error](1).MapLeft(error.Error))
	})
	t.Run("BiMap applies the matching side only", func(t *testing.T) {
		assert.Equal(t, Left[string](5), Left[int]("short").BiMap(S.Size, strconv.Itoa))
		assert.Equal(t, Right[int]("7"), Right[string](7).BiMap(S.Size, strconv.Itoa))
	})
	t.Run("ChainLeft recovers a Left and may change the Left type", func(t *testing.T) {
		recoverAll := func(error) Either[string, int] { return Right[string](0) }
		assert.Equal(t, Right[string](0), Left[int](errMethod).ChainLeft(recoverAll))
	})
	t.Run("ChainLeft does not call f for Right", func(t *testing.T) {
		called := false
		result := Right[error](3).ChainLeft(func(error) Either[string, int] { called = true; return Right[string](0) })
		assert.Equal(t, Right[string](3), result)
		assert.False(t, called)
	})
	t.Run("Swap exchanges sides and round-trips", func(t *testing.T) {
		assert.Equal(t, Right[int](errMethod), Left[int](errMethod).Swap())
		assert.Equal(t, Left[error](1), Right[error](1).Swap())
		assert.Equal(t, Right[error](1), Right[error](1).Swap().Swap())
	})
}

func TestEitherMethods_ChainFirst(t *testing.T) {
	validate := func(n int) Either[error, string] {
		if n < 0 {
			return Left[string](errMethod)
		}
		return Right[error]("ok")
	}
	assert.Equal(t, Right[error](5), Right[error](5).ChainFirst(validate))
	assert.Equal(t, Left[int](errMethod), Right[error](-5).ChainFirst(validate))

	other := errors.New("other")
	assert.Equal(t, Left[int](other), Left[int](other).ChainFirst(validate))
}

func TestEitherMethods_Elimination(t *testing.T) {
	t.Run("Fold runs only the matching branch", func(t *testing.T) {
		leftCalled := false
		result := Right[error](3).Fold(func(error) int { leftCalled = true; return 0 }, N.Mul(2))
		assert.Equal(t, 6, result)
		assert.False(t, leftCalled)

		rightCalled := false
		result = Left[int](errMethod).Fold(F.Constant1[error](-1), func(n int) int { rightCalled = true; return n })
		assert.Equal(t, -1, result)
		assert.False(t, rightCalled)
	})
	t.Run("GetOrElse receives the Left value and is lazy for Right", func(t *testing.T) {
		called := false
		assert.Equal(t, 7, Right[string](7).GetOrElse(func(string) int { called = true; return 0 }))
		assert.False(t, called)
		assert.Equal(t, 5, Left[int]("short").GetOrElse(S.Size))
	})
	t.Run("Alt is lazy for Right and falls back for Left", func(t *testing.T) {
		called := false
		assert.Equal(t, Right[error](1), Right[error](1).Alt(func() Either[error, int] { called = true; return Right[error](2) }))
		assert.False(t, called)
		assert.Equal(t, Right[error](2), Left[int](errMethod).Alt(F.Constant(Right[error](2))))
	})
	t.Run("ToOption drops the Left value", func(t *testing.T) {
		assert.Equal(t, O.Some(1), Right[error](1).ToOption())
		assert.Equal(t, O.None[int](), Left[int](errMethod).ToOption())
	})
}

func TestEitherMethods_Filter(t *testing.T) {
	positive := N.MoreThan(0)
	assert.Equal(t, Right[string](1), Right[string](1).Filter(positive, "not positive"))
	assert.Equal(t, Left[int]("not positive"), Right[string](-1).Filter(positive, "not positive"))
	assert.Equal(t, Left[int]("original"), Left[int]("original").Filter(positive, "not positive"))
}

func BenchmarkEitherMethods_MapChainFold(b *testing.B) {
	e := Right[error](1)
	b.ReportAllocs()
	for b.Loop() {
		_ = e.
			Map(func(x int) int { return x + 1 }).
			Chain(func(x int) Either[error, int] { return Right[error](x * 2) }).
			Fold(func(error) int { return 0 }, func(x int) int { return x - 5 })
	}
}

func BenchmarkEitherMonad_MapChainFold(b *testing.B) {
	e := Right[error](1)
	b.ReportAllocs()
	for b.Loop() {
		_ = MonadFold(
			MonadChain(
				MonadMap(e, func(x int) int { return x + 1 }),
				func(x int) Either[error, int] { return Right[error](x * 2) },
			),
			func(error) int { return 0 },
			func(x int) int { return x - 5 },
		)
	}
}

func BenchmarkEitherPipe_MapChainFold(b *testing.B) {
	e := Right[error](1)
	b.ReportAllocs()
	for b.Loop() {
		_ = F.Pipe3(
			e,
			Map[error](func(x int) int { return x + 1 }),
			Chain(func(x int) Either[error, int] { return Right[error](x * 2) }),
			Fold(func(error) int { return 0 }, func(x int) int { return x - 5 }),
		)
	}
}
