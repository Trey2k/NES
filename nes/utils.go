package nes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
)

func Hex16(n uint16) string {
	return fmt.Sprintf("%X", n)
}

func Hex8(n uint8) string {
	return fmt.Sprintf("%X", n)
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func readNextBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func elemSize(container interface{}) uintptr {
	return reflect.TypeOf(container).Elem().Size()
}

func readInto(file *os.File, data interface{}) {
	rawData := readNextBytes(file, int(elemSize(data)))
	buffer := bytes.NewBuffer(rawData)
	err := binary.Read(buffer, binary.BigEndian, data)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

}
