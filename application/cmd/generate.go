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
	"path/filepath"
	"strings"
	"time"

	"github.com/admpub/httpscerts"
	"github.com/admpub/license_gen/lib"
	"github.com/spf13/cobra"
)

// HTTPS 证书生成工具

var (
	https   *bool
	output  *string
	rsaBits *int
	certKey *string
	privKey *string
	host    *string
	expDate *string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate an SSL or RSA certificate",
	RunE:  generateRunE,
}

func generateRunE(cmd *cobra.Command, args []string) error {
	if len(*output) < 1 {
		*output = `.`
	}
	saveCertFile := filepath.Join(*output, *certKey)
	saveKeyFile := filepath.Join(*output, *privKey)
	if *https {
		tlsConfig := httpscerts.NewClassicConfig(strings.Split(*host, `,`)...)
		tlsConfig.RsaBits = *rsaBits
		if len(*expDate) > 0 {
			date, err := time.Parse("2006-01-02", *expDate)
			if err != nil {
				return err
			}
			tlsConfig.ValidFor = time.Now().Sub(date)
		}
		return httpscerts.Generate(tlsConfig, saveCertFile, saveKeyFile)
	}
	fmt.Println("Generating x509 Certificate")
	return lib.GenerateCertificate(saveCertFile, saveKeyFile, *rsaBits)
}

func init() {
	rootCmd.AddCommand(generateCmd)
	https = generateCmd.Flags().Bool("https", false, "Generate HTTPS Certificate")
	output = generateCmd.Flags().String("output", ".", "Output to directory")
	rsaBits = generateCmd.Flags().Int("bits", 2048, "Size of RSA key to generate")
	certKey = generateCmd.Flags().String("cert", "cert.pem", "Public certificate key")
	privKey = generateCmd.Flags().String("key", "key.pem", "Certificate key file")
	host = generateCmd.Flags().String("host", "localhost,127.0.0.1,::1", "Comma-separated hostnames and IPs to generate a certificate for")
	expDate = generateCmd.Flags().String("expiry", "", "Expiry date for the License. Expected format is 2006-01-02")
}
