package commons

import (
	"bytes"
	"compress/zlib"
	"io"
)

func zCompress(data []byte) (compressedData []byte, err error) {
	var buf bytes.Buffer
	writer := zlib.NewWriter(&buf)
	_, err = writer.Write(data)
	if err == nil {
		err = writer.Close()
		if err == nil {
			return buf.Bytes(), nil
		}
	}
	return
}

func zDecompress(compressedData []byte) (data []byte, err error) {
	reader, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err == nil {
		defer reader.Close()
		return io.ReadAll(reader)
	}
	return
}

func Compress(data []byte) (compressedBytes []byte, err error) {
	dataSize := len(data)
	compressedData, err := zCompress(data)
	if err == nil {
		compressedSize := len(compressedData)
		if dataSize <= compressedSize {
			compressedBytes = append([]byte{0}, data...)
		} else {
			compressedBytes = append([]byte{1}, compressedData...)
		}
	}
	return compressedBytes, err
}

func Decompress(compressedBytes []byte) (data []byte, err error) {
	preProcessedBytes := compressedBytes[1:]
	if compressedBytes[0] == 1 {
		return zDecompress(preProcessedBytes)
	}
	return preProcessedBytes, nil
}
