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
	"fmt"
	"strconv"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
)

// ExampleOption_Map demonstrates transforming the value of a Some.
func ExampleOption_Map() {
	fmt.Println(Some(21).Map(N.Mul(2)))
	fmt.Println(None[int]().Map(N.Mul(2)))
	// Output:
	// Some[int](42)
	// None[int]
}

// ExampleOption_MapTo demonstrates replacing the value of a Some.
func ExampleOption_MapTo() {
	fmt.Println(Some(1).MapTo("done"))
	// Output: Some[string](done)
}

// ExampleOption_Chain demonstrates sequencing operations that may produce None.
func ExampleOption_Chain() {
	half := FromPredicate(func(n int) bool { return n%2 == 0 })

	fmt.Println(Some(8).Chain(half).Map(N.Div(2)))
	fmt.Println(Some(7).Chain(half).Map(N.Div(2)))
	// Output:
	// Some[int](4)
	// None[int]
}

// ExampleOption_ChainFirst demonstrates validating a value without replacing it.
func ExampleOption_ChainFirst() {
	nonEmpty := FromPredicate(func(s string) bool { return s != "" })

	fmt.Println(Some("id-1").ChainFirst(nonEmpty))
	// Output: Some[string](id-1)
}

// ExampleOption_ChainTo demonstrates continuing with a fixed Option when a value is present.
func ExampleOption_ChainTo() {
	fmt.Println(Some(1).ChainTo(Some("next")))
	fmt.Println(None[int]().ChainTo(Some("next")))
	// Output:
	// Some[string](next)
	// None[string]
}

// ExampleOption_Fold demonstrates eliminating an Option into a plain value.
func ExampleOption_Fold() {
	describe := func(o Option[int]) string {
		return o.Fold(F.Constant("none"), strconv.Itoa)
	}

	fmt.Println(describe(Some(3)))
	fmt.Println(describe(None[int]()))
	// Output:
	// 3
	// none
}

// ExampleOption_GetOrElse demonstrates extracting a value with a fallback.
func ExampleOption_GetOrElse() {
	fmt.Println(Some(5).GetOrElse(F.Constant(0)))
	fmt.Println(None[int]().GetOrElse(F.Constant(0)))
	// Output:
	// 5
	// 0
}

// ExampleOption_Alt demonstrates falling back to another Option.
func ExampleOption_Alt() {
	fmt.Println(None[int]().Alt(F.Constant(Some(9))))
	// Output: Some[int](9)
}

// ExampleOption_Filter demonstrates discarding values that fail a predicate.
func ExampleOption_Filter() {
	fmt.Println(Some(4).Filter(N.MoreThan(0)))
	fmt.Println(Some(-4).Filter(N.MoreThan(0)))
	// Output:
	// Some[int](4)
	// None[int]
}
