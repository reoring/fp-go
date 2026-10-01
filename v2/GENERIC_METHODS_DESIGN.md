# Generic Methods Design (Go 1.27)

Status: implemented on branch `go1.27-generic-methods`, except
`ComposeTraversal` (deferred). Scope: `v2` only; the root (v1) module is frozen.

## Table of Contents

- [Goal](#goal)
- [Go 1.27 Language Changes](#go-127-language-changes)
- [Current State](#current-state)
- [Constraints](#constraints)
- [Design](#design)
- [Implementation Plan](#implementation-plan)
- [Open Questions](#open-questions)
- [Non-Goals](#non-goals)

## Goal

Add a method-chaining entry point to the defined data types of fp-go v2 using
Go 1.27 generic methods, without changing the existing data-last function API.

```go
// Existing (stays the canonical form for building pipelines)
result := F.Pipe2(opt, O.Map(parse), O.Chain(validate))

// New (equivalent, for transforming a value in place)
result := opt.Map(parse).Chain(validate)
```

A method is always a thin delegation to the existing free function with the
same name and semantics. Methods never introduce behavior of their own.

## Go 1.27 Language Changes

From the [Go 1.27 release notes](https://go.dev/doc/go1.27), "Changes to the language":

1. **Generic methods** (proposal [#77273](https://github.com/golang/go/issues/77273)).
   A method may declare its own type parameters in addition to those of the
   receiver: `func (o Option[A]) Map[B any](f func(A) B) Option[B]`. Method
   type arguments are inferred from call arguments exactly like for functions.
   Restrictions that matter here:
   - Interface methods cannot have type parameters, and a generic method never
     satisfies an interface method.
   - Receiver type parameters must be plain identifiers. A receiver such as
     `Option[Option[A]]` or `Option[func(A) B]` is not expressible.
   - Method expressions require an instantiated receiver:
     `Option[int].Map[string]` is valid, `Option.Map` is not.
   - Generic methods are invisible to `reflect`.
2. **Generalized function type inference** ([#77245](https://github.com/golang/go/issues/77245)).
   A generic function used in a struct/array/slice/map literal, channel send,
   or conversion has its type arguments inferred from the target type.
3. Struct literal keys may be promoted-field selectors. Not relevant to fp-go.

Only (1) changes the design space for fp-go. (2) removes a handful of explicit
instantiations in test map literals (for example
`effect/effect_contract_test.go:40`) and nothing in library code.

## Current State

### Already adopted

Twenty generic methods exist, all behind `//go:build go1.27`:

| File | Methods |
|------|---------|
| `internal/common/lens_go1_27.go` | `Lens.Compose`, `ComposeIso`, `ComposePrism`, `ComposeOptional` |
| `internal/common/prism_go1_27.go` | `Prism.Compose`, `ComposeLens`, `ComposeIso`, `ComposeOptional` |
| `internal/common/optional_go1_27.go` | `Optional.Compose`, `ComposeLens`, `ComposePrism` |
| `internal/common/iso_go1_27.go` | `Iso.Compose`, `ComposeLens` |
| `optics/codec/go1_27.go` | `Type.Pipe`, `PipeIso`, `PipeRefinement`, `PipePrism`, `ApSL`, `ApSO`, `Bind` |

Established conventions from these files:

- Each method is a one-liner delegating to the free function
  (`optics/codec/go1_27.go:60`: `return Pipe[O, I](ab)(t)`).
- Doc comments show "method form" and "free-function form" side by side.
- Tests assert that the method and the free function produce the same result
  (`internal/common/lens_go1_27_test.go`).

These files carry the tag on purpose: `v2/go.mod` declares `go 1.24.0` and CI
(`.github/workflows/build.yml`) tests v2 on Go 1.24.x through 1.27.x. A
`//go:build go1.27` constraint raises the language version of that single file
to 1.27, so the methods exist on Go 1.27+ toolchains while the rest of the
module still builds on 1.24. Verified: with `go 1.24.0` in `go.mod`, the Go
1.27.0 toolchain builds `internal/common` and `optics/...` and the
`*_go1_27_test.go` tests pass.

### Documentation drift

`optics/README.md:176` and `optics/doc.go:158` claim every optic has
`ComposeLens`, `ComposePrism`, `ComposeOptional`, `ComposeIso` and
`ComposeTraversal`. In fact `Lens` has `Compose` (not `ComposeLens`),
`Iso` lacks `ComposePrism`/`ComposeOptional`, `Optional` lacks `ComposeIso`,
and `ComposeTraversal` does not exist anywhere.

### Dead build variants

`either/core_any.go` (`//go:build either_any`), `either/core_pointers.go`
(`//go:build either_pointers`) and their `idiomatic/result` counterparts
declare their own `Either` struct. Since `either/core.go` was refactored to
alias `common.Either` (commit `57a075a6`, 2026-09-02), both variants fail:

```
$ go build -tags either_pointers ./either/
either/core_pointers.go:21:6: Either redeclared in this block
        either/core.go:22:2: other declaration of Either
```

They were orphaned code and have been deleted. `Either` therefore has exactly
one implementation and methods need to be written once, in `internal/common`.

## Constraints

### Go version support

fp-go keeps supporting older toolchains (CI matrix: v1 on 1.20–1.27, v2 on
1.24–1.27, gen/v2 on 1.26–1.27). `go.mod` directives are not raised. Every
generic method and its tests live in a `*_go1_27.go` / `*_go1_27_test.go` file
starting with `//go:build go1.27`. Users on Go 1.24–1.26 keep the full
free-function API and simply do not see the methods. New `Monad*` helper
functions that methods delegate to are ordinary untagged code, available on
all supported versions.

### Only defined types can carry methods

fp-go v2 uses generic type aliases (see `DESIGN.md`, "Generic Type Aliases").
An alias of a function type cannot have methods.

| Type | Definition | Methods possible |
|------|-----------|------------------|
| `Option[A]` | `internal/common/option_core.go:40`, struct | yes, declared in `internal/common` |
| `Either[E, A]` | `internal/common/either_core.go:20`, struct | yes, declared in `internal/common` |
| `Result[A]` | `result/types.go:51`, `= Either[error, A]` | inherits `Either` methods |
| `Pair[L, R]` | `pair/types.go:61`, struct | yes |
| `Const[E, A]` | `constant/const.go:72`, struct | yes |
| `Builder[S, A]` | `optics/builder/types.go`, struct | yes |
| `Tuple1..15` | `tuple/`, structs | yes, low value |
| `Lens`, `Prism`, `Optional`, `Iso`, `codec.Type` | structs | already done |
| `IO`, `Reader*`, `IOEither`, `IOResult`, `State`, `Writer`, `Traversal`, `SeqResult` | aliases of `func` types (`io/types.go:21`: `IO[A any] = func() A`) | **no** |
| `idiomatic/*` | `(A, error)` / `(A, bool)` tuples | **no** |

Giving the function-typed monads a fluent API would require wrapping them in
structs, which contradicts the alias design and the performance rationale in
`CHAINING_PERFORMANCE_ANALYSIS.md`. Out of scope.

### Shapes that cannot be methods

- `Flatten`: needs receiver `Option[Option[A]]`. Use `Chain(F.Identity)`.
- `Ap`: needs a parameter or receiver of the form `T[func(A) B]`. A method on
  `Option[A]` whose signature mentions `Option[func(A) B]` makes the method set
  of `Option[A]` refer to `Option[func(A) B]`, whose `Ap` refers to
  `Option[func(func(A) B) C]`, and so on. The compiler rejects this:
  `instantiation cycle: A instantiated as func(A) B`. The same holds for
  `Either.Ap` and `Pair.Ap`. `Ap` stays a free function only.
- Type-class interfaces (`monad.Monad`, `internal/chain.Chainable`,
  `internal/functor.Functor`): interface methods stay non-generic, so these
  keep their lifted type parameters.

### Generic aliases in method signatures (Go 1.27.0 compiler bug)

A generic method whose signature uses a generic alias that refers back to the
receiver type (here `OptionKleisli[A, B] = func(A) Option[B]`) compiles in its
own package, but every importing package crashes the compiler with
`fatal error: all goroutines are asleep - deadlock!` in
`types2.(*Named).unpack` via `newAliasInstance`. Minimal reproduction:

```go
// package a
type O[A any] struct{ v A }
type K[A, B any] = func(A) O[B]
func (o O[A]) Chain[B any](f K[A, B]) O[B] { return f(o.v) }

// package b
import "aliasbug/a"
var X a.O[int] // go build: deadlock
```

Writing the function type literally (`f func(A) O[B]`) avoids it. Rule: method
signatures spell out `func(A) Option[B]`, `func(A) Either[E, B]`, etc., instead
of the `*Kleisli`/`*Operator` aliases. Aliases that do not mention the receiver
type (`Predicate[A]`) are unaffected.

### Performance

`internal/common/option.go:201` `OptionMonadMap` is branch-only with no
allocation. The curried `OptionMap` returns a closure.
`CHAINING_PERFORMANCE_ANALYSIS.md` identifies deep chains on struct-based
`Option` as the remaining hot spot. Methods therefore delegate to the
data-first `Monad*` functions, never to the curried forms. The existing codec
methods delegate through the curried form (`Pipe[O, I](ab)(t)`); codec
composition runs once when a codec is built, not per value, and the codec
package has no data-first variants, so they are left as is.

### Documentation visibility

Methods declared on `common.Option` are not listed under `option.Option` on
pkg.go.dev because `option.Option` is an alias. `Lens` and friends already
have this problem. Package docs (`option/doc.go`, `either/doc.go`,
`pair/doc.go`) must enumerate the methods.

## Design

### Principle

A method `x.Op(args...)` is defined as `MonadOp(x, args...)` and documented as
equivalent to `F.Pipe1(x, Op(args...))`. Names, argument order and laziness
match the free functions exactly. `DESIGN.md` gains a fourth form next to the
curried and monadic forms:

```go
Map[A, B any](f func(A) B) func(Option[A]) Option[B]      // curried, data last
MonadMap[A, B any](fa Option[A], f func(A) B) Option[B]   // monadic, data first
func (fa Option[A]) Map[B any](f func(A) B) Option[B]     // method, receiver first
```

Guidance: use free functions when the pipeline itself is the value
(`F.Flow`, reusable operators); use methods when transforming a value you
already hold.

### `Option[A]` — `internal/common/option_go1_27.go`

| Method | Delegates to |
|--------|--------------|
| `Map[B](f func(A) B) Option[B]` | `OptionMonadMap` |
| `MapTo[B](b B) Option[B]` | `OptionMonadMapTo` |
| `Chain[B](f func(A) Option[B]) Option[B]` | `OptionMonadChain` |
| `ChainFirst[B](f func(A) Option[B]) Option[A]` | `OptionMonadChainFirst` |
| `ChainTo[B](mb Option[B]) Option[B]` | `OptionMonadChainTo` |
| ~~`Ap`~~ | not expressible, see "Shapes that cannot be methods" |
| `Fold[B](onNone func() B, onSome func(A) B) B` | `OptionMonadFold` |
| `GetOrElse(onNone func() A) A` | `OptionMonadGetOrElse` |
| `Alt(that func() Option[A]) Option[A]` | `OptionMonadAlt` |
| `Filter(p func(A) bool) Option[A]` | new `OptionMonadFilter`; only `OptionFilter(p)` (curried) exists today |

Every type argument is inferred from the arguments. `Flatten` is omitted; its
doc comment on the free function points to `Chain(F.Identity)`.

### `Either[E, A]` — `internal/common/either_go1_27.go`

| Method | Delegates to |
|--------|--------------|
| `Map[B](f func(A) B) Either[E, B]` | `EitherMonadMap` |
| `MapTo[B](b B) Either[E, B]` | `EitherMonadMapTo` |
| `MapLeft[E2](f func(E) E2) Either[E2, A]` | `EitherMonadMapLeft` |
| `BiMap[E2, B](f func(E) E2, g func(A) B) Either[E2, B]` | `EitherMonadBiMap` |
| `Chain[B](f func(A) Either[E, B]) Either[E, B]` | `EitherMonadChain` |
| `ChainFirst[B](f func(A) Either[E, B]) Either[E, A]` | `EitherMonadChainFirst` |
| ~~`ChainTo`~~ | not a method: `EitherMonadChainTo` ignores its first argument by design (documented and tested), so a method would ignore its receiver and contradict `Option.ChainTo`, which propagates `None` |
| `ChainLeft[E2](f func(E) Either[E2, A]) Either[E2, A]` | `EitherMonadChainLeft` |
| ~~`Ap`~~ | not expressible, see "Shapes that cannot be methods" |
| `Fold[B](onLeft func(E) B, onRight func(A) B) B` | `EitherMonadFold` |
| `GetOrElse(onLeft func(E) A) A` | new `EitherMonadGetOrElse`; only `EitherGetOrElse(onLeft)` (curried) exists today |
| `Alt(that func() Either[E, A]) Either[E, A]` | `EitherMonadAlt` |
| `Swap() Either[A, E]` | `EitherSwap` |
| `ToOption() Option[A]` | `EitherToOption` |
| `Filter(p func(A) bool, empty E) Either[E, A]` | new `EitherMonadFilter`; only `EitherFilter(p, empty)` (curried) exists today |

`Result[A]` inherits all of these. `Result`-specific helpers that require
`E = error` (`UnwrapError`, `Eitherize*`) stay free functions because a method
cannot constrain a receiver type parameter.

`OrElse` is not added: the free function `EitherOrElse` is the curried form of
`ChainLeft`, and one method name per operation is enough.

### `Pair[L, R]` — `pair/pair_go1_27.go`

| Method | Delegates to |
|--------|--------------|
| `MapHead[L1](f func(L) L1) Pair[L1, R]` | `MonadMapHead` |
| `MapTail[R1](f func(R) R1) Pair[L, R1]` | `MonadMapTail` |
| `Map[R1](f func(R) R1) Pair[L, R1]` | `MonadMap` |
| `BiMap[L1, R1](f func(L) L1, g func(R) R1) Pair[L1, R1]` | `MonadBiMap` |
| `Swap() Pair[R, L]` | `Swap` |
| `Chain[R1](sg Semigroup[L], f func(R) Pair[L, R1]) Pair[L, R1]` | `MonadChain` |
| ~~`Ap`~~ | not expressible, see "Shapes that cannot be methods" |

### Optics completion

- ✅ `Iso.ComposePrism` (free function `IsoComposePrism` existed).
- ✅ Free functions `IsoComposeOptional`, `OptionalComposeIso`, public
  `iso.ComposeOptional`, `optional.ComposeIso`, and the methods
  `Iso.ComposeOptional`, `Optional.ComposeIso`. The composition matrix is
  complete for the four struct optics.
- `ComposeTraversal`: `Traversal` is a function alias and has no methods, but
  `Lens.ComposeTraversal` returning a `Traversal` is expressible. Requires
  locating the existing traversal composition in `optics/traversal` first;
  deferred until then.
- ✅ `optics/README.md` and `optics/doc.go` list the real method matrix.

### Cleanup

- Keep `//go:build go1.27` and the `*_go1_27.go` naming (see "Go version
  support").
- Delete `either/core_any.go`, `either/core_pointers.go`,
  `idiomatic/result/core_any.go`, `idiomatic/result/core_pointers.go`.

## Implementation Plan

Order, one PR each; `go vet ./... && go test ./...` at every step
(`TestRetryHttp` in the `*/http` packages hits the network and fails
independently of this work).

1. ✅ Cleanup (dead `Either` variants, optics doc drift).
2. ✅ `Option` methods + tests + benchmark. Gate passed (see "Benchmark
   results").
3. ✅ `Either` methods + tests (`Result` inherits them).
4. ✅ `Pair` methods + tests.
5. ✅ Optics matrix completion and doc fixes.
6. ✅ `DESIGN.md`, `option/doc.go`, `either/doc.go`, `result/doc.go`, `pair/doc.go`.

### Files

| Path | Change |
|------|--------|
| `internal/common/option_go1_27.go` | ✅ new, tagged |
| `option/option_go1_27_test.go`, `option/examples_methods_go1_27_test.go` | ✅ new, tagged; tests live in `option` so examples document the public type |
| `internal/common/option.go` | ✅ new untagged `OptionMonadFilter` |
| `internal/common/either_go1_27.go` | ✅ new, tagged |
| `either/either_go1_27_test.go`, `either/examples_methods_go1_27_test.go`, `result/examples_methods_go1_27_test.go` | ✅ new, tagged |
| `internal/common/either.go`, `either_filterable.go` | ✅ new untagged `EitherMonadGetOrElse`, `EitherMonadFilter` |
| `pair/pair_go1_27.go`, `pair_go1_27_test.go` | ✅ new, tagged |
| `internal/common/iso.go`, `lens_optional_compose.go`, `iso_optional_compose_test.go` | ✅ new untagged `IsoComposeOptional`, `OptionalComposeIso` + law tests |
| `optics/iso/optional.go`, `optics/optional/optional.go` | ✅ public `iso.ComposeOptional`, `optional.ComposeIso` |
| `internal/common/iso_go1_27.go`, `optional_go1_27.go` | ✅ `Iso.ComposePrism`, `Iso.ComposeOptional`, `Optional.ComposeIso` |
| `either/core_any.go`, `core_pointers.go`, `idiomatic/result/core_any.go`, `core_pointers.go` | ✅ deleted |
| `DESIGN.md`, `option/doc.go`, `either/doc.go`, `result/doc.go`, `pair/doc.go`, `optics/README.md`, `optics/doc.go` | ✅ docs |

### Tests

Per `AGENTS.md` "Checklist for New Code":

- `ExampleOption_Map`-style examples with `// Output:` for every method.
- Behavioral tests, not delegation tests: `Some`/`None` and `Left`/`Right`
  branches, laziness of `Alt`/`GetOrElse`/`Fold` callbacks (the untaken branch
  is not invoked), `ChainFirst` preserves the original value, `MapLeft`
  leaves `Right` untouched, `Swap` round-trips.
- Benchmark in `option`: three-step `Map → Chain → Fold` via methods versus
  nested `MonadMap`/`MonadChain`/`MonadFold` calls, run with `-benchmem`.
  Expected: identical ns/op and 0 B/op difference.

### Benchmark results

`go test ./option/ -bench 'Option(Methods|Monad|Pipe)_' -benchmem -count 5`,
Go 1.27.0, linux/amd64, mean of 5 runs:

| Benchmark | ns/op | B/op | allocs/op |
|-----------|------:|-----:|----------:|
| `Methods_MapChainFold` (`opt.Map(f).Chain(g).Fold(h, k)`) | 2.63 | 0 | 0 |
| `Monad_MapChainFold` (nested `MonadMap`/`MonadChain`/`MonadFold`) | 2.39 | 0 | 0 |
| `Pipe_MapChainFold` (`F.Pipe3(opt, Map(f), Chain(g), Fold(h, k))`) | 77.11 | 144 | 5 |
| `Methods_Chain5` | 19.32 | 0 | 0 |
| `Monad_Chain5` | 20.08 | 0 | 0 |

Methods cost the same as direct `Monad*` calls and are ~30× faster than the
curried `F.Pipe` style, which allocates one closure per operator.

For `Either` (`go test ./either/ -bench 'Either(Methods|Monad|Pipe)_'`), all
three forms measure 0.43 ns/op with 0 allocations: the curried `Either`
operators inline completely, so methods are at parity rather than faster.

Compatibility check: with `GOTOOLCHAIN=go1.24.0`, `go vet` and `go test` for
`option`, `internal/common` and `optics/...` pass (tagged files are excluded).

## Open Questions

1. ~~**`Ap` orientation.**~~ Resolved: `Ap` cannot be a method (instantiation
   cycle), see "Shapes that cannot be methods".
2. ~~**Docs for alias types.**~~ Resolved: package docs (`option`, `either`,
   `result`, `pair`) have a "Methods (Go 1.27+)" section listing the methods.
3. ~~**Dead build variants.**~~ Resolved: deleted.
4. **`Const`, `Builder`, `Tuple*`.** Expressible but low demand. Defer until
   requested.
5. **`Either.ChainTo` semantics.** `either.MonadChainTo(Left(e), mb)` returns
   `mb`, while `option.MonadChainTo(None, mb)` returns `None`. The Either
   behavior is documented and tested, but it diverges from the usual
   sequencing semantics (`*>`/`>>`). Changing it is a breaking behavior change
   and out of scope here.
6. **`ComposeTraversal`.** Deferred until the traversal composition in
   `optics/traversal` is mapped.

## Non-Goals

- Fluent API for function-aliased monads (`IO`, `Reader*`, `State`, ...).
- Replacing the data-last API or the `F.Pipe`/`F.Flow` style in existing
  examples.
- Changes to type-class interfaces or the `*/generic` HKT-emulation packages.
- Removing explicit instantiations enabled by generalized inference.
- Any change to the root (v1) module.
