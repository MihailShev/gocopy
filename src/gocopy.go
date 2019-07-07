package main

import (
	"flag"
	"io"
	"os"
)

import (
	"fmt"
)

var args = struct {
	from      string
	to        string
	offsetArg int64
	limit     uint
}{}

//nolint
func init() {
	flag.StringVar(&args.from, "from", "", "file to read from")
	flag.StringVar(&args.to, "to", "", "file to write")
	flag.Int64Var(&args.offsetArg, "offset", 0, "offset in bytes from the beginning of the file")
	flag.UintVar(&args.limit, "limit", 0, "limit in bytes to read from the file")
}

func main() {
	flag.Parse()
	dst, src, err := getSource()
	defer closeSource([]*os.File{dst, src})

	if err == nil {
		err = copySource(dst, src)
	}

	if err != nil {
		fmt.Println(err)
	}
}

func copySource(dst, src *os.File) error {
	buf, length, copyErr := makeBuf(dst)

	if copyErr != nil {
		return copyErr
	}

	offset := args.offsetArg

	if offset > length {
		fmt.Println("Offset is more then the file size or limit, 0 bytes copied")
	} else {
		printProgress(offset, 0)
	}

	for offset < length {
		read, err := dst.Read(buf[offset:])

		if err == nil {
			_, err = src.Write(buf[offset : int64(read)+offset])
			offset += int64(read)
			printProgress(length, offset)
		}

		if err != nil {
			if err != io.EOF {
				copyErr = err
			}
			break
		}
	}

	return copyErr
}

func printProgress(length, read int64) {
	var result float64

	if read == 0 {
		result = 0
	} else {
		result = float64(read) / float64(length) * 100
	}

	fmt.Printf("%s%6.2f%s", "progress:\t", result, "%\n")
}

func makeBuf(dst *os.File) (buf []byte, size int64, err error) {
	size, err = calcBufSize(dst)

	if err != nil {
		return nil, size, err
	}

	return make([]byte, size), size, err
}

func calcBufSize(file *os.File) (int64, error) {
	fileSize, err := getFileSize(file)
	limit := int64(args.limit)

	if limit > 0 && limit <= fileSize && err == nil {
		return limit, nil
	}

	return fileSize, err
}

func getFileSize(file *os.File) (int64, error) {
	stat, err := file.Stat()

	if err != nil {
		return 0, err
	}
	return stat.Size(), nil

}

func getSource() (dst, src *os.File, err error) {
	dst, err = os.Open(args.from)

	if err != nil {
		return
	}

	src, err = os.Create(args.to)
	return
}

func closeSource(s []*os.File) {
	for _, f := range s {
		if f != nil {
			err := f.Close()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
