/*
 * Copyright 1999-2020 Alibaba Group Holding Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package nacos_error

import (
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/common/constant"
)

type NacosError struct {
	errorCode   string
	errMsg      string
	originError error
}

func NewNacosError(errorCode string, errMsg string, originError error) *NacosError {
	return &NacosError{
		errorCode:   errorCode,
		errMsg:      errMsg,
		originError: originError,
	}

}

func (err *NacosError) Error() (str string) {
	nacosErrMsg := fmt.Sprintf("[%s] %s", err.ErrorCode(), err.errMsg)
	if err.originError != nil {
		return nacosErrMsg + "\ncaused by:\n" + err.originError.Error()
	}
	return nacosErrMsg
}

func (err *NacosError) ErrorCode() string {
	if err.errorCode == "" {
		return constant.DefaultClientErrorCode
	} else {
		return err.errorCode
	}
}
