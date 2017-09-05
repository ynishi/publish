// Copyright 2017 Yutaka Nishimura. All rights reserved.
// Use of this source code is governed by a Apache License 2.0
// license that can be found in the LICENSE file.

package publish

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"strings"

	"fmt"

	"context"

	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/awstesting/unit"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
)

const (
	bucket_test     = `bkt`
	region_test     = `mock-region`
	endpoint_test   = `mock-region.test-amazonaws.com`
	awss3_toml_test = `
[AwsS3]
Bucket = "bkt"
Key = "doc.md"
Accesskey = "a"
Secretkey = "s"
Region = "mock-region"
[AwsS3POI]
ServerSideEncryption = "AES256"
ContentType = "text/plain"
`
)

type s3BucketTest struct {
	bucket  string
	url     string
	errCode string
}

func TestInitConfAwsS3(t *testing.T) {

	testAwsS3 := &PublishAwsS3Opts{
		Bucket:    "bkt",
		Key:       "doc.md",
		Accesskey: "a",
		Secretkey: "s",
		Region:    "mock-region",
	}
	testAwsS3POI := &s3.PutObjectInput{
		ServerSideEncryption: aws.String("AES256"),
		ContentType:          aws.String("text/plain"),
	}

	publishAwsS3 := &PublishAwsS3{}
	c := viper.New()
	c.SetConfigType("toml")
	c.ReadConfig(strings.NewReader(awss3_toml_test))

	InitConfAwsS3(publishAwsS3, c)

	if !reflect.DeepEqual(publishAwsS3.AwsS3, testAwsS3) {
		t.Fatalf("AwsS3 not matchted.\n want: %q,\n have: %q\n", testAwsS3, publishAwsS3.AwsS3)
	}
	if !reflect.DeepEqual(publishAwsS3.AwsS3POI, testAwsS3POI) {
		t.Fatalf("AwsS3POI not matchted.\n want: %q,\n have: %q\n", testAwsS3POI, publishAwsS3.AwsS3POI)
	}

}

func TestPublishAwsS3(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "PUT" {
				t.Fatalf("method not PUT: %q", r.Method)
			}
			if strings.Contains(r.URL.Host, endpoint_test) {
				t.Fatalf("endpoint not matched: %q", r.URL.Host)
			}
			if r.URL.Path != "/"+filename {
				t.Fatalf("filename not matched: %q", r.URL.Path)
			}
			fmt.Fprintln(w, "OK")
			return
		}))

	sess := unit.Session
	svc := s3.New(sess, &aws.Config{
		//require name to localhost below by change /etc/hosts or and so on.
		Endpoint:   aws.String(region_test + ".test-amazonaws.com:" + server.URL[17:]),
		DisableSSL: aws.Bool(true),
	})

	publishAwsS3 := &PublishAwsS3{}
	c := viper.New()
	c.SetConfigType("toml")
	c.ReadConfig(strings.NewReader(awss3_toml_test))

	InitConfAwsS3(publishAwsS3, c)
	r := strings.NewReader(content)
	publishAwsS3.Svc = svc

	errChan := make(chan error, 1)
	ctx := context.Background()

	go func() {
		errChan <- publishAwsS3.Publish(ctx, r)
	}()
	select {
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	case err := <-errChan:
		if err != nil {
			t.Fatal(err)
		}
	}

	ctxc, cancel := context.WithCancel(ctx)
	defer cancel()
	go publishAwsS3.Publish(ctxc, r)
	cancel()
	select {
	case <-ctxc.Done():
		// do nothing
	default:
		t.Fatal("failed cancel: %q", ctx)
	}
}
