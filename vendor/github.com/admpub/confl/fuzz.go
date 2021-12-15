// +build gofuzz

package confl

/*
Fuzz testing support files

https://github.com/dvyukov/go-fuzz

Usage:

    go-fuzz-build github.com/lytics/confl
    mkdir fuzz
    cp _examples/*.conf fuzz/
    go-fuzz -bin=confl-fuzz.zip -workdir=fuzz

See fuzz/crashers for results.
*/

func Fuzz(data []byte) int {
	var v map[string]interface{}
	if err := Unmarshal(data, &v); err != nil {
		return 0
	}
	return 1
}
