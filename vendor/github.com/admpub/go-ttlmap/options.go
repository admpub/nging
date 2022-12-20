package ttlmap

// Options for initializing a new Map.
type Options struct {
	InitialCapacity int
	OnWillExpire    func(key string, item Item)
	OnWillEvict     func(key string, item Item)
}

// KeyExistMode represents a restriction on the existence of a key for the
// operation to succeed.
type KeyExistMode int

const (
	// KeyExistDontCare can be used to ignore wether a key exists or not.
	KeyExistDontCare KeyExistMode = 0
	// KeyExistNotYet fails the operation if the key exists already.
	KeyExistNotYet KeyExistMode = 1
	// KeyExistAlready fails the opration if the key does not exist already.
	KeyExistAlready KeyExistMode = 2
)

// SetOptions for setting items on a Map.
type SetOptions struct {
	KeyExist KeyExistMode
}

func (opts *SetOptions) keyExist() KeyExistMode {
	if opts == nil {
		return KeyExistDontCare
	}
	return opts.KeyExist
}

// UpdateOptions for updating items on a Map.
type UpdateOptions struct {
	KeepValue      bool
	KeepExpiration bool
}
