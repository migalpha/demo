package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	end           int = 16777215
	numGoRoutines int = 10
)

// Simple approach the file is written sequentially.
// func main() {
// 	f, err := os.Create("numbers.txt")
// 	if err != nil {
// 		log.Fatalf("err: %s", err.Error())
// 	}
// 	defer f.Close()

// 	for i := 0; i < end+1; i++ {
// 		f.WriteString(fmt.Sprintf("%06x\n", i))
// 	}
// }

// second approach using go routines but bottleneck writing file.
// func main() {
// 	f, err := os.Create("numbers2.txt")
// 	if err != nil {
// 		log.Fatalf("err: %s", err.Error())
// 	}
// 	defer f.Close()

// 	doneCh := make(chan struct{})
// 	for i := 0; i <= end; i = i + (end / numGoRoutines) + 1 {
// 		next := i + (end / numGoRoutines)
// 		if next > end {
// 			next = end
// 		}
// 		go writeString(i, next, f, doneCh)
// 	}

// 	var numDone int
// 	for numDone < numGoRoutines {
// 		<-doneCh
// 		numDone++
// 		fmt.Printf("finish go routine # %d\n", numDone)
// 	}
// 	fmt.Println("finished!!")
// }

// func writeString(start, end int, f *os.File, doneCh chan struct{}) {
// 	for i := start; i < end; i++ {
// 		_, err := f.WriteString(fmt.Sprintf("%06x\n", i))
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// 	doneCh <- struct{}{}
// }

// third approach using go routines properly.
func main() {
	f, err := os.Create("numbers3.txt")
	if err != nil {
		log.Fatalf("err: %s", err.Error())
	}
	defer f.Close()

	writerCh := make(chan string)
	doneWriterCh := make(chan struct{})
	go func() {
		for s := range writerCh {
			f.WriteString(s)
		}
		doneWriterCh <- struct{}{}
	}()

	numGoRoutines := 10
	doneCh := make(chan struct{})
	for i := 0; i <= end; i = i + (end / numGoRoutines) + 1 {
		next := i + (end / numGoRoutines)
		if next > end {
			next = end
		}
		go stringBuilder(i, next, writerCh, doneCh)
	}

	var doneNum int
	for doneNum < numGoRoutines {
		<-doneCh
		doneNum++
		fmt.Printf("finish %d go routine\n", doneNum)
	}
	close(writerCh)
	<-doneWriterCh
	fmt.Println("finished!!")
}

func stringBuilder(start, end int, result chan string, done chan struct{}) {
	var builder strings.Builder
	for i := start; i <= end; i++ {
		fmt.Fprintf(&builder, "%06x\n", i)
	}
	result <- builder.String()
	done <- struct{}{}
}
