/* Copyright 2012 Marc-Antoine Ruel. Licensed under the Apache License, Version
2.0 (the "License"); you may not use this file except in compliance with the
License.  You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0. Unless required by applicable law or
agreed to in writing, software distributed under the License is distributed on
an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
or implied. See the License for the specific language governing permissions and
limitations under the License. */

package rs

import (
	"github.com/maruel/rs/internal/gf256"
)

// The Galois Field for QR codes. See http://research.swtch.com/field for more
// information.
//
// x^8 + x^4 + x^3 + x^2 + 1
var QRCodeField256 = NewField(0x11D, 2)

// Field is a wrapper to gf256.Field so the type doesn't leak in.
type Field struct {
	f *gf256.Field
}

// NewField wraps gf256.NewField(). It is safe to use the premade
// QRCodeField256 all the time.
func NewField(poly int, α byte) *Field {
	return &Field{gf256.NewField(poly, int(α))}
}
