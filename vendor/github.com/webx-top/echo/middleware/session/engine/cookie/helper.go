package cookie

// KeyPairs Convert hashKey and blockKey to bytes
// @param hashKey
// @param blockKey
func KeyPairs(keys ...string) [][]byte {
	var hashKey, blockKey string
	if len(keys) > 0 {
		hashKey = keys[0]
	}
	if len(keys) > 1 {
		blockKey = keys[1]
	}
	keyPairs := [][]byte{}
	if len(hashKey) > 0 {
		keyPairs = append(keyPairs, []byte(hashKey))

		if len(blockKey) > 0 && blockKey != hashKey {
			keyPairs = append(keyPairs, []byte(blockKey))
		}
	}
	return keyPairs
}
