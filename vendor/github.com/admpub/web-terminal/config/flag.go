package config

import (
	"flag"
)

func FlagParse() {
	flag.StringVar(&Default.Listen, "listen", ":37079", "the port of http")
	flag.BoolVar(&Default.Debug, "debug", false, "show debug message.")
	flag.StringVar(&Default.MIBSDir, "mibs_dir", "", "set mibs directory.")
	flag.StringVar(&Default.SHExecute, "sh_execute", "bash", "the shell path")
	flag.StringVar(&Default.APPRoot, "url_prefix", "/", "url prefix")

	flag.StringVar(&Default.Password, "pw", "", "")
	flag.StringVar(&Default.IDFile, "i", "", "")
	flag.StringVar(&Default.SHFile, "f", "", "")

	flag.Parse()
}
