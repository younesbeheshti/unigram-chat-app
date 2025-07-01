package utils

import "encoding/base64"


// Base64Encode
func Base64Encode(data []byte) []byte {
	dst := make([]byte, base64.RawStdEncoding.EncodedLen(len(data)))
	base64.RawStdEncoding.Encode(dst, data)
	return dst
}

// Base64Decode
func Base64Decode(data []byte) ([]byte, error) {
	dst := make([]byte, base64.RawStdEncoding.DecodedLen(len(data)))
	n, err := base64.RawStdEncoding.Decode(dst, data)
	if err !=  nil{
		return nil, err
	}
	return dst[:n], nil
}