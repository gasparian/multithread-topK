package io

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func checkValidPath(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		return err
	}
	return err
}

// ParseInputPath parses filepath from the stdin
func ParseInputPath() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter path to file: ")
	path, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	path = strings.TrimRight(path, "\n")
	err = checkValidPath(path)
	if err != nil {
		return "", err
	}
	return path, nil
}

// PrintResult prints slice of strings to stdout
func PrintResult(res []string) {
	for _, r := range res {
		fmt.Println(r)
	}
}

// FileSegmentPointer represents starting byte index and length of data segment in bytes
type FileSegmentPointer struct {
	Start int64
	Len   int64
}

// GetFileSegments reads file and returns segments pointers of ~`segmentSize` based on provided delimiter
// TODO: submit segments pointers in channel and return it, to not keep segments in memory
func GetFileSegments(f *os.File, bufSize int, segmentSize int64, delimiter byte) ([]FileSegmentPointer, error) {
	var (
		pointer     int64 = 0
		seek        int64 = 0
		start       int64 = 0
		chunkLength int64 = 0
		err         error
		n           int
		segment     FileSegmentPointer
	)
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	fsize := fi.Size()
	if segmentSize <= 0 {
		segmentSize = fsize
	}
	buf := make([]byte, bufSize)
	segments := make([]FileSegmentPointer, 0)
	for err == nil {
		chunkLength += segmentSize
		seek = pointer + chunkLength
		if seek >= fsize {
			segment.Start = pointer
			segment.Len = fsize - pointer - 1
			segments = append(segments, segment)
		}
		f.Seek(seek, 0)
		n, err = f.Read(buf)
		if n > 0 {
			for _, b := range buf[:n] {
				if b == delimiter && (seek+chunkLength) < (fsize-1) {
					start = pointer
					segment.Start = start
					segment.Len = chunkLength
					pointer += chunkLength + 1
					chunkLength = 0
					segments = append(segments, segment)
					break
				}
				chunkLength++
			}
		}
	}
	return segments, nil
}
