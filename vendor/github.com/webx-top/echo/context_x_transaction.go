package echo

func (c *xContext) SetTransaction(t Transaction) {
	c.transaction = NewTransaction(t)
}

func (c *xContext) Transaction() Transaction {
	return c.transaction.Transaction
}

func (c *xContext) Begin() error {
	return c.transaction.Begin(c)
}

func (c *xContext) Rollback() error {
	return c.transaction.Rollback(c)
}

func (c *xContext) Commit() error {
	return c.transaction.Commit(c)
}

func (c *xContext) End(succeed bool) error {
	return c.transaction.End(c, succeed)
}
