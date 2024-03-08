package utils

import (
	"hash/crc64"
)

var crc64EcmaTable = crc64.MakeTable(crc64.ECMA)

// Crc64EcmaTable returns our standardized ECMA table
func Crc64EcmaTable() *crc64.Table {
	return crc64EcmaTable
}

// Crc64String returns the CRC-64 ECMA value of the string
func Crc64String(s string) uint64 {
	return crc64.Checksum([]byte(s), crc64EcmaTable)
}
