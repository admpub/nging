/*
Package confl provides facilities for decoding and encoding TOML/NGINX configuration
files via reflection.
*/
package confl

import (
	"encoding/json"
	"fmt"
)

func printJSON(v interface{}) {
	b, _ := json.MarshalIndent(v, ``, `  `)
	fmt.Println(string(b))
}
