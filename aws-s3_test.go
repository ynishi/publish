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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/awstesting/unit"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
)

const (
	bucket_test     = `bkt`
	region_test     = `mock-region`
	endpoint_test   = `mock-region.test-amazonaws.com`
	awss3_test_toml = `
[AwsS3]
Bucket = "bkt"
Key = "doc.md"
Accesskey = "a"
Secretkey = "s"
Region = "mock-region"
[AwsS3POI]
ServerSideEncryption = "AES256"
Contenttype = "text/plain"
`
)

type s3BucketTest struct {
	bucket  string
	url     string
	errCode string
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

	publishAwsS3 := &PublishAwsS3{
		Conf: viper.New(),
	}
	publishAwsS3.Conf.SetConfigType("toml")
	publishAwsS3.Conf.ReadConfig(strings.NewReader(awss3_test_toml))

	ctx := context.Background()
	r := strings.NewReader(content)
	publishAwsS3.Svc = svc
	err := publishAwsS3.Publish(ctx, r)
	if err != nil {
		t.Fatal(err)
	}
}
