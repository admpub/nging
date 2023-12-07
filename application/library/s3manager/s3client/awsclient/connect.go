/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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

package awsclient

import (
	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func Connect(m *dbschema.NgingCloudStorage, bucketName string) (client *AWSClient, err error) {
	var sess *session.Session
	sess, err = NewSession(m)
	if err != nil {
		return
	}
	return &AWSClient{S3: s3.New(sess), session: sess, bucketName: bucketName}, nil
}

func NewSession(m *dbschema.NgingCloudStorage) (*session.Session, error) {
	isSecure := m.Secure == `Y`
	config := &aws.Config{
		DisableSSL:  aws.Bool(!isSecure),
		Endpoint:    aws.String(m.Endpoint),
		Credentials: credentials.NewStaticCredentials(m.Key, m.Secret, ""),
	}
	if len(m.Region) > 0 {
		config.Region = aws.String(m.Region)
	}
	return session.NewSession(config)
}
