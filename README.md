# libarchive

Golang bindings for the [libarchive](http://libarchive.org) library.

## Simple Extraction Example
```golang
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	ar "github.com/CiscoSecurityServices/go-libarchive"
)

func printContents(filename string) {
	fmt.Println("file ", filename)
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error while opening file:\n %s\n", err)
		return
	}
	reader, err := ar.NewReader(file)
	if err != nil {
		fmt.Printf("Error on NewReader\n %s\n", err)
	}
	defer reader.Close()
	for {
		entry, err := reader.Next()
		if err == io.EOF {
			break
		}
		if errors.Is(err, ErrArchiveWarn) {
			// do something with the warning
			fmt.Println(err)
			continue
		}
		if err != nil {
			fmt.Printf("Error on reader.Next():\n%s\n", err)
			return
		}
		fmt.Printf("Name %s\n", entry.PathName())
		var buf bytes.Buffer
		size, err := buf.ReadFrom(reader)


		if err != nil {
			fmt.Printf("Error on reading entry from archive:\n%s\n", err)
		}
		if size > 0 {
			fmt.Println("Contents:\n***************", buf.String(), "*********************")
		}
	}
}

func main() {
	for _, filename := range os.Args[1:] {
		printContents(filename)
	}
}
```

## Docker
```
docker build -t lib .
docker run --rm -it -v $PWD:/app lib:latest go test ./...
```


## Acknowledgments
based on [mstoykov's go libarchive](https://github.com/mstoykov/go-libarchive)
which is based on
based on [robxu9's go-libarchive](https://github.com/robxu9/go-libarchive)

