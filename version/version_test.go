package version_test

import (
	"testing"

	"web-layout/utils/version"
)

func TestInfo_Version(t *testing.T) {
	v := &version.Info{
		Name:   "Test_Server",
		Tag:    "0.65.2",
		Commit: "23fd345df",
		Branch: "master",
	}

	t.Log(v.Version())
}
