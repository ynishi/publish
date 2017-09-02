// Copyright 2017 Yutaka Nishimura. All rights reserved.
// Use of this source code is governed by a Apache License 2.0
// license that can be found in the LICENSE file.

package cmd

import (
	"bytes"
	"testing"
)

func TestRootCmd(t *testing.T) {
	buf := new(bytes.Buffer)

	RootCmd.SetOutput(buf)
	RootCmd.SetArgs([]string{
		"--content=test.md",
		"--config=test.toml",
		"--timeout=120",
	})
	err := RootCmd.Execute()
	if err != nil {
		t.Fatalf("failed do RootCmd. %v", err)
	}
}
