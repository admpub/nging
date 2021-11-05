// added by swh@admpub.com

package factory

import (
	"context"
	"errors"
	"fmt"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
)

var (
	ErrNotFoundKey   = errors.New(`not found the key`)
	ErrNotFoundTable = errors.New(`not found the table`)
	ErrNotFoundField = errors.New(`not found the field`)
)

func New() *Factory {
	f := &Factory{
		databases: []*Cluster{},
		names:     map[string]int{},
	}
	f.Transaction = &Transaction{
		factory: f,
	}
	return f
}

type Factory struct {
	*Transaction
	databases []*Cluster
	names     map[string]int
	cacher    Cacher
}

func (f *Factory) Debug() bool {
	return db.DefaultSettings.LoggingEnabled()
}

func (f *Factory) CountCluster() int {
	return len(f.databases)
}

func (f *Factory) SetDebug(on bool, elapsedMs ...uint32) *Factory {
	db.DefaultSettings.SetLogging(on, elapsedMs...)
	for _, cluster := range f.databases {
		for _, master := range cluster.masters {
			master.SetLogging(on, elapsedMs...)
		}
		for _, slave := range cluster.slaves {
			slave.SetLogging(on, elapsedMs...)
		}
	}
	return f
}

func (f *Factory) SetCacher(cc Cacher) *Factory {
	f.cacher = cc
	return f
}

func (f *Factory) Cacher() Cacher {
	return f.cacher
}

// AddDB : add the master database
func (f *Factory) AddDB(databases ...db.Database) *Factory {
	if len(f.databases) > 0 {
		f.databases[0].AddMaster(databases...)
	} else {
		c := NewCluster()
		c.AddMaster(databases...)
		f.databases = append(f.databases, c)
	}
	return f
}

func (f *Factory) AddNamedDB(name string, databases ...db.Database) *Factory {
	if len(f.databases) > 0 {
		f.databases[0].AddMaster(databases...)
	} else {
		c := NewCluster()
		c.AddMaster(databases...)
		f.databases = append(f.databases, c)
	}
	f.names[name] = 0
	return f
}

// AddSlaveDB : add the slave database
func (f *Factory) AddSlaveDB(databases ...db.Database) *Factory {
	if len(f.databases) > 0 {
		f.databases[0].AddSlave(databases...)
	} else {
		c := NewCluster()
		c.AddSlave(databases...)
		f.databases = append(f.databases, c)
	}
	return f
}

// AddSlaveDB : add the slave database
func (f *Factory) AddNamedSlaveDB(name string, databases ...db.Database) *Factory {
	if len(f.databases) > 0 {
		f.databases[0].AddSlave(databases...)
	} else {
		c := NewCluster()
		c.AddSlave(databases...)
		f.databases = append(f.databases, c)
	}
	f.names[name] = 0
	return f
}

func (f *Factory) SetCluster(index int, cluster *Cluster, names ...string) *Factory {
	size := len(f.databases)
	if size > index {
		f.databases[index] = cluster
	} else if size == index {
		f.AddCluster(cluster)
	}
	if len(names) > 0 {
		f.names[names[0]] = index
	}
	return f
}

func (f *Factory) SetIndexName(index int, name string) *Factory {
	if len(f.databases) >= index {
		panic("[Factory.SetIndexName(" + fmt.Sprintf("%d, %q", index, name) + ")] index out of bounds: " + fmt.Sprintf("%d", len(f.databases)-1))
	}
	f.names[name] = index
	return f
}

func (f *Factory) AddCluster(clusters ...*Cluster) *Factory {
	f.databases = append(f.databases, clusters...)
	return f
}

func (f *Factory) AddNamedCluster(name string, clusters *Cluster) *Factory {
	f.names[name] = len(f.databases)
	f.databases = append(f.databases, clusters)
	return f
}

func (f *Factory) Cluster(index int) *Cluster {
	if len(f.databases) > index {
		return f.databases[index]
	}
	if index == 0 {
		panic(`Not connected to any database`)
	}
	return f.Cluster(0)
}

// GetCluster alias for Cluster
func (f *Factory) GetCluster(index int) *Cluster {
	return f.Cluster(index)
}

func (f *Factory) ClusterByName(name string) *Cluster {
	index, ok := f.names[name]
	if !ok {
		panic(`[Factory.ClusterByName] The cluster named ` + name + ` cannot be found`)
	}
	return f.Cluster(index)
}

func (f *Factory) IndexByName(name string) int {
	index, ok := f.names[name]
	if !ok {
		panic(`[Factory.IndexByName] The cluster named ` + name + ` cannot be found`)
	}
	return index
}

func (f *Factory) Tx(param *Param, ctx context.Context) error {
	if param.middlewareTx == nil {
		return nil
	}
	c := f.Cluster(param.index)
	trans := &Transaction{
		cluster: c,
		factory: f,
	}
	fn := func(tx sqlbuilder.Tx) error {
		trans.tx = tx
		return param.middlewareTx(trans)
	}
	if rdb, ok := c.Master().(sqlbuilder.Database); ok {
		return rdb.Tx(ctx, fn)
	}
	return db.ErrUnsupported
}

func (f *Factory) NewTx(ctx context.Context, args ...int) (trans *Transaction, err error) {
	var index int
	if len(args) > 0 {
		index = args[0]
	}
	c := f.Cluster(index)
	trans = &Transaction{
		cluster: c,
		factory: f,
	}
	if rdb, ok := c.Master().(sqlbuilder.Database); ok {
		trans.tx, err = rdb.NewTx(ctx)
	} else {
		err = db.ErrUnsupported
	}
	return
}

func (f *Factory) CloseAll() {
	for _, cluster := range f.databases {
		cluster.CloseAll()
	}
	f.databases = f.databases[0:0]
	for name := range f.names {
		delete(f.names, name)
	}
	f.Transaction = &Transaction{factory: f}
}
