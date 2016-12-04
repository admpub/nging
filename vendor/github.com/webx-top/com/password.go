/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package com

// 生成哈希值
func Hash(str string) string {
	return Sha256(str)
}

// 盐值加密
func Salt() string {
	return Hash(RandStr(64))
}

// 创建密码
func MakePassword(password string, salt string) string {
	return Hash(salt + password)
}

// 检查密码(密码原文，数据库中保存的哈希过后的密码，数据库中保存的盐值)
func CheckPassword(rawPassword string, hashedPassword string, salt string) bool {
	return MakePassword(rawPassword, salt) == hashedPassword
}
