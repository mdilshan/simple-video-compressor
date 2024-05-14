package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func readOutput(scanner *bufio.Scanner) {
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func main() {
	inputFile := flag.String("i", "", "Input video file path")
	outputFile := flag.String("o", "", "Output video file path")
	// non required flags
	logLevel := flag.String("ll", "error", "Log levels: quiet, panic, fatal, error, warning, info, verbose, debug, trace")

	flag.Parse()

	if *inputFile == "" || *outputFile == "" {
		flag.PrintDefaults()
		return
	}

	cmd := exec.Command("ffmpeg", "-i", *inputFile, "-c:v", "libx264", "-crf", "28", "-preset", "medium", "-c:a", "aac", "-b:a", "96k", *outputFile, "-loglevel", *logLevel)

	// ffmpeg write the log to stderr only. The log level can be set using the -loglevel flag
	// https://lists.ffmpeg.org/pipermail/libav-user/2013-July/005251.html
	stderr, errPipeError := cmd.StderrPipe()

	if errPipeError != nil {
		fmt.Println("Error creating StderrPipe: ", errPipeError)
		return
	}

	startError := cmd.Start()

	if startError != nil {
		fmt.Println("Error starting command: ", startError)
		return
	}

	scanner := bufio.NewScanner(stderr)

	go readOutput(scanner)

	fmt.Println("\033[96m", "Starting conversion...", "\033[0m")
	fmt.Println("\033[93m", "This could take a while depending on the length of the video file...", "\033[0m")
	cmd.Wait()

	fmt.Println("Process exited")
}
