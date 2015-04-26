package dds

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Format uint32

const (
	DXT1 = Format(0x31545844)
	DXT3 = Format(0x33545844)
	DXT5 = Format(0x35545844)
)

func LoadFile(filename string) (uint32, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return Load(bufio.NewReader(file))
}

func Load(r io.Reader) (uint32, error) {
	var err error

	var magic [4]byte

	_, err = io.ReadFull(r, magic[:])
	if err != nil {
		return 0, err
	}
	if string(magic[:]) != "DDS " {
		return 0, errors.New("Not DDS file")
	}

	var buf [124]byte
	_, err = io.ReadFull(r, buf[:])
	if err != nil {
		return 0, err
	}

	enc := binary.LittleEndian

	height := int32(enc.Uint32(buf[8:]))
	width := int32(enc.Uint32(buf[12:]))
	linearSize := enc.Uint32(buf[16:])
	mipMapCount := int32(enc.Uint32(buf[24:]))
	fourCC := Format(enc.Uint32(buf[80:]))

	bufsize := linearSize
	if mipMapCount > 1 {
		bufsize = 2 * linearSize
	}

	buffer := make([]byte, bufsize)
	n, _ := io.ReadFull(r, buffer)
	buffer = buffer[:n]

	format := uint32(0)
	switch fourCC {
	case DXT1:
		format = gl.COMPRESSED_RGBA_S3TC_DXT1_EXT
	case DXT3:
		format = gl.COMPRESSED_RGBA_S3TC_DXT3_EXT
	case DXT5:
		format = gl.COMPRESSED_RGBA_S3TC_DXT5_EXT
	default:
		return 0, fmt.Errorf("Unimplemented format 0x%x", fourCC)
	}

	var textureID uint32
	gl.GenTextures(1, &textureID)

	gl.BindTexture(gl.TEXTURE_2D, textureID)
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)

	blockSize := int32(8)
	if fourCC != DXT1 {
		blockSize = 16
	}
	offset := 0

	for level := int32(0); level < mipMapCount && (width > 0 || height > 0); level++ {
		size := ((width + 3) / 4) * ((height + 3) / 4) * blockSize
		gl.CompressedTexImage2D(gl.TEXTURE_2D, level, format, width, height, 0, size, gl.Ptr(buffer[offset:]))

		offset += int(size)
		width /= 2
		height /= 2

		if width < 1 {
			width = 1
		}
		if height < 1 {
			height = 1
		}
	}

	return textureID, nil
}
