package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pingcap/failpoint"
)

func main() {
	failpoint.Enable(_curpkg_("open-file-failed"), `return("injected")`)
	if err := WriteFile([]byte("hello failpoint")); err != nil {
		fmt.Println(err)
	}
	fmt.Println(failpoint.Status(_curpkg_("open-file-failed")))
}

func WriteFile(data []byte) error {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	filePath := filepath.Join(dir, "failpoint")
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	failpoint.Inject("open-file-failed", func(v failpoint.Value) error {
		return fmt.Errorf("open file failed: %v", v.(string))
	})

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	failpoint.Inject("write-file-failed", func(v failpoint.Value) error {
		return fmt.Errorf("write file failed: %v", v.(string))
	})

	return nil
}
