package echo

// Get retrieves data from the context.
func (c *xContext) Get(key string, defaults ...interface{}) interface{} {
	return c.store.Get(key, defaults...)
}

// Set saves data in the context.
func (c *xContext) Set(key string, val interface{}) {
	c.store.Set(key, val)
}

// Delete saves data in the context.
func (c *xContext) Delete(keys ...string) {
	c.store.Delete(keys...)
}

func (c *xContext) Stored() Store {
	return c.store.CloneStore()
}
