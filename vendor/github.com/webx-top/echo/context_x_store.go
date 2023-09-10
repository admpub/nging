package echo

import "github.com/webx-top/echo/param"

// Get retrieves data from the context.
func (c *xContext) Get(key string, defaults ...interface{}) interface{} {
	return c.store.Get(key, defaults...)
}

// Set saves data in the context.
func (c *xContext) Set(key string, val interface{}) {
	c.store.Set(key, val)
}

// Incr Increment the value and return the new value
func (c *xContext) Incr(key string, n interface{}, defaults ...interface{}) int64 {
	return c.store.Incr(key, param.AsInt64(n), defaults...)
}

// Decr Decrement the value and return the new value
func (c *xContext) Decr(key string, n interface{}, defaults ...interface{}) int64 {
	return c.store.Decr(key, param.AsInt64(n), defaults...)
}

// Delete saves data in the context.
func (c *xContext) Delete(keys ...string) {
	c.store.Delete(keys...)
}

func (c *xContext) Stored() Store {
	return c.store.CloneStore()
}
