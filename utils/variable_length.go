package utils

import "io"

// Variable length encoding, with add 1 encoding
// Only used to decode base object offset on OBJ_OFS_DELTA raw type object
func ReadVariableOffset(reader io.ByteReader) int64 {
	b, _ := reader.ReadByte()

	var v = int64(b & 0x7F)
	for b&0x80 > 0 {
		v++

		b, _ = reader.ReadByte()
		v = (v << 7) + int64(b&0x7F)
	}

	return v
}

// Variable size encoding, without 1 encoding, little endian.
// This is used for decoding deflated base object size and deflated object size on OBJ_REF_DELTA object
func ReadVariableSize(reader io.ByteReader) int64 {
	b, _ := reader.ReadByte()

	var v = int64(b & 0x7F)
	shift := 7
	for b&0x80 > 0 {

		b, _ = reader.ReadByte()
		v = int64(b&0x7F)<<shift + v
		shift += 7
	}

	return v
}
