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

package result_test

import (
	"errors"
	"fmt"
	"strconv"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	R "github.com/IBM/fp-go/v2/result"
)

// ExampleResult_methods demonstrates that Result, as an alias of Either[error, A],
// supports the Either methods on Go 1.27 and later.
func ExampleResult_methods() {
	parse := R.Eitherize1(strconv.Atoi)
	positive := R.FromPredicate(N.MoreThan(0), F.Constant1[int](errors.New("not positive")))

	process := func(s string) int {
		return parse(s).
			Chain(positive).
			Map(N.Mul(10)).
			GetOrElse(F.Constant1[error](-1))
	}

	fmt.Println(process("4"))
	fmt.Println(process("-4"))
	fmt.Println(process("x"))
	// Output:
	// 40
	// -1
	// -1
}
