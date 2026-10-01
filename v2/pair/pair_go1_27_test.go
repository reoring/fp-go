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

package pair_test

import (
	"fmt"
	"strconv"
	"testing"

	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/pair"
	SG "github.com/IBM/fp-go/v2/semigroup"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestPairMethods_Mapping(t *testing.T) {
	p := pair.MakePair("key", 21)

	t.Run("Map and MapTail change only the tail", func(t *testing.T) {
		assert.Equal(t, pair.MakePair("key", "21"), p.Map(strconv.Itoa))
		assert.Equal(t, pair.MakePair("key", "21"), p.MapTail(strconv.Itoa))
	})
	t.Run("MapHead changes only the head", func(t *testing.T) {
		assert.Equal(t, pair.MakePair(3, 21), p.MapHead(S.Size))
	})
	t.Run("BiMap changes both sides independently", func(t *testing.T) {
		assert.Equal(t, pair.MakePair(3, 42), p.BiMap(S.Size, N.Mul(2)))
	})
	t.Run("Swap exchanges sides and round-trips", func(t *testing.T) {
		assert.Equal(t, pair.MakePair(21, "key"), p.Swap())
		assert.Equal(t, p, p.Swap().Swap())
	})
}

func TestPairMethods_Chain(t *testing.T) {
	concat := SG.MakeSemigroup(S.Monoid.Concat)
	step := func(n int) pair.Pair[string, int] { return pair.MakePair("+step", n+1) }

	t.Run("combines the original head before the new head", func(t *testing.T) {
		assert.Equal(t, pair.MakePair("start+step", 2), pair.MakePair("start", 1).Chain(concat, step))
	})
	t.Run("accumulates heads across a chain", func(t *testing.T) {
		result := pair.MakePair("start", 1).Chain(concat, step).Chain(concat, step)
		assert.Equal(t, pair.MakePair("start+step+step", 3), result)
	})
}

// ExamplePair_Map demonstrates transforming the tail of a pair.
func ExamplePair_Map() {
	fmt.Println(pair.MakePair("port", 8080).Map(strconv.Itoa))
	// Output: Pair[string, string](port, 8080)
}

// ExamplePair_MapHead demonstrates transforming the head of a pair.
func ExamplePair_MapHead() {
	fmt.Println(pair.MakePair("port", 8080).MapHead(S.Size))
	// Output: Pair[int, int](4, 8080)
}

// ExamplePair_MapTail demonstrates transforming the tail of a pair.
func ExamplePair_MapTail() {
	fmt.Println(pair.MakePair("port", 8080).MapTail(N.Add(1)))
	// Output: Pair[string, int](port, 8081)
}

// ExamplePair_BiMap demonstrates transforming both sides of a pair.
func ExamplePair_BiMap() {
	fmt.Println(pair.MakePair("port", 8080).BiMap(S.Size, strconv.Itoa))
	// Output: Pair[int, string](4, 8080)
}

// ExamplePair_Chain demonstrates accumulating a log in the head while computing in the tail.
func ExamplePair_Chain() {
	concat := SG.MakeSemigroup(S.Monoid.Concat)
	double := func(n int) pair.Pair[string, int] { return pair.MakePair(" doubled", n*2) }

	fmt.Println(pair.MakePair("start", 5).Chain(concat, double).Chain(concat, double))
	// Output: Pair[string, int](start doubled doubled, 20)
}

// ExamplePair_Swap demonstrates exchanging head and tail.
func ExamplePair_Swap() {
	fmt.Println(pair.MakePair("port", 8080).Swap())
	// Output: Pair[int, string](8080, port)
}
