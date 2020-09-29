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
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

// 导出软著源码

var (
	sourceCodeFileIgnore     *string
	sourceCodeMaxLines       *uint
	sourceCodeFileDir        *string
	sourceCodeFileExtensions *string
	exportSourceCodeToFile   *string
)

var exportSourceCodeCmd = &cobra.Command{
	Use:   "exportsc",
	Short: `Export source files without comments`,
	RunE:  exportSourceCodeRunE,
}

func exportSourceCodeRunE(cmd *cobra.Command, args []string) error {
	extentions := strings.Split(strings.ToLower(*sourceCodeFileExtensions), `,`)
	f, err := os.Open(*exportSourceCodeToFile)
	if err != nil {
		if os.IsNotExist(err) {
			f, err = os.Create(*exportSourceCodeToFile)
		}
		if err != nil {
			return err
		}
	}
	defer f.Close()
	var re *regexp.Regexp
	re, err = regexp.Compile(*sourceCodeFileIgnore)
	if err != nil {
		return err
	}
	var root string
	root, err = filepath.Abs(*sourceCodeFileDir)
	if err != nil {
		return err
	}
	var lines uint
	err = filepath.Walk(root, func(pPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if re.MatchString(pPath) {
			return nil
		}
		if *sourceCodeMaxLines > 0 && lines > *sourceCodeMaxLines {
			return echo.ErrExit
		}
		extension := strings.ToLower(filepath.Ext(info.Name()))
		var ok bool
		for _, ext := range extentions {
			if extension == ext {
				ok = true
				break
			}
		}
		if !ok {
			return nil
		}
		var content string
		var commentStarted bool
		fmt.Println(`reading file:`, pPath)
		err = com.SeekFileLines(pPath, func(line string) error {
			lineClean := strings.TrimSpace(line)
			if len(lineClean) == 0 {
				return nil
			}
			if strings.HasPrefix(lineClean, `//`) {
				return nil
			}
			switch extension {
			case `.php`:
				if strings.HasPrefix(lineClean, `#`) {
					return nil
				}
			}
			if commentStarted {
				if strings.HasSuffix(lineClean, `*/`) {
					commentStarted = false
				}
				return nil
			}
			if strings.HasPrefix(lineClean, `/*`) {
				commentStarted = true
				return nil
			}
			content += line + "\n"
			lines++
			return nil
		})
		if err != nil {
			return err
		}
		_, err = f.WriteString("// Location: " + strings.TrimPrefix(pPath, root) + "\n" + content)
		return err
	})
	if err == echo.ErrExit {
		err = nil
	}
	return err
}

func init() {
	rootCmd.AddCommand(exportSourceCodeCmd)
	sourceCodeMaxLines = exportSourceCodeCmd.Flags().Uint("maxLines", 3500, "导出的最大行数")
	sourceCodeFileIgnore = exportSourceCodeCmd.Flags().String("ignore", "/vendor/|/license/|bindata_assetfs\\.go|/\\.git/|/dbschema/", "忽略文件或文件夹正则表达式")
	sourceCodeFileDir = exportSourceCodeCmd.Flags().String("dir", "./", "源文件所在文件夹")
	sourceCodeFileExtensions = exportSourceCodeCmd.Flags().String("ext", ".go", "源文件扩展名，多个扩展名用(,)分隔，例如: .go,.php")
	exportSourceCodeToFile = exportSourceCodeCmd.Flags().String("output", "./sourceCode.txt", "导出后保存的文件名称")
}
