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
	if _, _err_ := failpoint.Eval(_curpkg_("open-file-failed")); _err_ == nil {
		return fmt.Errorf("open file failed")
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	if _, _err_ := failpoint.Eval(_curpkg_("write-file-failed")); _err_ == nil {
		return fmt.Errorf("write file failed")
	}

	return nil
}
