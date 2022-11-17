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

import (
	"context"
	"fmt"
	"sync/atomic"
)

type Transaction interface {
	Begin(ctx context.Context) error
	Rollback(ctx context.Context) error
	Commit(ctx context.Context) error
	End(ctx context.Context, succeed bool) error
}

type UnwrapTransaction interface {
	Unwrap() Transaction
}

var (
	DefaultNopTransaction               = NewTransaction(nil)
	DefaultDebugTransaction Transaction = &DebugTransaction{}
)

type DebugTransaction struct {
}

func (b *DebugTransaction) Begin(ctx context.Context) error {
	fmt.Println(`DebugTransaction: Begin`)
	return nil
}

func (b *DebugTransaction) Rollback(ctx context.Context) error {
	fmt.Println(`DebugTransaction: Rollback`)
	return nil
}

func (b *DebugTransaction) Commit(ctx context.Context) error {
	fmt.Println(`DebugTransaction: Commit`)
	return nil
}

func (b *DebugTransaction) End(ctx context.Context, succeed bool) error {
	if succeed {
		return b.Commit(ctx)
	}
	return b.Rollback(ctx)
}

func NewTransaction(trans Transaction) *BaseTransaction {
	return &BaseTransaction{Transaction: trans}
}

type BaseTransaction struct {
	i int32
	Transaction
}

func (b *BaseTransaction) Begin(ctx context.Context) error {
	if b.Transaction == nil {
		return nil
	}
	if atomic.AddInt32(&b.i, 1) > 1 {
		return nil
	}
	return b.Transaction.Begin(ctx)
}

func (b *BaseTransaction) Rollback(ctx context.Context) error {
	if b.Transaction == nil {
		return nil
	}
	if atomic.LoadInt32(&b.i) <= 0 {
		return nil
	}
	atomic.SwapInt32(&b.i, 0)
	return b.Transaction.Rollback(ctx)
}

func (b *BaseTransaction) Commit(ctx context.Context) error {
	if b.Transaction == nil {
		return nil
	}
	if atomic.LoadInt32(&b.i) < 1 {
		panic(`transaction has already been committed or rolled back`)
	}
	if atomic.AddInt32(&b.i, -1) > 0 {
		return nil
	}
	return b.Transaction.Commit(ctx)
}

func (b *BaseTransaction) End(ctx context.Context, succeed bool) error {
	if succeed {
		return b.Commit(ctx)
	}
	return b.Rollback(ctx)
}

func (b *BaseTransaction) Unwrap() Transaction {
	return b.Transaction
}
