package echo

// Get retrieves data from the context.
func (c *xContext) Get(key string, defaults ...interface{}) interface{} {
	c.storeLock.RLock()
	v := c.store.Get(key, defaults...)
	c.storeLock.RUnlock()
	return v
}

// Set saves data in the context.
func (c *xContext) Set(key string, val interface{}) {
	c.storeLock.Lock()
	c.store.Set(key, val)
	c.storeLock.Unlock()
}

// Delete saves data in the context.
func (c *xContext) Delete(keys ...string) {
	c.storeLock.Lock()
	c.store.Delete(keys...)
	c.storeLock.Unlock()
}

func (c *xContext) Stored() Store {
	c.storeLock.Lock()
	copied := c.store.Clone()
	c.storeLock.Unlock()
	return copied
}
