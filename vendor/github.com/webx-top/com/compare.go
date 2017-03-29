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
	"regexp"
	"strconv"
)

var (
	regexpNotNumber = regexp.MustCompile(`[^0-9]+`)
)

const (
	// VersionCompareGt 左大于右
	VersionCompareGt = 1
	// VersionCompareEq 左等于右
	VersionCompareEq = 0
	// VersionCompareLt 左小于右
	VersionCompareLt = -1
)

// VersionCompare compare two versions in x.y.z form
// @param  {string} a     version string
// @param  {string} b     version string
// @return {int}          1 = a is higher, 0 = equal, -1 = b is higher
func VersionCompare(a, b string) (ret int) {
	as := regexpNotNumber.Split(a, -1)
	bs := regexpNotNumber.Split(b, -1)
	al := len(as)
	bl := len(bs)
	loopMax := bl

	if al > bl {
		loopMax = al
	}

	for i := 0; i < loopMax; i++ {
		var x, y string

		if al > i {
			x = as[i]
		}

		if bl > i {
			y = bs[i]
		}

		xi, _ := strconv.Atoi(x)
		yi, _ := strconv.Atoi(y)

		if xi > yi {
			ret = VersionCompareGt
		} else if xi < yi {
			ret = VersionCompareLt
		}

		if ret != 0 {
			break
		}
	}
	return
}
