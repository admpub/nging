package pester

type ApplyOptions func(*Client)

func Concurrency(concurrency int) ApplyOptions {
	return func(c *Client) {
		c.Concurrency = concurrency
	}
}

func MaxRetries(maxRetries int) ApplyOptions {
	return func(c *Client) {
		c.MaxRetries = maxRetries
	}
}

func Backoff(mackoff BackoffStrategy) ApplyOptions {
	return func(c *Client) {
		c.Backoff = mackoff
	}
}

func KeepLog(keepLog bool) ApplyOptions {
	return func(c *Client) {
		c.KeepLog = keepLog
	}
}

func RetryOnHTTP429(retryOnHTTP429 bool) ApplyOptions {
	return func(c *Client) {
		c.RetryOnHTTP429 = retryOnHTTP429
	}
}