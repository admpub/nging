package mysql

import (
	"testing"

	"fmt"

	"github.com/webx-top/com"
)

func TestMapx(t *testing.T) {
	com.Dump(ParseFormName("a[b][c][d]"))
	com.Dump(ParseFormName("a[[b][c][d]"))
	com.Dump(ParseFormName("a][[b][c][d]"))
	com.Dump(ParseFormName("a[][b][c][d]"))
	data := map[string][]string{
		"a[d]":   []string{"first"},
		"a[e]":   []string{"second"},
		"a[f]":   []string{"third"},
		"a[g]":   []string{"fourth"},
		"b[]":    []string{"index_0", "index_1"},
		"c[][a]": []string{"index 0.a"},
		"c[][b]": []string{"index 1.b"},
	}
	mx := &Mapx{}
	mx.Parse(data)
	com.Dump(mx)
	fmt.Println(`a[d] Value: "first" ==`, mx.Value("a", "d"))
	fmt.Println(`a[e] Value: "second" ==`, mx.Value("a", "e"))
	fmt.Println(`a[f] Value: "third" ==`, mx.Value("a", "f"))
	fmt.Println(`a[g] Value: "fourth" ==`, mx.Value("a", "g"))
	fmt.Println(`b[] Value: [index_0 index_1] ==`, mx.Values("b"))
	fmt.Println(`c[][a] Value: "index 0.a" ==`, mx.Value("c", "0", "a"))
	fmt.Println(`c[][b] Value: "index 1.b" ==`, mx.Value("c", "1", "b"))
}
