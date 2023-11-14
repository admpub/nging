// Copyright 2013 com authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package com

import (
	"errors"
	"math/rand"
	"reflect"
	"sort"
	"strings"
	"time"
)

func WithPrefix(strs []string, prefix string) []string {
	for k, v := range strs {
		strs[k] = prefix + v
	}
	return strs
}

func WithSuffix(strs []string, suffix string) []string {
	for k, v := range strs {
		strs[k] = v + suffix
	}
	return strs
}

// AppendStr appends string to slice with no duplicates.
func AppendStr(strs []string, str string) []string {
	for _, s := range strs {
		if s == str {
			return strs
		}
	}
	return append(strs, str)
}

// CompareSliceStr compares two 'string' type slices.
// It returns true if elements and order are both the same.
func CompareSliceStr(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}

	return true
}

// CompareSliceStrU compares two 'string' type slices.
// It returns true if elements are the same, and ignores the order.
func CompareSliceStrU(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		for j := len(s2) - 1; j >= 0; j-- {
			if s1[i] == s2[j] {
				s2 = append(s2[:j], s2[j+1:]...)
				break
			}
		}
	}
	if len(s2) > 0 {
		return false
	}
	return true
}

// IsSliceContainsStr returns true if the string exists in given slice.
func IsSliceContainsStr(sl []string, str string) bool {
	str = strings.ToLower(str)
	for _, s := range sl {
		if strings.ToLower(s) == str {
			return true
		}
	}
	return false
}

// IsSliceContainsInt64 returns true if the int64 exists in given slice.
func IsSliceContainsInt64(sl []int64, i int64) bool {
	for _, s := range sl {
		if s == i {
			return true
		}
	}
	return false
}

// ==============================
type reducetype func(interface{}) interface{}
type filtertype func(interface{}) bool

func InSlice(v string, sl []string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func InSet(v string, sl string, seperator ...string) bool {
	var sep string
	if len(seperator) > 0 {
		sep = seperator[0]
	}
	if len(sep) == 0 {
		sep = `,`
	}
	for _, vv := range strings.Split(sl, sep) {
		if vv == v {
			return true
		}
	}
	return false
}

func InSliceIface(v interface{}, sl []interface{}) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func InStringSlice(v string, sl []string) bool {
	return InSlice(v, sl)
}

func InInterfaceSlice(v interface{}, sl []interface{}) bool {
	return InSliceIface(v, sl)
}

func InIntSlice(v int, sl []int) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func InInt32Slice(v int32, sl []int32) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func InInt16Slice(v int16, sl []int16) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func InInt64Slice(v int64, sl []int64) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func InUintSlice(v uint, sl []uint) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func InUint32Slice(v uint32, sl []uint32) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func InUint16Slice(v uint16, sl []uint16) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func InUint64Slice(v uint64, sl []uint64) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func SliceRandList(min, max int) []int {
	if max < min {
		min, max = max, min
	}
	length := max - min + 1
	t0 := time.Now()
	rand.Seed(int64(t0.Nanosecond()))
	list := rand.Perm(length)
	for index := range list {
		list[index] += min
	}
	return list
}

func SliceMerge(slice1, slice2 []interface{}) (c []interface{}) {
	c = append(slice1, slice2...)
	return
}

func SliceReduce(slice []interface{}, a reducetype) (dslice []interface{}) {
	for _, v := range slice {
		dslice = append(dslice, a(v))
	}
	return
}

func SliceRand(a []interface{}) (b interface{}) {
	randnum := rand.Intn(len(a))
	b = a[randnum]
	return
}

func SliceSum(intslice []int64) (sum int64) {
	for _, v := range intslice {
		sum += v
	}
	return
}

func SliceFilter(slice []interface{}, a filtertype) (ftslice []interface{}) {
	for _, v := range slice {
		if a(v) {
			ftslice = append(ftslice, v)
		}
	}
	return
}

func SliceDiff(slice1, slice2 []interface{}) (diffslice []interface{}) {
	for _, v := range slice1 {
		if !InSliceIface(v, slice2) {
			diffslice = append(diffslice, v)
		}
	}
	return
}

func StringSliceDiff(slice1, slice2 []string) (diffslice []string) {
	for _, v := range slice1 {
		if !InSlice(v, slice2) {
			diffslice = append(diffslice, v)
		}
	}
	return
}

func SliceExtract(parts []string, recv ...*string) {
	recvEndIndex := len(recv) - 1
	if recvEndIndex < 0 {
		return
	}
	for index, value := range parts {
		if index > recvEndIndex {
			break
		}
		*recv[index] = value
	}
}

func UintSliceDiff(slice1, slice2 []uint) (diffslice []uint) {
	for _, v := range slice1 {
		if !InUintSlice(v, slice2) {
			diffslice = append(diffslice, v)
		}
	}
	return
}

func IntSliceDiff(slice1, slice2 []int) (diffslice []int) {
	for _, v := range slice1 {
		if !InIntSlice(v, slice2) {
			diffslice = append(diffslice, v)
		}
	}
	return
}

func Uint64SliceDiff(slice1, slice2 []uint64) (diffslice []uint64) {
	for _, v := range slice1 {
		if !InUint64Slice(v, slice2) {
			diffslice = append(diffslice, v)
		}
	}
	return
}

func Int64SliceDiff(slice1, slice2 []int64) (diffslice []int64) {
	for _, v := range slice1 {
		if !InInt64Slice(v, slice2) {
			diffslice = append(diffslice, v)
		}
	}
	return
}

func SliceIntersect(slice1, slice2 []interface{}) (diffslice []interface{}) {
	for _, v := range slice1 {
		if !InSliceIface(v, slice2) {
			diffslice = append(diffslice, v)
		}
	}
	return
}

func SliceChunk(slice []interface{}, size int) (chunkslice [][]interface{}) {
	if size >= len(slice) {
		chunkslice = append(chunkslice, slice)
		return
	}
	end := size
	for i := 0; i <= (len(slice) - size); i += size {
		chunkslice = append(chunkslice, slice[i:end])
		end += size
	}
	return
}

func SliceRange(start, end, step int64) (intslice []int64) {
	for i := start; i <= end; i += step {
		intslice = append(intslice, i)
	}
	return
}

func SlicePad(slice []interface{}, size int, val interface{}) []interface{} {
	if size <= len(slice) {
		return slice
	}
	for i := 0; i < (size - len(slice)); i++ {
		slice = append(slice, val)
	}
	return slice
}

func SliceUnique(slice []interface{}) (uniqueslice []interface{}) {
	for _, v := range slice {
		if !InSliceIface(v, uniqueslice) {
			uniqueslice = append(uniqueslice, v)
		}
	}
	return
}

func SliceShuffle(slice []interface{}) []interface{} {
	size := len(slice)
	for i := 0; i < size; i++ {
		a := rand.Intn(size)
		b := rand.Intn(size)
		if a == b {
			continue
		}
		slice[a], slice[b] = slice[b], slice[a]
	}
	return slice
}

var ErrNotSliceType = errors.New("expects a slice type")

// Shuffle 打乱数组
func Shuffle(arr interface{}) error {
	contentType := reflect.TypeOf(arr)
	if contentType.Kind() != reflect.Slice {
		return ErrNotSliceType
	}
	contentValue := reflect.ValueOf(arr)
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	len := contentValue.Len()
	for i := len - 1; i > 0; i-- {
		j := random.Intn(i + 1)
		x, y := contentValue.Index(i).Interface(), contentValue.Index(j).Interface()
		contentValue.Index(i).Set(reflect.ValueOf(y))
		contentValue.Index(j).Set(reflect.ValueOf(x))
	}
	return nil
}

func SliceInsert(slice, insertion []interface{}, index int) []interface{} {
	result := make([]interface{}, len(slice)+len(insertion))
	at := copy(result, slice[:index])
	at += copy(result[at:], insertion)
	copy(result[at:], slice[index:])
	return result
}

// SliceRemove SliceRomove(a,4,5) //a[4]
func SliceRemove(slice []interface{}, start int, args ...int) []interface{} {
	var end int
	if len(args) == 0 {
		end = start + 1
	} else {
		end = args[0]
	}
	if end > len(slice)-1 {
		return slice[:start]
	}
	return append(slice[:start], slice[end:]...)
}

func SliceGet(slice []interface{}, index int, defautls ...interface{}) interface{} {
	if index >= 0 && index < len(slice) {
		return slice[index]
	}
	if len(defautls) > 0 {
		if fn, ok := defautls[0].(func() interface{}); ok {
			return fn()
		}
		return defautls[0]
	}
	return nil
}

func StrSliceGet(slice []string, index int, defautls ...string) string {
	if index >= 0 && index < len(slice) {
		return slice[index]
	}
	if len(defautls) > 0 {
		return defautls[0]
	}
	return ``
}

func IntSliceGet(slice []int, index int, defautls ...int) int {
	if index >= 0 && index < len(slice) {
		return slice[index]
	}
	if len(defautls) > 0 {
		return defautls[0]
	}
	return 0
}

func UintSliceGet(slice []uint, index int, defautls ...uint) uint {
	if index >= 0 && index < len(slice) {
		return slice[index]
	}
	if len(defautls) > 0 {
		return defautls[0]
	}
	return 0
}

func Int32SliceGet(slice []int32, index int, defautls ...int32) int32 {
	if index >= 0 && index < len(slice) {
		return slice[index]
	}
	if len(defautls) > 0 {
		return defautls[0]
	}
	return 0
}

func Uint32SliceGet(slice []uint32, index int, defautls ...uint32) uint32 {
	if index >= 0 && index < len(slice) {
		return slice[index]
	}
	if len(defautls) > 0 {
		return defautls[0]
	}
	return 0
}

func Int64SliceGet(slice []int64, index int, defautls ...int64) int64 {
	if index >= 0 && index < len(slice) {
		return slice[index]
	}
	if len(defautls) > 0 {
		return defautls[0]
	}
	return 0
}

func Uint64SliceGet(slice []uint64, index int, defautls ...uint64) uint64 {
	if index >= 0 && index < len(slice) {
		return slice[index]
	}
	if len(defautls) > 0 {
		return defautls[0]
	}
	return 0
}

// SliceRemoveCallback : 根据条件删除
// a=[]int{1,2,3,4,5,6}
//
//	SliceRemoveCallback(len(a), func(i int) func(bool)error{
//		if a[i]!=4 {
//		 	return nil
//		}
//		return func(inside bool)error{
//			if inside {
//				a=append(a[0:i],a[i+1:]...)
//			}else{
//				a=a[0:i]
//			}
//			return nil
//		}
//	})
func SliceRemoveCallback(length int, callback func(int) func(bool) error) error {
	for i, j := 0, length-1; i <= j; i++ {
		if removeFunc := callback(i); removeFunc != nil {
			var err error
			if i+1 <= j {
				err = removeFunc(true)
			} else {
				err = removeFunc(false)
			}
			if err != nil {
				return err
			}
			i--
			j--
		}
	}
	return nil
}

func SplitKVRows(rows string, seperator ...string) map[string]string {
	sep := `=`
	if len(seperator) > 0 && len(seperator[0]) > 0 {
		sep = seperator[0]
	}
	res := map[string]string{}
	for _, row := range strings.Split(rows, StrLF) {
		parts := strings.SplitN(row, sep, 2)
		if len(parts) != 2 {
			continue
		}
		parts[0] = strings.TrimSpace(parts[0])
		if len(parts[0]) == 0 {
			continue
		}
		parts[1] = strings.TrimSpace(parts[1])
		res[parts[0]] = parts[1]
	}
	return res
}

func SplitKVRowsCallback(rows string, callback func(k, v string) error, seperator ...string) (err error) {
	sep := `=`
	if len(seperator) > 0 && len(seperator[0]) > 0 {
		sep = seperator[0]
	}
	for _, row := range strings.Split(rows, StrLF) {
		parts := strings.SplitN(row, sep, 2)
		if len(parts) != 2 {
			continue
		}
		parts[0] = strings.TrimSpace(parts[0])
		if len(parts[0]) == 0 {
			continue
		}
		parts[1] = strings.TrimSpace(parts[1])
		err = callback(parts[0], parts[1])
		if err != nil {
			return
		}
	}
	return
}

func JoinKVRows(value interface{}, seperator ...string) string {
	m, y := value.(map[string]string)
	if !y {
		return ``
	}
	sep := `=`
	if len(seperator) > 0 && len(seperator[0]) > 0 {
		sep = seperator[0]
	}
	r := make([]string, 0, len(m))
	for k, v := range m {
		r = append(r, k+sep+v)
	}
	sort.Strings(r)
	return strings.Join(r, "\n")
}

func TrimSpaceForRows(rows string) []string {
	rowSlice := strings.Split(rows, StrLF)
	res := make([]string, 0, len(rowSlice))
	for _, row := range rowSlice {
		row = strings.TrimSpace(row)
		if len(row) == 0 {
			continue
		}
		res = append(res, row)
	}
	return res
}
