package main

type Example struct {
	Strings []string `json:"strings,omitempty"` // Should not be flagged
	Things  []*Thing `json:"things,omitempty"`  // Should not be flagged
}

type Thing struct {
}

func main() {
	slice1 := []string{} // Should not be flagged
	_ = slice1
	slice2 := []*Thing{} // Should not be flagged
	_ = slice2
}
