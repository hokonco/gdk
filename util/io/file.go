package io

import (
	"io/ioutil"
	"os"
	"path"
)

// FileReadString ...
func FileReadString(paths ...string) string {
	return string(FileReadBytes(paths...))
}

// FileReadBytes ...
func FileReadBytes(paths ...string) []byte {
	b, err := ioutil.ReadFile(path.Clean(path.Join(paths...)))
	die(err)

	return b
}

// FileWalkDir ...
func FileWalkDir(fn func(dir string, file os.FileInfo), paths ...string) {
	var f *os.File
	var fs []os.FileInfo
	var err error
	var path = path.Clean(path.Join(paths...))
	f, err = os.Open(path)
	die(err)

	fs, err = f.Readdir(-1)
	die(err)

	f.Close()
	for _, file := range fs {
		fn(path, file)
		if file.IsDir() {
			FileWalkDir(fn, path, file.Name())
		}
	}
	return
}

func die(err error) {
	if err != nil {
		panic(err)
	}
}
