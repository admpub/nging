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

package codec

import "crypto/des"

var (
	_ Codec = NewDesCBCCrypto()
	_ Codec = NewDesECBCrypto()
)

func DesGenKey(key []byte) []byte {
	size := des.BlockSize
	kkey := make([]byte, 0, size)
	ede2Key := []byte(key)
	length := len(ede2Key)
	if length > size {
		kkey = append(kkey, ede2Key[:size]...)
	} else {
		div := size / length
		mod := size % length
		for i := 0; i < div; i++ {
			kkey = append(kkey, ede2Key...)
		}
		kkey = append(kkey, ede2Key[:mod]...)
	}
	return kkey
}
