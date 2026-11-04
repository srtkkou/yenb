/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package main

import (
	"fmt"
)

var (
	// Commit hash of VCS.
	CommitHash string
	// Commit date of VCS.
	CommitAt string
)

// Main
func main() {
	fmt.Println("yenb")
	fmt.Printf("CommitHash=%s\n", CommitHash)
	fmt.Printf("CommitAt=%s\n", CommitAt)
}
