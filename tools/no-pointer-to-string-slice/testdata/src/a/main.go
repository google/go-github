package main

type Example struct {
	Names []*string // want `use \[\]string instead of \[\]\*string in struct fields`
}

func main() {
	slice := []*string{} // want `use \[\]string instead of \[\]\*string`
	_ = slice
}
