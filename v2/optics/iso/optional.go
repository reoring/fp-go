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

package iso

import "github.com/IBM/fp-go/v2/internal/common"

// ComposeOptional composes an Iso[S, A] with an Optional[A, B] to produce an Optional[S, B].
//
// The Iso is used as a total lens: GetOption converts S to A and then applies the
// Optional; Set converts S to A, sets B through the Optional and converts back to S.
//
// Type Parameters:
//   - S: The new source type (introduced by the Iso)
//   - A: The intermediate type (target of the Iso; source of the Optional)
//   - B: The focus type of the resulting Optional
//
// Parameters:
//   - ab: Optional[A, B] that optionally focuses on B inside A
//
// Returns:
//   - A function that takes an Iso[S, A] and returns an Optional[S, B]
//
// See Also:
//   - ComposePrism: for composing with a Prism[A, B]
func ComposeOptional[S, A, B any](ab common.Optional[A, B]) func(Iso[S, A]) common.Optional[S, B] {
	return common.IsoComposeOptional[S](ab)
}
