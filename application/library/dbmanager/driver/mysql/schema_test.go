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
	data := map[string][]string{
		"a[d]": []string{"first"},
		"a[e]": []string{"second"},
		"a[f]": []string{"third"},
		"a[g]": []string{"fourth"},
	}
	mx := &Mapx{}
	mx.Parse(data)
	com.Dump(mx)
	fmt.Println(`a[d] Value:`, mx.Value("a", "d"))
	fmt.Println(`a[e] Value:`, mx.Value("a", "e"))
	fmt.Println(`a[f] Value:`, mx.Value("a", "f"))
	fmt.Println(`a[g] Value:`, mx.Value("a", "g"))
}
