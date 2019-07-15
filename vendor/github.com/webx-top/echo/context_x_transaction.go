package echo

func (c *xContext) SetTransaction(t Transaction) {
	c.transaction = t
}

func (c *xContext) Transaction() Transaction {
	return c.transaction
}

func (c *xContext) Begin() error {
	return c.Transaction().Begin(c)
}

func (c *xContext) Rollback() error {
	return c.Transaction().Rollback(c)
}

func (c *xContext) Commit() error {
	return c.Transaction().Commit(c)
}

func (c *xContext) End(succeed bool) error {
	return c.Transaction().End(c, succeed)
}
