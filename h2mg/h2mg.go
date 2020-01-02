// Package h2mg provides support for 2MG Image Formats. Documentation at:
// http://apple2.org.za/gswv/a2zine/Docs/DiskImage_2MG_Info.txt
//
// This package doesn't attempt to implement complete support, just the
// bare minimum for the read/write functions I need.
package h2mg

import (
	"encoding/binary"
	"fmt"
)

// HeaderSize is the fixed number of bytes in a .2mg header v1
const HeaderSize = 64

// BlockSize is the size of a ProDOS disk block
const BlockSize = 512

const (
	// FormatDOS3 is the format ID for a DOS 3.3 image
	FormatDOS3 = 0
	// FormatProDOS is the format name for a ProDOS image
	FormatProDOS = 1
	// FormatNIB is the format name for a NIB image
	FormatNIB = 2
)

// Header2MG is the struct containing parsed .2MG header data
type Header2MG struct {
	Magic       string // Magic Number
	HeaderSize  uint16 // Size of header, in bytes
	Version     uint16 // Version number of 2MG format
	ImageFormat uint32 // Image Format Choices
	BlockCount  uint32 // Number of blocks in image
	Offset      uint32 // Offset to stored disk data
	Length      uint32 // Length of stored disk data
}

// Validate returns an error if the image is invalid, or nil if it
// correctly validates.
func (h2 Header2MG) Validate() error {
	if h2.Magic != "2IMG" {
		return fmt.Errorf("2MG Magic did not match (got %s)", h2.Magic)
	}
	if h2.HeaderSize != HeaderSize {
		return fmt.Errorf("2MG: header size incorrect (expected %d, got %d)", HeaderSize, h2.HeaderSize)
	}
	if h2.Version != 1 {
		return fmt.Errorf("2MG: only support version 1 (got %d)", h2.Version)
	}
	if h2.ImageFormat != FormatProDOS {
		return fmt.Errorf("2MG: only ProDOS format supported (got %d)", h2.ImageFormat)
	}
	if h2.BlockCount*BlockSize != h2.Length {
		return fmt.Errorf("2MG: block count and length did not match (expected length %d bytes, got %d)",
			h2.BlockCount*BlockSize, h2.Length)
	}

	return nil
}

// Parse2MG gets a .2mg disk header byte set and returns that data
// structure plus an error
func Parse2MG(data []uint8) (Header2MG, error) {
	var result Header2MG
	if len(data) < HeaderSize {
		return result, fmt.Errorf("2mg header too short: expected %d bytes, got %d", HeaderSize, len(data))
	}
	result.Magic = string(data[0x00 : 0x03+1]) // Magic Text
	result.HeaderSize = binary.LittleEndian.Uint16(data[0x08 : 0x09+1])
	result.Version = binary.LittleEndian.Uint16(data[0x0a : 0x0b+1])
	result.ImageFormat = binary.LittleEndian.Uint32(data[0x0c : 0x0f+1])
	result.BlockCount = binary.LittleEndian.Uint32(data[0x14 : 0x17+1])
	result.Offset = binary.LittleEndian.Uint32(data[0x18 : 0x1b+1])
	result.Length = binary.LittleEndian.Uint32(data[0x1c : 0x1f+1])

	return result, nil
}
