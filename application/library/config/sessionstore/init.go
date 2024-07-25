package sessionstore

import (
	_ "github.com/admpub/nging/v5/application/library/config/sessionstore/file"

	_ "github.com/admpub/nging/v5/application/library/config/sessionstore/redis"

	_ "github.com/admpub/nging/v5/application/library/config/sessionstore/bolt"

	_ "github.com/admpub/nging/v5/application/library/config/sessionstore/mysql"
)
