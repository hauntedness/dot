// This is the package comment.
package p

import (
	json2 "encoding/json"
	"fmt"
)

// This comment is associated with the Greet function.
func Greet(who string) {
	fmt.Printf("Hello, %s!\n", who)
}

// Book is book
//
// book has its name and words
//
//go:generate go fmt ./...
type Book struct {
	// book name
	Name string
	// all words
	Words []string
	// Reader
	Reader *json2.Marshaler
}
