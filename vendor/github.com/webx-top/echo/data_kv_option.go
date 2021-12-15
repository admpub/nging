package echo

import "context"

type KVOption func(*KV)

func KVOptK(k string) KVOption {
	return func(a *KV) {
		a.K = k
	}
}

func KVOptV(v string) KVOption {
	return func(a *KV) {
		a.V = v
	}
}

func KVOptH(h H) KVOption {
	return func(a *KV) {
		a.H = h
	}
}

func KVOptHKV(k string, v interface{}) KVOption {
	return func(a *KV) {
		if a.H == nil {
			a.H = H{}
		}
		a.H.Set(k, v)
	}
}

func KVOptX(x interface{}) KVOption {
	return func(a *KV) {
		a.X = x
	}
}

func KVOptFn(fn func(context.Context) interface{}) KVOption {
	return func(a *KV) {
		a.fn = fn
	}
}
