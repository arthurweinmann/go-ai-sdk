package utils

import (
	"fmt"
	"hash/crc64"
	"io"
	"strings"
)

// ComputeCRC64 computes the CRC64 hash of lowercased str
func ComputeCRC64(str string) (string, error) {
	str = strings.ToLower(str)

	// Create a new table with the ECMA polynomial
	table := crc64.MakeTable(crc64.ECMA)

	// Create a new hash interface
	hash := crc64.New(table)

	// Write our data to it and check for any error
	_, err := io.WriteString(hash, str)
	if err != nil {
		return "", err
	}

	// Return the hash as a hexadecimal string
	return fmt.Sprintf("%x", hash.Sum64()), nil
}
