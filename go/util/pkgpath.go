package util

import "reflect"

type pkgpath struct{}

func PkgPath(name string) string {
	return reflect.TypeOf(pkgpath{}).PkgPath() + "/" + name
}
