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

// The methods in this file are the receiver forms of the package-level Option
// functions. They require Go 1.27, which added type parameters on methods.
// Each method delegates to the data-first OptionMonad* function so that
// chaining methods allocates no closures.

// Map applies f to the value of a Some and returns None unchanged.
//
// Type Parameters:
//   - B: the result value type
//
// Parameters:
//   - f: the function applied to the contained value
//
// Returns:
//   - Option[B]: Some(f(a)) for Some(a), None otherwise
//
// See Also:
//   - OptionMonadMap: the equivalent data-first function
//
//go:inline
func (fa Option[A]) Map[B any](f func(A) B) Option[B] {
	return OptionMonadMap(fa, f)
}

// MapTo replaces the value of a Some with b and returns None unchanged.
//
// Type Parameters:
//   - B: the replacement value type
//
// Parameters:
//   - b: the replacement value
//
// Returns:
//   - Option[B]: Some(b) for any Some, None otherwise
//
// See Also:
//   - OptionMonadMapTo: the equivalent data-first function
//
//go:inline
func (fa Option[A]) MapTo[B any](b B) Option[B] {
	return OptionMonadMapTo(fa, b)
}

// Chain applies f to the value of a Some and returns its result; None is returned unchanged.
//
// Type Parameters:
//   - B: the value type of the Option returned by f
//
// Parameters:
//   - f: the Kleisli arrow applied to the contained value
//
// Returns:
//   - Option[B]: f(a) for Some(a), None otherwise
//
// See Also:
//   - OptionMonadChain: the equivalent data-first function
//
//go:inline
func (fa Option[A]) Chain[B any](f func(A) Option[B]) Option[B] {
	return OptionMonadChain(fa, f)
}

// ChainFirst runs f on the value of a Some for its Option effect and keeps the
// original value. The result is None if fa is None or f returns None.
//
// Type Parameters:
//   - B: the value type of the Option returned by f, which is discarded
//
// Parameters:
//   - f: the Kleisli arrow whose success is required
//
// Returns:
//   - Option[A]: fa when fa is Some and f returns Some, None otherwise
//
// See Also:
//   - OptionMonadChainFirst: the equivalent data-first function
//
//go:inline
func (fa Option[A]) ChainFirst[B any](f func(A) Option[B]) Option[A] {
	return OptionMonadChainFirst(fa, f)
}

// ChainTo returns mb if fa is Some and None otherwise.
//
// Type Parameters:
//   - B: the value type of mb
//
// Parameters:
//   - mb: the Option returned when fa is Some
//
// Returns:
//   - Option[B]: mb for any Some, None otherwise
//
// See Also:
//   - OptionMonadChainTo: the equivalent data-first function
//
//go:inline
func (fa Option[A]) ChainTo[B any](mb Option[B]) Option[B] {
	return OptionMonadChainTo(fa, mb)
}

// Fold eliminates the Option: onSome receives the value of a Some, onNone is
// called for None. Only the selected callback runs.
//
// Type Parameters:
//   - B: the result type
//
// Parameters:
//   - onNone: produces the result for None
//   - onSome: produces the result from the contained value
//
// Returns:
//   - B: the result of the selected callback
//
// See Also:
//   - OptionMonadFold: the equivalent data-first function
//
//go:inline
func (fa Option[A]) Fold[B any](onNone func() B, onSome func(A) B) B {
	return OptionMonadFold(fa, onNone, onSome)
}

// GetOrElse returns the value of a Some, or the result of onNone for None.
// onNone is not called for Some.
//
// Parameters:
//   - onNone: produces the fallback value
//
// Returns:
//   - A: the contained value or the fallback
//
// See Also:
//   - OptionMonadGetOrElse: the equivalent data-first function
//
//go:inline
func (fa Option[A]) GetOrElse(onNone func() A) A {
	return OptionMonadGetOrElse(fa, onNone)
}

// Alt returns fa if it is Some, otherwise the Option produced by that.
// that is not called for Some.
//
// Parameters:
//   - that: produces the alternative Option
//
// Returns:
//   - Option[A]: fa when Some, otherwise that()
//
// See Also:
//   - OptionMonadAlt: the equivalent data-first function
//
//go:inline
func (fa Option[A]) Alt(that func() Option[A]) Option[A] {
	return OptionMonadAlt(fa, that)
}

// Filter keeps a Some whose value satisfies pred and returns None otherwise.
//
// Parameters:
//   - pred: the predicate the contained value must satisfy
//
// Returns:
//   - Option[A]: fa when it is Some and pred holds, None otherwise
//
// See Also:
//   - OptionMonadFilter: the equivalent data-first function
//
//go:inline
func (fa Option[A]) Filter(pred Predicate[A]) Option[A] {
	return OptionMonadFilter(fa, pred)
}
