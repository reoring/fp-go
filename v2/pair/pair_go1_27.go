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

package pair

// The methods in this file are the receiver forms of the package-level Pair
// functions. They require Go 1.27, which added type parameters on methods.
// Each method delegates to the data-first function of the same name.
//
// Signatures spell out func(R) Pair[L, R1] instead of the Kleisli alias: with
// Go 1.27.0, a generic alias that refers to the receiver type in a method
// signature deadlocks the compiler in importing packages.

// MapHead applies f to the head (left) value and keeps the tail.
//
// Type Parameters:
//   - L1: the new head type
//
// Parameters:
//   - f: the function applied to the head
//
// Returns:
//   - Pair[L1, R]: a pair of f(head) and the original tail
//
// See Also:
//   - MonadMapHead: the equivalent data-first function
//
//go:inline
func (p Pair[L, R]) MapHead[L1 any](f func(L) L1) Pair[L1, R] {
	return MonadMapHead(p, f)
}

// MapTail applies f to the tail (right) value and keeps the head.
//
// Type Parameters:
//   - R1: the new tail type
//
// Parameters:
//   - f: the function applied to the tail
//
// Returns:
//   - Pair[L, R1]: a pair of the original head and f(tail)
//
// See Also:
//   - MonadMapTail: the equivalent data-first function
//
//go:inline
func (p Pair[L, R]) MapTail[R1 any](f func(R) R1) Pair[L, R1] {
	return MonadMapTail(p, f)
}

// Map applies f to the tail (right) value and keeps the head. It is the same
// operation as MapTail, matching the package-level Map.
//
// Type Parameters:
//   - R1: the new tail type
//
// Parameters:
//   - f: the function applied to the tail
//
// Returns:
//   - Pair[L, R1]: a pair of the original head and f(tail)
//
// See Also:
//   - MonadMap: the equivalent data-first function
//
//go:inline
func (p Pair[L, R]) Map[R1 any](f func(R) R1) Pair[L, R1] {
	return MonadMap(p, f)
}

// BiMap applies f to the head and g to the tail.
//
// Type Parameters:
//   - L1: the new head type
//   - R1: the new tail type
//
// Parameters:
//   - f: the function applied to the head
//   - g: the function applied to the tail
//
// Returns:
//   - Pair[L1, R1]: a pair of f(head) and g(tail)
//
// See Also:
//   - MonadBiMap: the equivalent data-first function
//
//go:inline
func (p Pair[L, R]) BiMap[L1, R1 any](f func(L) L1, g func(R) R1) Pair[L1, R1] {
	return MonadBiMap(p, f, g)
}

// Chain applies f to the tail and combines the heads of both pairs with sg.
//
// Type Parameters:
//   - R1: the tail type of the pair returned by f
//
// Parameters:
//   - sg: combines the original head with the head returned by f
//   - f: the function applied to the tail
//
// Returns:
//   - Pair[L, R1]: a pair of sg.Concat(head, head of f(tail)) and the tail of f(tail)
//
// See Also:
//   - MonadChain: the equivalent data-first function
//
//go:inline
func (p Pair[L, R]) Chain[R1 any](sg Semigroup[L], f func(R) Pair[L, R1]) Pair[L, R1] {
	return MonadChain(sg, p, f)
}

// Swap exchanges head and tail.
//
// Returns:
//   - Pair[R, L]: a pair of the original tail and head
//
// See Also:
//   - Swap: the equivalent package-level function
//
//go:inline
func (p Pair[L, R]) Swap() Pair[R, L] {
	return Swap(p)
}
