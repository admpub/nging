/*

   Copyright 2019 Wenhui Shen <www.webx.top>

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

package echo

import "sync/atomic"

type Transaction interface {
	Begin() error
	Rollback() error
	Commit() error
	End(succeed bool) error
}

var (
	_                     Transaction = NewTransaction(nil)
	DefaultNopTransaction Transaction = &NopTransaction{}
)

type NopTransaction struct {
}

func (b *NopTransaction) Begin() error {
	return nil
}

func (b *NopTransaction) Rollback() error {
	return nil
}

func (b *NopTransaction) Commit() error {
	return nil
}

func (b *NopTransaction) End(succeed bool) error {
	if succeed {
		return b.Commit()
	}
	return b.Rollback()
}

func NewTransaction(trans Transaction) *BaseTransaction {
	return &BaseTransaction{Transaction: trans}
}

type BaseTransaction struct {
	i uint64
	Transaction
}

func (b *BaseTransaction) Begin() error {
	atomic.AddUint64(&b.i, 1)
	return b.Transaction.Begin()
}

func (b *BaseTransaction) Rollback() error {
	newValue := atomic.LoadUint64(&b.i) - 1
	atomic.SwapUint64(&b.i, newValue)
	if newValue > 0 {
		return nil
	}
	return b.Transaction.Rollback()
}

func (b *BaseTransaction) Commit() error {
	newValue := atomic.LoadUint64(&b.i) - 1
	atomic.SwapUint64(&b.i, newValue)
	if newValue > 0 {
		return nil
	}
	return b.Transaction.Commit()
}

func (b *BaseTransaction) End(succeed bool) error {
	if succeed {
		return b.Commit()
	}
	return b.Rollback()
}
