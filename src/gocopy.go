package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var args = struct {
	from      string
	to        string
	offsetArg int64
	limit     uint
}{}

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

	if err != nil {
		fmt.Println(err)
		return
	}

	copySource(dst, src)
}

func copySource(dst, src *os.File) {
	buf, err := read(dst)

	if err == nil {
		_, err = src.Write(buf)
	} else {
		fmt.Println(err)
	}
}

func read(dst *os.File) ([]byte, error) {
	buf, length, readErr := makeBuf(dst)

	if readErr != nil {
		return nil, readErr
	}

	offset := args.offsetArg

	for offset < length {
		read, err := dst.Read(buf[offset:])

		if err != nil {
			if err != io.EOF {
				readErr = err
			}
			break
		}

		offset += int64(read)
	}

	return buf, readErr
}

func progres(length int64, readed int)  {
	fmt.Println(length / int64(readed) * 100)
}

func makeBuf(dst *os.File) ([]byte, int64, error) {
	size, err := calcBufSize(*dst)

	if err != nil {
		return nil, 0, err
	}

	return make([]byte, size), size, err
}

func calcBufSize(file os.File) (int64, error) {

	if args.limit > 0 {
		return int64(args.limit), nil
	} else {
		return getFileSize(file)
	}
}

func getFileSize(file os.File) (int64, error) {
	stat, err := file.Stat()

	if err != nil {
		return 0, err
	} else {
		return stat.Size(), nil
	}
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
