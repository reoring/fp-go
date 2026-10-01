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

package common

// The methods in this file are the receiver forms of the package-level Either
// functions. They require Go 1.27, which added type parameters on methods.
// Each method delegates to the data-first EitherMonad* function so that
// chaining methods allocates no closures.
//
// Signatures spell out function types such as func(A) Either[E, B] instead of
// the EitherKleisli alias: with Go 1.27.0, a generic alias that refers to the
// receiver type in a method signature deadlocks the compiler in importing
// packages.

// Map applies f to the value of a Right and returns a Left unchanged.
//
// Type Parameters:
//   - B: the new Right type
//
// Parameters:
//   - f: the function applied to the Right value
//
// Returns:
//   - Either[E, B]: Right(f(a)) for Right(a), the original Left otherwise
//
// See Also:
//   - EitherMonadMap: the equivalent data-first function
//
//go:inline
func (fa Either[E, A]) Map[B any](f func(A) B) Either[E, B] {
	return EitherMonadMap(fa, f)
}

// MapTo replaces the value of a Right with b and returns a Left unchanged.
//
// Type Parameters:
//   - B: the replacement Right type
//
// Parameters:
//   - b: the replacement value
//
// Returns:
//   - Either[E, B]: Right(b) for any Right, the original Left otherwise
//
// See Also:
//   - EitherMonadMapTo: the equivalent data-first function
//
//go:inline
func (fa Either[E, A]) MapTo[B any](b B) Either[E, B] {
	return EitherMonadMapTo(fa, b)
}

// MapLeft applies f to the value of a Left and returns a Right unchanged.
//
// Type Parameters:
//   - E2: the new Left type
//
// Parameters:
//   - f: the function applied to the Left value
//
// Returns:
//   - Either[E2, A]: Left(f(e)) for Left(e), the original Right otherwise
//
// See Also:
//   - EitherMonadMapLeft: the equivalent data-first function
//
//go:inline
func (fa Either[E, A]) MapLeft[E2 any](f func(E) E2) Either[E2, A] {
	return EitherMonadMapLeft(fa, f)
}

// BiMap applies f to a Left value or g to a Right value.
//
// Type Parameters:
//   - E2: the new Left type
//   - B: the new Right type
//
// Parameters:
//   - f: the function applied to the Left value
//   - g: the function applied to the Right value
//
// Returns:
//   - Either[E2, B]: Left(f(e)) for Left(e), Right(g(a)) for Right(a)
//
// See Also:
//   - EitherMonadBiMap: the equivalent data-first function
//
//go:inline
func (fa Either[E, A]) BiMap[E2, B any](f func(E) E2, g func(A) B) Either[E2, B] {
	return EitherMonadBiMap(fa, f, g)
}

// Chain applies f to the value of a Right and returns its result; a Left is returned unchanged.
//
// Type Parameters:
//   - B: the Right type of the Either returned by f
//
// Parameters:
//   - f: the Kleisli arrow applied to the Right value
//
// Returns:
//   - Either[E, B]: f(a) for Right(a), the original Left otherwise
//
// See Also:
//   - EitherMonadChain: the equivalent data-first function
//
//go:inline
func (fa Either[E, A]) Chain[B any](f func(A) Either[E, B]) Either[E, B] {
	return EitherMonadChain(fa, f)
}

// ChainFirst runs f on the value of a Right for its Either effect and keeps the
// original value. A Left from either side is returned.
//
// Type Parameters:
//   - B: the Right type of the Either returned by f, which is discarded
//
// Parameters:
//   - f: the Kleisli arrow whose success is required
//
// Returns:
//   - Either[E, A]: fa when fa is Right and f returns Right, otherwise the first Left
//
// See Also:
//   - EitherMonadChainFirst: the equivalent data-first function
//
//go:inline
func (fa Either[E, A]) ChainFirst[B any](f func(A) Either[E, B]) Either[E, A] {
	return EitherMonadChainFirst(fa, f)
}

// ChainLeft applies f to the value of a Left, typically to recover from an
// error; a Right is returned unchanged.
//
// Type Parameters:
//   - E2: the Left type of the Either returned by f
//
// Parameters:
//   - f: the function applied to the Left value
//
// Returns:
//   - Either[E2, A]: f(e) for Left(e), the original Right otherwise
//
// See Also:
//   - EitherMonadChainLeft: the equivalent data-first function
//
//go:inline
func (fa Either[E, A]) ChainLeft[E2 any](f func(E) Either[E2, A]) Either[E2, A] {
	return EitherMonadChainLeft(fa, f)
}

// Fold eliminates the Either: onLeft receives a Left value, onRight a Right
// value. Only the selected callback runs.
//
// Type Parameters:
//   - B: the result type
//
// Parameters:
//   - onLeft: produces the result from the Left value
//   - onRight: produces the result from the Right value
//
// Returns:
//   - B: the result of the selected callback
//
// See Also:
//   - EitherMonadFold: the equivalent data-first function
//
//go:inline
func (fa Either[E, A]) Fold[B any](onLeft func(E) B, onRight func(A) B) B {
	return EitherMonadFold(fa, onLeft, onRight)
}

// GetOrElse returns the Right value, or the result of onLeft for a Left.
// onLeft is not called for a Right.
//
// Parameters:
//   - onLeft: computes the fallback from the Left value
//
// Returns:
//   - A: the Right value or the fallback
//
// See Also:
//   - EitherMonadGetOrElse: the equivalent data-first function
//
//go:inline
func (fa Either[E, A]) GetOrElse(onLeft func(E) A) A {
	return EitherMonadGetOrElse(fa, onLeft)
}

// Alt returns fa if it is a Right, otherwise the Either produced by that.
// that is not called for a Right.
//
// Parameters:
//   - that: produces the alternative Either
//
// Returns:
//   - Either[E, A]: fa when Right, otherwise that()
//
// See Also:
//   - EitherMonadAlt: the equivalent data-first function
//
//go:inline
func (fa Either[E, A]) Alt(that func() Either[E, A]) Either[E, A] {
	return EitherMonadAlt(fa, that)
}

// Filter keeps a Right whose value satisfies p and turns a Right that fails p
// into Left(empty). A Left is returned unchanged.
//
// Parameters:
//   - p: the predicate the Right value must satisfy
//   - empty: the Left value used when p fails
//
// Returns:
//   - Either[E, A]: fa when it is Left or p holds, Left(empty) otherwise
//
// See Also:
//   - EitherMonadFilter: the equivalent data-first function
//
//go:inline
func (fa Either[E, A]) Filter(p Predicate[A], empty E) Either[E, A] {
	return EitherMonadFilter(fa, p, empty)
}

// Swap exchanges the Left and Right sides.
//
// Returns:
//   - Either[A, E]: Right(e) for Left(e), Left(a) for Right(a)
//
// See Also:
//   - EitherSwap: the equivalent function
//
//go:inline
func (fa Either[E, A]) Swap() Either[A, E] {
	return EitherSwap(fa)
}

// ToOption converts a Right to Some and a Left to None, discarding the Left value.
//
// Returns:
//   - Option[A]: Some(a) for Right(a), None for any Left
//
// See Also:
//   - EitherToOption: the equivalent function
//
//go:inline
func (fa Either[E, A]) ToOption() Option[A] {
	return EitherToOption(fa)
}
