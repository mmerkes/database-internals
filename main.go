package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"syscall"
)

const KEY_SIZE int = 255
const VALUE_SIZE int = 255
const POINTER_SIZE int = 4
const CHECKSUM_SIZE int = 32
const HEADER_SIZE uint32 = 32 + 4 + 4

type PageHeader struct {
	Checksum    [CHECKSUM_SIZE]byte
	LowerOffset uint32
	UpperOffset uint32
}

func (header *PageHeader) Print() {
	fmt.Printf("Checksum: %v\nLowerOffset: %d\nUpperOffset: %d\n",
		header.Checksum, header.LowerOffset, header.UpperOffset)
}

func (header *PageHeader) Write(to []byte) {
	var i int
	for i = 0; i < CHECKSUM_SIZE; i++ {
		to[i] = header.Checksum[i]
	}

	binary.LittleEndian.PutUint32(to[i:], header.LowerOffset)
	i += 4

	binary.LittleEndian.PutUint32(to[i:], header.UpperOffset)
	i += 4
}

func (header *PageHeader) Read(from []byte) {
	header.Checksum = [CHECKSUM_SIZE]byte{}
	var i int
	for i = 0; i < CHECKSUM_SIZE; i++ {
		header.Checksum[i] = from[i]
	}

	header.LowerOffset = binary.LittleEndian.Uint32(from[i:])
	i += 4

	header.UpperOffset = binary.LittleEndian.Uint32(from[i:])
	i += 4
}

func main() {
	var stat syscall.Statfs_t
	if err := syscall.Statfs("/", &stat); err != nil {
		panic(err)
	}

	blockSize := stat.Bsize

	page := make([]byte, blockSize)

	checksum := sha256.Sum256([]byte("hello world\n"))

	header := PageHeader{
		Checksum:    checksum,
		LowerOffset: HEADER_SIZE,
		UpperOffset: blockSize,
	}
	fmt.Println("Original:")
	header.Print()

	header.Write(page)

	h2 := PageHeader{}

	h2.Read(page)

	fmt.Println("Result:")
	h2.Print()
}
