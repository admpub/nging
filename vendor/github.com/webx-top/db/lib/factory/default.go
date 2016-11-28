package factory

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
)

var DefaultFactory = New()

func Default() *Factory {
	return DefaultFactory
}

func Debug() bool {
	return DefaultFactory.Debug()
}

func SetDebug(on bool) *Factory {
	DefaultFactory.SetDebug(on)
	return DefaultFactory
}

func SetCacher(cc Cacher) *Factory {
	DefaultFactory.SetCacher(cc)
	return DefaultFactory
}

func AddDB(databases ...db.Database) *Factory {
	DefaultFactory.AddDB(databases...)
	return DefaultFactory
}

func AddSlaveDB(databases ...db.Database) *Factory {
	DefaultFactory.AddSlaveDB(databases...)
	return DefaultFactory
}

func AddDBToCluster(index int, databases ...db.Database) *Factory {
	DefaultFactory.AddDBToCluster(index, databases...)
	return DefaultFactory
}

func AddSlaveDBToCluster(index int, databases ...db.Database) *Factory {
	DefaultFactory.AddSlaveDBToCluster(index, databases...)
	return DefaultFactory
}

func SetCluster(index int, cluster *Cluster) *Factory {
	DefaultFactory.SetCluster(index, cluster)
	return DefaultFactory
}

func AddCluster(clusters ...*Cluster) *Factory {
	DefaultFactory.AddCluster(clusters...)
	return DefaultFactory
}

func GetCluster(index int) *Cluster {
	return DefaultFactory.GetCluster(index)
}

func Tx(param *Param) error {
	return DefaultFactory.Tx(param)
}

func NewTx(args ...int) (*Transaction, error) {
	return DefaultFactory.NewTx(args...)
}

func CloseAll() {
	DefaultFactory.CloseAll()
}

func Result(param *Param) db.Result {
	return DefaultFactory.Result(param)
}

// ================================
// API
// ================================

// Read ==========================

func SelectAll(param *Param) error {
	return DefaultFactory.SelectAll(param)
}

func SelectOne(param *Param) error {
	return DefaultFactory.SelectOne(param)
}

func SelectCount(param *Param) (int64, error) {
	return DefaultFactory.SelectCount(param)
}

func SelectList(param *Param) (func() int64, error) {
	return DefaultFactory.SelectList(param)
}

func Select(param *Param) sqlbuilder.Selector {
	return DefaultFactory.Select(param)
}

func All(param *Param) error {
	return DefaultFactory.All(param)
}

func List(param *Param) (func() int64, error) {
	return DefaultFactory.List(param)
}

func One(param *Param) error {
	return DefaultFactory.One(param)
}

func Count(param *Param) (int64, error) {
	return DefaultFactory.Count(param)
}

// Write ==========================

func Insert(param *Param) (interface{}, error) {
	return DefaultFactory.Insert(param)
}

func Update(param *Param) error {
	return DefaultFactory.Update(param)
}

func Upsert(param *Param, beforeUpsert ...func()) (interface{}, error) {
	return DefaultFactory.Upsert(param, beforeUpsert...)
}

func Delete(param *Param) error {
	return DefaultFactory.Delete(param)
}
