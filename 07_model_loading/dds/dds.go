package dds

import (
	"bufio"
	"errors"
	"io"
	"os"
)

type Format uint32

const (
	DXT1 Format = iota
	DXT3
	DXT5
)

type DDS struct {
	Height      uint32
	Width       uint32
	LinearSize  uint32
	MipMapCount uint32
	FourCC      uint32
	Format      Format
	Buffer      []byte
}

func LoadFile(filename string) (*DDS, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return Load(bufio.NewReader(file))
}

func Load(r io.Reader) (*DDS, error) {
	var err error

	var magic [4]byte

	_, err = io.ReadFull(r, magic[:])
	if err != nil {
		return nil, err
	}
	if string(magic[:]) != "DDS " {
		return nil, errors.New("Not DDS file")
	}

	var header [124]byte
	_, err = io.ReadFull(r, header[:])
	if err != nil {
		return nil, err
	}

	dds := &DDS{
		Height:      read32(header[8:]),
		Width:       read32(header[12:]),
		LinearSize:  read32(header[16:]),
		MipMapCount: read32(header[24:]),
		FourCC:      read32(header[80:]),
	}

	return dds, nil
}

func read32(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}
