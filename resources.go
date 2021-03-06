// Code generated by "esc -private -local-prefix-cwd -pkg=main -o=resources.go templates/"; DO NOT EDIT.

package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(path.Join("/Users/martin/go/src/github.com/programmfabrik/errorgen", f.local))
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	if !f.isDir {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is not directory", f.name)
	}

	fis, ok := _escDirs[f.local]
	if !ok {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is directory, but we have no info about content of this dir, local=%s", f.name, f.local)
	}
	limit := count
	if count <= 0 || limit > len(fis) {
		limit = len(fis)
	}

	if len(fis) == 0 && count > 0 {
		return nil, io.EOF
	}

	return fis[0:limit], nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// _escFS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func _escFS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// _escDir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func _escDir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// _escFSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func _escFSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		_ = f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// _escFSMustByte is the same as _escFSByte, but panics if name is not present.
func _escFSMustByte(useLocal bool, name string) []byte {
	b, err := _escFSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// _escFSString is the string version of _escFSByte.
func _escFSString(useLocal bool, name string) (string, error) {
	b, err := _escFSByte(useLocal, name)
	return string(b), err
}

// _escFSMustString is the string version of _escFSMustByte.
func _escFSMustString(useLocal bool, name string) string {
	return string(_escFSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/templates/codegen.go.tmpl": {
		name:    "codegen.go.tmpl",
		local:   "templates/codegen.go.tmpl",
		size:    3575,
		modtime: 1612787204,
		compressed: `
H4sIAAAAAAAC/5xXS4/bNhA+R79iqm4CKXAp5BrAh7abBntoYtRJe0gDhJZGXsISpZKUdw1B/73gU5It
u4+9rEnODL/55kW1ND/QPULfA7knG7cahihiddsIBUkEEJe1ivX/3UmhNL8UPqtMYd1WVKHZkUowvren
ouOK1ZgVuOv2eqfvfwBB+R7hztl9u4Y7ck8e7GoYtJqW8ufD4PWQF/o8jaIsg49HFE+CKQT1yCSoBhjP
q65A8GDILx3Pf6UtMK4aUI8IAnmBGlx0pEJ7qi/eCCzZMwzDJ6en1eSFlShSpxYvlN4J0Yj3yPUlKEqa
I/QRgNlOUrBcRC8+8ydB2yQF1AcRwIYKWsskHfX6AbIMfkPVCS6hNefQSSygbMQEujP+c1NguEBrfqA1
QgZ5UyA0pfHXX7ZVND8kKXz5quM2vUbqk7n4EEV972NUYKntruAOhbCYfcDusaRdpaQNWZZdUKPXTt/K
LBI4l5rQEQHAVbr9oYPZcgvSsGYAjnDN1frvPSptzsjCMCSpsW5UyKdRbHsmNhdK/c0uGQGGaFxNmEMh
VsYzz5eBfpstFMKem4VWbAXjqoT4pXwp46nCmex2LmzuirWNUWQzF7HseJlpdHCrAyGV6HIFNgwtFchV
yCfw6WlxDoM1Zk5sRtlMM+zMTG8uTE+CaPkyOPUPch5A0xVabcElzcaEBL6F5Y8wDN88u2bn3qvPQzbG
wGSdcOVAgeMTMC4V5bkpoyDyxNQjFDbn4UirDrW/ZcfzIJOk8Hok0Lpnovgq7PYWDBJL0xpMYySuQP87
IcBK5+jv47a2H7jS+Rui/v0xnoino5lAzfnKMgMhlMbhZJIo6UI36qeqF/eHzL1i0LbeJIXkdWjCvjev
bA6m7gplqutmIzdyrNSi6zVwVjlVo7y+aPMuQIMz/yW2nTv+CmsD1bfwYMUTRGyNTLTdSbjhAz4lMY9T
4jQNvESVYWNDhcTEjU/ySbB6MXjkIwzDCmL4k8dpeoPHs7Y/i4rLwBva8xHmCTcB0KQjcWHyBOv97+YM
t5SzPClrRbbWiyQ2ViFv6pZV2qxn5y28/Cs2xtN0wuGuMxE2zw3yU1eWKHwJCaHDR949Y94pTF7tunIF
SPxg/V+40Bj7F7gci7uuJFvDTxICMdbvAU+3KzjL/FA64Gnah/Q0Nj3Gj+YDnozGQpxmFpJ03hqXspTW
kowaE6eyzA8/B0ei+mcor6cZN8NynPfpNKBZgAFrOM6hzHFMOaG8mHTsvGlPsxcMJJ1E82jKHynjjO/T
q4Bvob3o5beBz3ql9WPyMDAuuegH7FVlBykqFBKohJq2i8W4+Fo8q+fWz2B/2x+CtrBnR+RhcC9RoMUS
XSSus164rWvIdbflQoo/c7qrUL/Bda+Evr9DPzFpJZAWJ8BnJk1hWUMknlZTML/WIBbnTpbBgwwtVXRo
evojQkulfiPbyDObJpLWeMXbBzn1ddc01pkX+GYFzcE0NiFIcu3dGbpKc7iorTdkOgrXukVONi6bR0kr
iVFwzw4aYHVbYY1c6cd5UbGdxSrJgwTZtfqT6Ipr82+MhfSw48m+Un1e/h0AAP//BVMfvPcNAAA=
`,
	},

	"/templates/doc.html.tmpl": {
		name:    "doc.html.tmpl",
		local:   "templates/doc.html.tmpl",
		size:    1043,
		modtime: 1603442925,
		compressed: `
H4sIAAAAAAAC/4SUwW7bMAyG73kKIuuxjZAWu7iMgQLeddsh2J2x6EiYbAmytjUw8u6DJdlxmgK5GBL/
jwLJnzCq0JpyhYpJlisA7MPJ8HgCOFh5giEeARrbhaeGWm1OBZDXZF6jco7fQAfDEOSMH6yX7AvYunfo
rdESDobq369ZdiSl7o4FfHXvt+9cP/JUW2PI9VzAdFqmbFxHLV/X+Y/1UYUCDtbIJUvFX93rwPIRaM6o
rbG+WJQ3sijyIFCk0eA4jXFQ2/Kb99b3KNS2XA0DeOqODA/8nVp+hAeuuIFiB5tqk0A4xwfVc4kEWu7W
w5BpOJ/XoDw3u/WXq2C5vKGgEoV6jva4JFXcbKqouRhWL+WeW2coMAr1klF0ni/8j8R7LhdZP8lTy4Fj
PykvepA2AMO0F+nmp2OUUjKKoK7D+5Pj22jFfe21C9p2t+Jb19lAn2sVN/THhKWAYipljM4lYkgmJehi
jcvWuGxNmkesPrvzSYPjVUJtqO9367hjyRY32xLkBzrp49v7e0B1D3i7B/y6AS5Tic1zJ6fmUMyTQZH9
vRAo8m6L+DP4HwAA//9hbeQpEwQAAA==
`,
	},

	"/templates": {
		name:  "templates",
		local: `templates/`,
		isDir: true,
	},
}

var _escDirs = map[string][]os.FileInfo{

	"templates/": {
		_escData["/templates/codegen.go.tmpl"],
		_escData["/templates/doc.html.tmpl"],
	},
}
