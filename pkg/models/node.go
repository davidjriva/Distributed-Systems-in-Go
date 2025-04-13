package models

import (
	"io"
)

type Node struct {
	// Stdin is for reading messages in from the Maelstrom network.
	Stdin io.Reader

	// Stdout is for writing messages out to the Maelstrom network.
	Stdout io.Writer // contains filtered or unexported fields
}