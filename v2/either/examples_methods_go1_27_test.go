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

package either_test

import (
	"errors"
	"fmt"
	"strconv"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	S "github.com/IBM/fp-go/v2/string"
)

// ExampleEither_Map demonstrates transforming the value of a Right.
func ExampleEither_Map() {
	fmt.Println(E.Right[string](21).Map(N.Mul(2)))
	fmt.Println(E.Left[int]("failed").Map(N.Mul(2)))
	// Output:
	// Right[int](42)
	// Left[string](failed)
}

// ExampleEither_MapTo demonstrates replacing the value of a Right.
func ExampleEither_MapTo() {
	fmt.Println(E.Right[string](1).MapTo("done"))
	// Output: Right[string](done)
}

// ExampleEither_MapLeft demonstrates transforming the error side.
func ExampleEither_MapLeft() {
	fmt.Println(E.Left[int](errors.New("boom")).MapLeft(error.Error))
	// Output: Left[string](boom)
}

// ExampleEither_BiMap demonstrates transforming whichever side is present.
func ExampleEither_BiMap() {
	fmt.Println(E.Right[string](7).BiMap(S.Size, strconv.Itoa))
	fmt.Println(E.Left[int]("short").BiMap(S.Size, strconv.Itoa))
	// Output:
	// Right[string](7)
	// Left[int](5)
}

// ExampleEither_Chain demonstrates sequencing computations that may fail.
func ExampleEither_Chain() {
	parse := E.Eitherize1(strconv.Atoi)

	fmt.Println(E.Right[error]("21").Chain(parse).Map(N.Mul(2)))
	// Output: Right[int](42)
}

// ExampleEither_ChainFirst demonstrates validating a value without replacing it.
func ExampleEither_ChainFirst() {
	nonEmpty := E.FromPredicate(S.IsNonEmpty, F.Constant1[string]("empty"))

	fmt.Println(E.Right[string]("id-1").ChainFirst(nonEmpty))
	fmt.Println(E.Right[string]("").ChainFirst(nonEmpty))
	// Output:
	// Right[string](id-1)
	// Left[string](empty)
}

// ExampleEither_ChainLeft demonstrates recovering from a Left.
func ExampleEither_ChainLeft() {
	useDefault := F.Constant1[string](E.Right[string](0))

	fmt.Println(E.Left[int]("missing").ChainLeft(useDefault))
	// Output: Right[int](0)
}

// ExampleEither_Fold demonstrates eliminating an Either into a plain value.
func ExampleEither_Fold() {
	describe := func(e E.Either[string, int]) string {
		return e.Fold(S.Format[string]("error: %s"), strconv.Itoa)
	}

	fmt.Println(describe(E.Right[string](3)))
	fmt.Println(describe(E.Left[int]("bad input")))
	// Output:
	// 3
	// error: bad input
}

// ExampleEither_GetOrElse demonstrates extracting a value with a fallback computed from the Left.
func ExampleEither_GetOrElse() {
	fmt.Println(E.Right[string](5).GetOrElse(S.Size))
	fmt.Println(E.Left[int]("oops").GetOrElse(S.Size))
	// Output:
	// 5
	// 4
}

// ExampleEither_Alt demonstrates falling back to another Either.
func ExampleEither_Alt() {
	fmt.Println(E.Left[int]("primary down").Alt(F.Constant(E.Right[string](9))))
	// Output: Right[int](9)
}

// ExampleEither_Filter demonstrates turning a Right that fails a predicate into a Left.
func ExampleEither_Filter() {
	fmt.Println(E.Right[string](4).Filter(N.MoreThan(0), "not positive"))
	fmt.Println(E.Right[string](-4).Filter(N.MoreThan(0), "not positive"))
	// Output:
	// Right[int](4)
	// Left[string](not positive)
}

// ExampleEither_Swap demonstrates exchanging the Left and Right sides.
func ExampleEither_Swap() {
	fmt.Println(E.Left[int]("x").Swap())
	// Output: Right[string](x)
}

// ExampleEither_ToOption demonstrates discarding the Left value.
func ExampleEither_ToOption() {
	fmt.Println(E.Right[string](1).ToOption())
	fmt.Println(E.Left[int]("x").ToOption())
	// Output:
	// Some[int](1)
	// None[int]
}
