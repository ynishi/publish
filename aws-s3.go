// Copyright 2017 Yutaka Nishimura. All rights reserved.
// Use of this source code is governed by a Apache License 2.0
// license that can be found in the LICENSE file.

package publish

import (
	"errors"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
)

type PublishAwsS3 struct {
	Publisher
	Svc  *s3.S3
	Conf *viper.Viper
}

type PublishAwsS3Opts struct {
	Bucket    string
	Key       string
	Accesskey string
	Secretkey string
	Token     string
	Region    string
}

func (p *PublishAwsS3) Publish(r io.Reader) error {

	if p.Conf == nil {
		return errors.New("error: conf is nil. pointer to viper is needed.")
	}
	var po *PublishAwsS3Opts
	var input *s3.PutObjectInput
	p.Conf.ReadInConfig()
	err := p.Conf.UnmarshalKey("AwsS3", &po)
	if err != nil {
		return err
	}
	err = p.Conf.UnmarshalKey("AwsS3POI", &input)
	if err != nil {
		return err
	}
	if po.Bucket != "" {
		input.Bucket = aws.String(po.Bucket)
	}
	if po.Key != "" {
		input.Key = aws.String(po.Key)
	}
	if input.Bucket == nil || input.Key == nil || po.Accesskey == "" || po.Secretkey == "" || po.Region == "" {
		return errors.New("error: cannot fetch conf vars.")
	}
	cred := credentials.NewStaticCredentials(
		po.Accesskey,
		po.Secretkey,
		po.Token,
	)
	if p.Svc == nil {
		sess, err := session.NewSession(&aws.Config{
			Credentials: cred,
			Region:      aws.String(po.Region),
		})
		if err != nil {
			return err
		}
		p.Svc = s3.New(sess)
	}

	input.Body = aws.ReadSeekCloser(r)

	_, err = p.Svc.PutObject(input)
	if err != nil {
		return err
	}
	return nil
}
