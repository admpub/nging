/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package license_test

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"golang.org/x/crypto/blake2b"

	//"github.com/admpub/highwayhash"
	//"github.com/admpub/metrohash"
	"github.com/admpub/snowflake"
	"github.com/admpub/sonyflake"
	"github.com/webx-top/com"
)

/*
最新Hash函数体验
警告：不要将非加密哈希函数用来处理来源不可信的数据，例如来自用户提交的数据
*/
var sf *sonyflake.Sonyflake

func init() {
	// 19位
	startTime, err := time.ParseInLocation(`2006-01-02 15:04:05`, `2018-08-08 08:08:08`, time.Local)
	if err != nil {
		panic(err)
	}
	st := sonyflake.Settings{
		StartTime: startTime,
		MachineID: func() (uint16, error) {
			return 12333, nil
		},
		CheckMachineID: func(id uint16) bool {
			return id == 12333
		},
	}
	sf = sonyflake.NewSonyflake(st)

	snowflake.Epoch = startTime.UnixNano() / 1e6 //设置纪元时间（毫秒），如不设置默认为2010年11月4日01:42:54的毫秒
}

func TestBlake2b(t *testing.T) {
	fmt.Println(`Blake2b`, `================================`)
	h, err := blake2b.New256(nil)
	if err != nil {
		panic(err)
	}
	// Write some data to the hash
	h.Write([]byte("Hello, World!!"))

	// Write some more data to the hash
	h.Write([]byte("How are you doing?"))

	// Get the current hash as a byte array
	b := h.Sum(nil)
	fmt.Println(hex.EncodeToString(b))
	//New512:  6d942c100b0abc3188836028a7d48b4551fda98237841766673fe3d37bdaf35f02aa7c0cc9977425d821e7d1dcfdfc8ac852375fb49d5cc38d6b8ee1d5e26071 (长度：128)
	//New256: 45b29b639d2dc6e3beb52970972c74af1c4ce08cbb19a047737efa768ce28e2d (长度：64)
}

// SonyFlake A distributed unique ID generator。灵感来自SnowFlake，一个分布式唯一ID生成器。
// StartTime是将Sonyflake时间定义为经过时间的时间。如果StartTime为0，则Sonyflake的开始时间设置为“2014-09-01 00:00:00 +0000 UTC”。如果StartTime超过当前时间，则不会创建Sonyflake。
// MachineID返回Sonyflake实例的唯一ID。如果MachineID返回错误，则不会创建Sonyflake。如果MachineID为nil，则使用默认的MachineID。默认MachineID返回私有IP地址的低16位。
// CheckMachineID验证机器ID的唯一性。如果CheckMachineID返回false，则不会创建Sonyflake。如果CheckMachineID为nil，则不进行验证。
func TestSonyFlake(t *testing.T) {
	fmt.Println(`SonyFlake`, `================================`)
	id, err := sf.NextID()
	if err != nil {
		panic(err)
	}
	fmt.Printf("ID: %v\n", id)
	com.Dump(sonyflake.Decompose(id))
}

/*/ MetroHash is a set of state-of-the-art hash functions for non-cryptographic use cases. They are notable for being algorithmically generated in addition to their exceptional performance. The set of published hash functions may be expanded in the future, having been selected from a very large set of hash functions that have been constructed this way.
// * Fastest general-purpose functions for bulk hashing.
// * Fastest general-purpose functions for small, variable length keys.
// * Robust statistical bias profile, similar to the MD5 cryptographic hash.
// * Hashes can be constructed incrementally (new)
// * 64-bit, 128-bit, and 128-bit CRC variants currently available.
// * Optimized for modern x86-64 microarchitectures.
// * Elegant, compact, readable functions.
//MetroHash 是一组用于非加密用例的最先进的哈希函数。
func TestMetroHash(t *testing.T) {

	fmt.Println(`MetroHash`, `================================`)
	// Create a new instance of the hash engine with default seed
	h := metrohash.NewMetroHash64()

	// Create a new instance of the hash engine with custom seed
	_ = metrohash.NewSeedMetroHash64(uint64(10))

	// Write some data to the hash
	h.Write([]byte("Hello, World!!"))

	// Write some more data to the hash
	h.Write([]byte("How are you doing?"))

	// Get the current hash as a byte array
	b := h.Sum(nil)
	fmt.Println(b)

	// Get the current hash as an integer (uint64) (little-endian)
	fmt.Println(h.Uint64()) // 20位数字

	// Get the current hash as a hexadecimal string (big-endian)
	fmt.Println(h.String())

	// Reset the hash
	h.Reset()

	// Output:
	// [205 190 61 93 89 212 164 71]
	// 14825354494498612295
	// cdbe3d5d59d4a447
}
*/
//Snowflake 一个分布式唯一ID生成器。
func TestSnowflake(t *testing.T) {
	fmt.Println(`Snowflake`, `================================`)
	// Create a new Node with a Node number of 1
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Generate a snowflake ID.
	id := node.Generate()

	// Print out the ID in a few different ways.
	fmt.Printf("Int64  ID: %d\n", id)
	fmt.Printf("String ID: %s\n", id)
	fmt.Printf("Base2  ID: %s\n", id.Base2())
	fmt.Printf("Base64 ID: %s\n", id.Base64())

	// Print out the ID's timestamp
	fmt.Printf("ID Time  : %d\n", id.Time())

	// Print out the ID's node number
	fmt.Printf("ID Node  : %d\n", id.Node())

	// Print out the ID's sequence number
	fmt.Printf("ID Step  : %d\n", id.Step())

	// Generate and print, all in one.
	fmt.Printf("ID       : %d\n", node.Generate().Int64())
	fmt.Printf("ID       : %d\n", node.Generate().Int64())
}

/*/ TestHighwayHash shows how to use HighwayHash-256 to compute fingerprints of files.
// HighwayHash 可用于防止散列泛滥攻击或验证短期消息。另外，它可以用作指纹识别功能。 HighwayHash不是通用加密哈希函数（例如Blake2b，SHA-3或SHA-2），如果需要强大的抗冲突性，则不应使用它。
func TestHighwayHash(t *testing.T) {
	fmt.Println(`HighwayHash`, `================================`)
	//highwayhash
	key, err := hex.DecodeString("000102030405060708090A0B0C0D0E0FF0E0D0C0B0A090807060504030201000")
	if err != nil {
		panic(err)
	}
	hash, err := highwayhash.New(key)
	if err != nil {
		panic(err)
	}

	hash.Write([]byte(`test`))

	// file, err := os.Open("./README.md") // specify your file here
	// if err != nil {
	// 	fmt.Printf("Failed to open the file: %v", err) // add error handling
	// 	return
	// }
	// defer file.Close()
	// if _, err = io.Copy(hash, file); err != nil {
	// 	fmt.Printf("Failed to read from file: %v", err) // add error handling
	// 	return
	// }

	checksum := hash.Sum(nil)
	fmt.Println(hex.EncodeToString(checksum))
	// 输出格式类似于: faaac029cdeeceacd4f74b1a392bcf5efb3183b9cec328d79fd6c1460608a608
}
*/
