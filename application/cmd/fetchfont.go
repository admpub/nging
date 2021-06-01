/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/marmot/miner"
)

// 下载 googleapis css文件中的字体

var (
	regexFontFile = regexp.MustCompile(`src\:(?:[ ]*local\('[^']+'\),)?[ ]*local\('([^']+)'\),[ ]*url\(([^\)]+)\)[ ]*format\('([^']+)'\)`)
	cssFile       *string
	fontOutput    *string
	fontDebug     *bool
)

var fetchFontCmd = &cobra.Command{
	Use:     "fetchfont",
	Short:   "Download the font in the googleapis css file",
	Example: "fetchfont https://fonts.googleapis.com/css?family=Covered+By+Your+Grace|Poppins:300,400,500,600,700",
	RunE:    fetchFontRunE,
}

func fetchFontRunE(cmd *cobra.Command, args []string) error {
	return FetchFont(*cssFile, *fontOutput, *fontDebug)
}

func FetchFont(_cssFile string, _outputDir string, _debug bool) error {
	worker, err := miner.NewWorker(nil)
	if err != nil {
		return err
	}
	var body []byte
	if com.IsURL(_cssFile) {
		worker.SetURL(_cssFile)
		body, err = worker.Get()
	} else {
		body, err = ioutil.ReadFile(_cssFile)
	}
	if err != nil {
		return err
	}
	if _debug {
		fmt.Println(string(body))
	}
	matches := regexFontFile.FindAllStringSubmatch(string(body), -1)
	err = com.MkdirAll(_outputDir, os.ModePerm)
	if err != nil {
		return err
	}
	for _, match := range matches {
		name := match[1]
		font := match[2]
		//format := match[3]
		fontBody, err := worker.SetURL(font).Get()
		if err == nil {
			destName := name + path.Ext(font)
			destFile := filepath.Join(_outputDir, destName)
			fmt.Println(font, `=>`, destFile)
			body = bytes.Replace(body, []byte(font), []byte(destName), -1)
			err = ioutil.WriteFile(destFile, fontBody, os.ModePerm)
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	destFile := filepath.Join(_outputDir, path.Base(_cssFile))
	err = ioutil.WriteFile(destFile, body, os.ModePerm)
	if _debug {
		echo.Dump(matches)
	}
	return err
}

func init() {
	rootCmd.AddCommand(fetchFontCmd)
	cssFile = fetchFontCmd.Flags().String("cssfile", "https://fonts.googleapis.com/css?family=Covered+By+Your+Grace|Poppins:300,400,500,600,700", "Generate HTTPS Certificate")
	fontOutput = fetchFontCmd.Flags().String("output", ".", "Output to directory")
	fontDebug = fetchFontCmd.Flags().Bool("debug", false, "Debug mode")
}
