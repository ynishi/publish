// Copyright 2017 Yutaka Nishimura. All rights reserved.
// Use of this source code is governed by a Apache License 2.0
// license that can be found in the LICENSE file.

package publish

import (
	"errors"
	"io"

	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
)

// PublishAwsS3 is Publisher for aws s3.
type PublishAwsS3 struct {
	Publisher
	Svc      *s3.S3
	AwsS3    *PublishAwsS3Opts
	AwsS3POI *s3.PutObjectInput
}

// PublishAwsS3Opts is option for PublishAwsS3.
type PublishAwsS3Opts struct {
	Bucket    string
	Key       string
	Accesskey string
	Secretkey string
	Token     string
	Region    string
}

// InitConfAwsS3 initialize config and set option to PublishAwsS3.
func InitConfAwsS3(as3 *PublishAwsS3, c *viper.Viper) (err error) {
	if c == nil {
		return errors.New("error: conf is nil. pointer to viper is needed.")
	}
	err = c.UnmarshalKey("AwsS3", &as3.AwsS3)
	if err != nil {
		return err
	}
	err = c.UnmarshalKey("AwsS3POI", &as3.AwsS3POI)
	if err != nil {
		return err
	}
	return nil
}

// String return name of PublishAwsS3.
func (p *PublishAwsS3) String() string {
	return "PublishAwsS3"
}

// Publish for PublishAwsS3 publish document to aws s3.
func (p *PublishAwsS3) Publish(ctx context.Context, r io.Reader) (err error) {

	logger.Println("start publish aws s3")

	po := p.AwsS3
	input := p.AwsS3POI

	if po == nil {
		return errors.New("error: awss3 conf read failed.")
	}

	if input == nil {
		return errors.New("error: awss3poi conf read failed.")
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

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	_, err = p.Svc.PutObjectWithContext(ctx, input)
	if err != nil {
		return err
	}
	logger.Printf("aws s3 put: %s/%s", *input.Bucket, *input.Key)

	logger.Println("end publish aws s3")
	return nil
}
