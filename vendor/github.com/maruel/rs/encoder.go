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

// Encoder can encode data into ecc codes.
type Encoder interface {
	// Encode calculates the ECC code for data and writes it into ecc.
	Encode(data []byte, ecc []byte)
}

type rSEncoder struct {
	r *gf256.RSEncoder
}

// NewEncoder generates a Reed-Solomon encoder that can generate ECC codes.
func NewEncoder(f *Field, c int) Encoder {
	return &rSEncoder{gf256.NewRSEncoder(f.f, c)}
}

func (r *rSEncoder) Encode(data []byte, ecc []byte) {
	r.r.ECC(data, ecc)
}
