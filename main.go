package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io"
	"log"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var (
	filename       = "./text.file"
	lastByte int64 = 0
	f        *os.File
)

func main() {
	// open file for reading
	var err error
	f, err = os.Open(filename)
	check(err)
	defer f.Close()

	//set pointer to the end
	lastByte, err = f.Seek(0, 2)
	check(err)

	//setup watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				//log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					//log.Println("modified file:", event.Name)
					logNewBytes()
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(filename)
	if err != nil {
		log.Fatal(err)
	}

	<-done
}

func logNewBytes() {
	data := make([]byte, 100)
	red := ""
	red_b := 0
	for {
		data = data[:cap(data)]
		n, err := f.Read(data)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("Error reading file: ", err)
			return
		}
		data = data[:n]
		red_b += n
		red = red + string(data)
	}
	if red_b == 0 {
		return
	}
	//	red = red
	fmt.Print(red)
}
