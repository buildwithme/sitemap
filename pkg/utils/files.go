package utils

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"strings"
)

func EnsureDirectories(path string) {
	index := strings.LastIndex(path, "/")
	if index <= 0 {
		return
	}

	path = path[:strings.LastIndex(path, "/")]

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
}

func WriteToFile(filePath string, dataSlice ...[]byte) error {
	// open output file
	fo, err := os.Create(filePath)
	if err != nil {
		log.Println(err)
		return err
	}

	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			log.Println(err)
			return
		}
	}()

	// make a write buffer
	w := bufio.NewWriter(fo)

	for i := range dataSlice {
		data := dataSlice[i]

		r := bytes.NewReader(data)

		// make a buffer to keep chunks that are read
		buf := make([]byte, 1024)
		for {
			// read a chunk
			n, err := r.Read(buf)
			if err != nil && err != io.EOF {
				log.Println(err)
				return err
			}
			if n == 0 {
				break
			}

			// write a chunk
			if _, err := w.Write(buf[:n]); err != nil {
				log.Println(err)
				return err
			}
		}
	}

	if err = w.Flush(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
