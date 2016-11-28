// added by swh@admpub.com
package factory

import (
	"log"
	"math/rand"

	"github.com/webx-top/db"
)

// NewCluster : database cluster
func NewCluster() *Cluster {
	return &Cluster{
		masters: []db.Database{},
		slaves:  []db.Database{},
	}
}

// Cluster : database cluster
type Cluster struct {
	masters []db.Database
	slaves  []db.Database
	prefix  string
}

// W : write
func (c *Cluster) W() (r db.Database) {
	length := len(c.masters)
	if length == 0 {
		panic(`Not connected to any database`)
	}
	if length > 1 {
		r = c.masters[rand.Intn(length-1)]
	} else {
		r = c.masters[0]
	}
	return
}

// Prefix : table prefix
func (c *Cluster) Prefix() string {
	return c.prefix
}

// Table : Table full name (including the prefix)
func (c *Cluster) Table(tableName string) string {
	return c.prefix + tableName
}

// SetPrefix : setting table prefix
func (c *Cluster) SetPrefix(prefix string) {
	c.prefix = prefix
}

// R : read-only
func (c *Cluster) R() (r db.Database) {
	length := len(c.slaves)
	if length == 0 {
		r = c.W()
	} else {
		if length > 1 {
			r = c.slaves[rand.Intn(length-1)]
		} else {
			r = c.slaves[0]
		}
	}
	return
}

// AddW : added writable database
func (c *Cluster) AddW(databases ...db.Database) *Cluster {
	c.masters = append(c.masters, databases...)
	return c
}

// AddR : added read-only database
func (c *Cluster) AddR(databases ...db.Database) *Cluster {
	c.slaves = append(c.slaves, databases...)
	return c
}

// SetW : set writable database
func (c *Cluster) SetW(index int, database db.Database) error {
	if len(c.masters) > index {
		c.masters[index] = database
		return nil
	}
	return ErrNotFoundKey
}

// SetR : set read-only database
func (c *Cluster) SetR(index int, database db.Database) error {
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
			log.Println(err.Error())
		}
	}
}

// CloseSlaves : Close all slave connections
func (c *Cluster) CloseSlaves() {
	for _, database := range c.slaves {
		if err := database.Close(); err != nil {
			log.Println(err.Error())
		}
	}
}

// CloseMaster : Close master connection
func (c *Cluster) CloseMaster(index int) bool {
	if len(c.masters) > index {
		if err := c.masters[index].Close(); err != nil {
			log.Println(err.Error())
		}
		return true
	}
	return false
}

// CloseSlave : Close slave connection
func (c *Cluster) CloseSlave(index int) bool {
	if len(c.slaves) > index {
		if err := c.slaves[index].Close(); err != nil {
			log.Println(err.Error())
		}
		return true
	}
	return false
}
