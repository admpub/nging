# qrcode
golang 编码/解码二维码  
1. 编码采用  github.com/boombuler/barcode
2. 解码支持采用golang或zbar

    zbar解码需要 `#include <zbar.h>` c语言库的支持（例如：pip install zbar）

例子
```golang
package main  

import (  
	"fmt"  
	"image/png"  
	"os"  

	"github.com/admpub/qrcode"
)

func main() {  
	err := png.EncodeToFile("test qrcode", 300, 300,"./test.png")  
	if err != nil {  
		fmt.Println(err)  
		return  
	}

	value := qrcode.DecodeFile("./test.png")  
	fmt.Println(value)
}
```