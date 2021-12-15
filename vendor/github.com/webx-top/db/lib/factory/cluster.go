// added by swh@admpub.com
package factory

import (
	stdlog "log"

	"github.com/webx-top/db"
)

type masterLogger struct {
}

func (lg *masterLogger) Log(q *db.QueryStatus) {
	if q.Err != nil {
		Log.Error("<master>\n\t" + q.Stringify("\n\t") + "\n")
		return
	}
	if q.Slow {
		Log.Warn("<master>\n\t" + q.Stringify("\n\t") + "\n")
		return
	}
	Log.Info("<master>\n\t" + q.Stringify("\n\t") + "\n")
}

type slaveLogger struct {
}

func (lg *slaveLogger) Log(q *db.QueryStatus) {
	if q.Err != nil {
		Log.Error("<slave>\n\t" + q.Stringify("\n\t") + "\n")
		return
	}
	if q.Slow {
		Log.Warn("<slave>\n\t" + q.Stringify("\n\t") + "\n")
		return
	}
	Log.Info("<slave>\n\t" + q.Stringify("\n\t") + "\n")
}

var (
	DefaultMasterLog = db.Logger(&masterLogger{})
	DefaultSlaveLog  = db.Logger(&slaveLogger{})
)

// NewCluster : database cluster
func NewCluster() *Cluster {
	return &Cluster{
		masters:        []db.Database{},
		slaves:         []db.Database{},
		masterSelecter: &SelectFirst{},
		slaveSelecter:  &RoundRobin{},
	}
}

// Cluster : database cluster
type Cluster struct {
	masters        []db.Database
	slaves         []db.Database
	masterSelecter Selecter
	slaveSelecter  Selecter
}

// Master : write
func (c *Cluster) Master() (r db.Database) {
	length := len(c.masters)
	if length == 0 {
		panic(`Not connected to any database`)
	}
	if length > 1 {
		r = c.masters[c.masterSelecter.Select(length)]
	} else {
		r = c.masters[0]
	}
	return
}

// Slave : read-only
func (c *Cluster) Slave() (r db.Database) {
	length := len(c.slaves)
	if length == 0 {
		r = c.Master()
	} else {
		if length > 1 {
			r = c.slaves[c.slaveSelecter.Select(length)]
		} else {
			r = c.slaves[0]
		}
	}
	return
}

func (c *Cluster) SetSlaveSelecter(selecter Selecter) *Cluster {
	c.slaveSelecter = selecter
	return c
}

func (c *Cluster) SlaveSelecter() Selecter {
	return c.slaveSelecter
}

func (c *Cluster) SetMasterSelecter(selecter Selecter) *Cluster {
	c.masterSelecter = selecter
	return c
}

func (c *Cluster) MasterSelecter() Selecter {
	return c.masterSelecter
}

// AddMaster : added writable database
func (c *Cluster) AddMaster(databases ...db.Database) *Cluster {
	c.setMasterLogger(databases...)
	c.masters = append(c.masters, databases...)
	return c
}

// AddW : Deprecated
func (c *Cluster) AddW(databases ...db.Database) *Cluster {
	return c.AddMaster(databases...)
}

// AddR : Deprecated
func (c *Cluster) AddR(databases ...db.Database) *Cluster {
	return c.AddSlave(databases...)
}

func (c *Cluster) setMasterLogger(databases ...db.Database) *Cluster {
	for _, v := range databases {
		v.SetLogger(DefaultMasterLog)
	}
	return c
}

func (c *Cluster) setSlaveLogger(databases ...db.Database) *Cluster {
	for _, v := range databases {
		v.SetLogger(DefaultSlaveLog)
	}
	return c
}

// AddSlave : added read-only database
func (c *Cluster) AddSlave(databases ...db.Database) *Cluster {
	c.setSlaveLogger(databases...)
	c.slaves = append(c.slaves, databases...)
	return c
}

// SetMaster : set writable database
func (c *Cluster) SetMaster(index int, database db.Database) error {
	c.setMasterLogger(database)
	if len(c.masters) > index {
		c.masters[index] = database
		return nil
	}
	return ErrNotFoundKey
}

// SetSlave : set read-only database
func (c *Cluster) SetSlave(index int, database db.Database) error {
	c.setSlaveLogger(database)
	if len(c.masters) > index {
		c.slaves[index] = database
		return nil
	}
	return ErrNotFoundKey
}

// CloseAll : Close all connections
func (c *Cluster) CloseAll() {
	c.CloseMasters()
	c.CloseSlaves()
}

// CloseMasters : Close all master connections
func (c *Cluster) CloseMasters() {
	for _, database := range c.masters {
		if err := database.Close(); err != nil {
			stdlog.Println(err.Error())
		}
	}
}

// CloseSlaves : Close all slave connections
func (c *Cluster) CloseSlaves() {
	for _, database := range c.slaves {
		if err := database.Close(); err != nil {
			stdlog.Println(err.Error())
		}
	}
}

// CloseMaster : Close master connection
func (c *Cluster) CloseMaster(index int) bool {
	if len(c.masters) > index {
		if err := c.masters[index].Close(); err != nil {
			stdlog.Println(err.Error())
		}
		return true
	}
	return false
}

// CloseSlave : Close slave connection
func (c *Cluster) CloseSlave(index int) bool {
	if len(c.slaves) > index {
		if err := c.slaves[index].Close(); err != nil {
			stdlog.Println(err.Error())
		}
		return true
	}
	return false
}

func (c *Cluster) CountSlave() int {
	return len(c.slaves)
}

func (c *Cluster) CountMaster() int {
	return len(c.masters)
}
