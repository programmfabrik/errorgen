package main

//go:generate ../errorgen -i simple.yml -o example.go

func main() {
	err := ErrFileNotFound().File("/tmp/testfile")
	println(err.Error())
}
