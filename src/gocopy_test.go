package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestGetFileSize(t *testing.T) {
	const fileSize = 1000
	const fileName = "testFileSize"

	file := createFile(fileName, fileSize)

	got, _ := getFileSize(file)

	if got != fileSize {
		t.Error("Failed: expected file size", fileSize, "got", got)
	} else {
		t.Log("Successful: expected file size", fileSize, "got", got)
	}

	_ = os.Remove(fileName)
}

func TestGetSource(t *testing.T) {
	const fileSize = 1000
	args.from = "from.txt"
	args.to = "to.txt"

	createFile(args.from, fileSize)

	src, dst, _ := getSource()

	srcInfo, _ := src.Stat()
	dstInfo, _ := dst.Stat()

	srcName := srcInfo.Name()
	dstName := dstInfo.Name()

	if args.from != srcName {
		t.Error("Failed: expected src name", args.from, "got", srcName)
	} else {
		t.Log("Successful: expected src name", args.from, "got", srcName)
	}

	if args.to != dstName {
		t.Error("Failed: expected dst name", args.to, "got", dstName)
	} else {
		t.Log("Successful: expected dst name", args.to, "got", dstName)
	}

	_ = os.Remove(args.from)
	_ = os.Remove(args.to)

	args.to = ""
	args.from = ""
}

func TestCalcBufSize(t *testing.T) {
	const fileSize int64 = 2000
	const fileName = "testCalcBufSize"

	file := createFile(fileName, fileSize)

	// test without limit args
	got, _ := calcBufSize(file)

	if got != fileSize {
		t.Error("Failed: expected buffer size", fileSize, "got", got)
	} else {
		t.Log("Successful: expected buffer size", fileSize, "got", got)
	}

	// test with limit arg
	limitSize := int64(50)

	args.limit = uint(limitSize)

	got, _ = calcBufSize(file)

	if got != limitSize {
		t.Error("Failed: expected buffer size", limitSize, "got", got)
	} else {
		t.Log("Successful: expected buffer size", limitSize, "got", got)
	}

	// test with limit more then file size
	limitSize = fileSize * 2
	args.limit = uint(limitSize)
	got, _ = calcBufSize(file)

	if got != fileSize {
		t.Error("Failed: expected buffer size", fileSize, "got", got)
	} else {
		t.Log("Successful: expected buffer size", fileSize, "got", got)
	}

	_ = os.Remove(fileName)
	args.limit = 0
}

func TestMakeBuf(t *testing.T) {
	const fileSize int64 = 1000
	const fileName = "testMakeBuf"
	file := createFile(fileName, fileSize)

	buf, size, _ := makeBuf(file)

	expectedType := "[]uint8"

	got := fmt.Sprintf("%T", buf)

	if expectedType != got {
		t.Error("Failed: expected buffer type", expectedType, "got", got)
	} else {
		t.Log("Successful: expected buffer type", expectedType, "got", got)
	}

	expectedSize := int64(len(buf))

	if expectedSize != size {
		t.Error("Failed: expected buffer size", expectedSize, "got", size)
	} else {
		t.Log("Successful: expected buffer size", expectedSize, "got", size)
	}

	_ = os.Remove(fileName)
}

func TestCopySource(t *testing.T) {
	const dstName = "dst.txt"
	const srcName = "src.txt"
	const text = "\"test copy source function\""

	file, _ := os.Create(dstName)
	_, _ = file.Write([]byte(text))
	_ = file.Close()

	args.from = dstName
	args.to = srcName

	dst, src, _ := getSource()

	_ = copySource(dst, src)

	str, _ := ioutil.ReadFile(srcName)
	got := string(str)

	if got != text {
		t.Error("Failed: expected text from copy:", text, "got:", got)
	} else {
		t.Log("Successful: expected text from copy:", text, "got:", got)
	}

	_ = os.Remove(dstName)
	_ = os.Remove(srcName)
}

func createFile(name string, size int64) *os.File {
	file, _ := os.Create(name)
	buf := make([]byte, size)

	_, _ = file.Write(buf)

	return file
}
