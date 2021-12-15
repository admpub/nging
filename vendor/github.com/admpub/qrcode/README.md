# qrcode
golang 编码/解码二维码  
1. 编码采用  github.com/boombuler/barcode
2. 解码支持采用golang或zbar(zbar更准确)

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

# 如何安装zbar
## Ubuntu
```
sudo apt-get install libzbar-dev
```

## CentOS
```shell
sudo yum -y install epel-release pygtk2.x86_64 zbar-pygtk.x86_64 pygtk2-devel.x86_64 pygtk2-doc.noarch pygobject2.x86_64 pygobject2-devel.x86_64 pygobject2-doc.x86_64 gtk2 gtk2-devel gtk2-devel-docs pdftk ImageMagick ImageMagick-devel ghostscript Python-imaging python-devel python-gtk2-dev libqt4-dev PyQt4.x86_64 PyQt4-devel.x86_64

wget http://ftp.gnome.org/pub/GNOME/sources/pygtk/2.24/pygtk-2.24.0.tar.gz

tar -zvxf pygtk-2.24.0.tar.gz

cd pygtk-2.24.0

./configure

make

sudo make install

wget http://downloads.sourceforge.net/project/zbar/zbar/0.10/zbar-0.10.tar.gz

tar -zvxf zbar-0.10.tar.gz

cd zbar-0.10

./configure --disable-video

make

sudo make install
```