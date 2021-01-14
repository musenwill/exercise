package main

import (
	"testing"

	"github.com/pingcap/failpoint"
	"github.com/stretchr/testify/assert"
)

func TestOk(t *testing.T) {
	assert.NoError(t, WriteFile(nil))
	assert.NoError(t, WriteFile([]byte("everything is ok")))
}

func TestOpenFailed(t *testing.T) {
	failpoint.Enable(_curpkg_("open-file-failed"), `return("injected")`)
	assert.Error(t, WriteFile([]byte("unable to open file")))
}

func TestWriteFailed(t *testing.T) {
	failpoint.Enable(_curpkg_("write-file-failed"), `return("injected")`)
	assert.Error(t, WriteFile([]byte("unable to open file")))
}
