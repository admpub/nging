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
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/s3manager/fileinfo"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

type AWSClient struct {
	*s3.S3
	session    *session.Session
	bucketName string
}

func (s *AWSClient) SetBucketName(bucketName string) *AWSClient {
	s.bucketName = bucketName
	return s
}

func (s *AWSClient) PutBucketCors(ctx echo.Context) error {
	rules := []*s3.CORSRule{
		// AllowedHeaders: []*string{aws.String(`*`)},
		// AllowedMethods: []*string{aws.String(`PUT`), aws.String(`POST`)},
		// AllowedOrigins: []*string{aws.String(`*`)},
		// ExposeHeaders:  []*string{aws.String(`ETag`)},
		//ID:*string,
		//MaxAgeSeconds:*int64,
	}
	rulesJSON := ctx.Form(`rules`)
	if len(rulesJSON) == 0 {
		return ctx.NewError(code.InvalidParameter, `CORS规则为JSON格式数据，不能为空`).SetZone(`rules`)
	}
	jsonBytes := com.Str2bytes(rulesJSON)
	err := json.Unmarshal(jsonBytes, &rules)
	if err != nil {
		return common.JSONBytesParseError(err, jsonBytes)
	}
	input := &s3.PutBucketCorsInput{
		Bucket: aws.String(s.bucketName),
		CORSConfiguration: &s3.CORSConfiguration{
			CORSRules: rules,
		},
	}
	output, err := s.S3.PutBucketCors(input)
	if err != nil {
		return err
	}
	log.Debug(output)
	return err
}

func (s *AWSClient) GetBucketCors() ([]*s3.CORSRule, error) {
	input := &s3.GetBucketCorsInput{
		Bucket: aws.String(s.bucketName),
	}
	output, err := s.S3.GetBucketCors(input)
	if err != nil {
		return nil, err
	}
	return output.CORSRules, err
}

func (s *AWSClient) CompleteMultipartUpload(ctx echo.Context, objectName string, uploadId string) error {
	objectName = strings.TrimPrefix(objectName, `/`)
	input := &s3.CompleteMultipartUploadInput{
		Bucket:          aws.String(s.bucketName),
		Key:             aws.String(objectName),
		MultipartUpload: &s3.CompletedMultipartUpload{},
		UploadId:        &uploadId,
	}
	input.MultipartUpload.Parts = []*s3.CompletedPart{}
	etags := ctx.FormValues(`etags`)
	var index int64
	for _, _etags := range etags {
		for _, etag := range strings.Split(_etags, `,`) {
			etag = strings.TrimSpace(etag)
			if len(etag) == 0 {
				continue
			}
			partNumber := int64(index) + 1
			input.MultipartUpload.Parts = append(input.MultipartUpload.Parts, &s3.CompletedPart{
				ETag:       aws.String(etag),
				PartNumber: aws.Int64(partNumber),
			})
			index++
		}
	}
	output, err := s.S3.CompleteMultipartUpload(input)
	if err != nil {
		return err
	}
	log.Debug(output)
	return err
}

func (s *AWSClient) ListPage(ctx echo.Context, objectPrefix string) (dirs []os.FileInfo, err error) {
	_, limit, pagination := common.PagingWithPosition(ctx)
	if limit < 1 {
		limit = 20
	}
	offset := ctx.Form(`offset`)
	prevOffset := ctx.Form(`prev`)
	var nextOffset string
	q := ctx.Request().URL().Query()
	q.Del(`offset`)
	q.Del(`prev`)
	q.Del(`_pjax`)
	pagination.SetURL(ctx.Request().URL().Path() + `?` + q.Encode() + `&offset={next}&prev={prev}`)
	input := &s3.ListObjectsInput{
		Bucket:    aws.String(s.bucketName),
		Prefix:    aws.String(objectPrefix),
		MaxKeys:   aws.Int64(int64(limit)),
		Delimiter: aws.String(`/`),
		Marker:    aws.String(offset),
	}
	var n int
	err = s.S3.ListObjectsPagesWithContext(ctx, input, func(p *s3.ListObjectsOutput, lastPage bool) bool {
		if p.NextMarker != nil {
			nextOffset = *p.NextMarker
		}
		for _, object := range p.CommonPrefixes {
			if object.Prefix == nil {
				continue
			}
			if len(objectPrefix) > 0 {
				key := strings.TrimPrefix(*object.Prefix, objectPrefix)
				object.Prefix = &key
			}
			if len(*object.Prefix) == 0 {
				continue
			}
			obj := fileinfo.NewStr(*object.Prefix)
			dirs = append(dirs, obj)
		}
		for _, object := range p.Contents {
			if object.Key == nil {
				continue
			}
			if len(objectPrefix) > 0 {
				key := strings.TrimPrefix(*object.Key, objectPrefix)
				object.Key = &key
			}
			if len(*object.Key) == 0 {
				continue
			}
			obj := fileinfo.NewS3(object)
			dirs = append(dirs, obj)
		}
		n += len(dirs)
		return n <= limit // continue paging
	})
	pagination.SetPosition(prevOffset, nextOffset, offset)
	ctx.Set(`pagination`, pagination)
	return
}

func (s *AWSClient) Upload(reader io.Reader, objectName string) (*s3manager.UploadOutput, error) {
	uploader := s3manager.NewUploader(s.session)
	return uploader.Upload(&s3manager.UploadInput{
		Body:   reader,
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objectName),
	})
}
