/*
  Copyright 2015 Adrian Stanescu. All rights reserved.
  Use of this source code is governed by the MIT License (MIT) that can be found in the LICENSE file.

  Go program that compares software versions in the x.y.z format
  Usage:
  x := "1"
  y := "1.0.1"
  z := "1.0"
  fmt.Println(VersionCompare(x, y)) // 1 = y
  fmt.Println(VersionCompare(x, z)) // 0 = equal
  fmt.Println(VersionCompare(x, a)) // -1 = x
*/

package com

import (
	"log"

	"github.com/hashicorp/go-version"
)

const (
	// VersionCompareGt 左大于右
	VersionCompareGt = 1
	// VersionCompareEq 左等于右
	VersionCompareEq = 0
	// VersionCompareLt 左小于右
	VersionCompareLt = -1
)

func SemVerCompare(a, b string) (int, error) {
	v1, err := version.NewVersion(a)
	if err != nil {
		return 0, err
	}

	v2, err := version.NewVersion(b)
	if err != nil {
		return 0, err
	}

	return v1.Compare(v2), nil
}

// VersionCompare compare two versions in x.y.z form
// @param  {string} a     version string
// @param  {string} b     version string
// @return {int}          1 = a is higher, 0 = equal, -1 = b is higher
func VersionCompare(a, b string) (ret int) {
	var err error
	ret, err = SemVerCompare(a, b)
	if err != nil {
		log.Printf(`failed to VersionCompare(%q, %q): %v`, a, b, err)
	}
	return
}

// VersionComparex compare two versions in x.y.z form
// @param  {string} a     version string
// @param  {string} b     version string
// @param  {string} op    <,<=,>,>=,= or lt,le,elt,gt,ge,egt,eq
// @return {bool}
func VersionComparex(a, b string, op string) bool {
	switch op {
	case `<`, `lt`:
		return VersionCompare(a, b) == VersionCompareLt
	case `<=`, `le`, `elt`:
		r := VersionCompare(a, b)
		return r == VersionCompareLt || r == VersionCompareEq
	case `>`, `gt`:
		return VersionCompare(a, b) == VersionCompareGt
	case `>=`, `ge`, `egt`:
		r := VersionCompare(a, b)
		return r == VersionCompareGt || r == VersionCompareEq
	case `=`, `eq`:
		return VersionCompare(a, b) == VersionCompareEq
	default:
		return false
	}
}
