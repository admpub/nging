/* Copyright 2012 Marc-Antoine Ruel. Licensed under the Apache License, Version
2.0 (the "License"); you may not use this file except in compliance with the
License.  You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0. Unless required by applicable law or
agreed to in writing, software distributed under the License is distributed on
an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
or implied. See the License for the specific language governing permissions and
limitations under the License. */

// Package rs implements Reed-Solomon error correcting codes.
//
// The code was inspired by ZXing's Java implementation but was reduced to only
// support 256 values space. Source: http://code.google.com/p/zxing/
//
// Much credit is due to Sean Owen, William Rucklidge since portions of this
// code are an indirect port of their Java or C++ Reed-Solomon implementations.
//
// Parts of ZXing's implementation have been replaced by Russ Cox's gf256
// library. Source: http://code.google.com/p/rsc/source/browse/#hg%2Fgf256
package rs

// BUG(maruel): It's far from complete. It can't accept total buffer size of
// more than 255 bytes. So len(data)+len(ecc) must be under 256 bytes.
