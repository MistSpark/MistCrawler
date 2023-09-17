package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {
	// Define command-line flags for the input file, output file, number of threads, silent mode, and help
	inputFile := flag.String("u", "", "Input file containing the list of URLs")
	outputFile := flag.String("o", "", "Output file to save valid URLs")
	threads := flag.Int("t", 1, "Number of concurrent threads for checking URLs")
	silent := flag.Bool("silent", false, "Run in silent mode (suppress stdout)")
	help := flag.Bool("h", false, "Display help message")
	flag.Parse()

	// Display help message and exit if the -h flag is provided
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// Check if the input file flag is provided
	if *inputFile == "" {
		fmt.Println("Usage: MistCrawler -u <input_file> -o <output_file> -t <threads> [-silent]")
		os.Exit(1)
	}

	// Check if the output file flag is provided
	if *outputFile == "" {
		fmt.Println("You must specify an output file using the -o flag.")
		os.Exit(1)
	}

	// Initialize the logger
	var logger *log.Logger
	if *silent {
		// Create a logger with a writer that discards stdout (no output)
		logger = log.New(io.Discard, "", 0)
	} else {
		// Create a logger that writes to stdout
		logger = log.New(os.Stdout, "", 0)
	}

	// Open the input file
	file, err := os.Open(*inputFile)
	if err != nil {
		logger.Printf("Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Create or clear the output file
	output, err := os.Create(*outputFile)
	if err != nil {
		logger.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer output.Close()

	// Scanner to read lines from the input file
	scanner := bufio.NewScanner(file)

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Create a channel for controlling the number of concurrent goroutines
	threadLimit := make(chan struct{}, *threads)

	// Loop through each URL in the input file
	for scanner.Scan() {
		url := scanner.Text()

		// Skip empty URLs
		if url == "" {
			continue
		}

		// Acquire a token from the channel to limit concurrent goroutines
		threadLimit <- struct{}{}

		// Increment the wait group counter
		wg.Add(1)

		// Goroutine for checking the URL
		go func(url string) {
			defer wg.Done()
			defer func() {
				// Release the token back to the channel when done
				<-threadLimit
			}()

			// Make an HTTP GET request
			resp, err := http.Get(url)
			if err != nil {
				logger.Printf("Error making HTTP request to %s: %v\n", url, err)
				return
			}
			defer resp.Body.Close()

			// Check if the response code is 200 (OK)
			if resp.StatusCode == http.StatusOK {
				// If it's 200, write the URL to the output file
				_, err := output.WriteString(url + "\n")
				if err != nil {
					logger.Printf("Error writing to output file: %v\n", err)
				}
				logger.Printf("URL %s is valid and added to %s\n", url, *outputFile)
			} else {
				logger.Printf("URL %s returned HTTP status code %d and is not valid.\n", url, resp.StatusCode)
			}
		}(url)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		logger.Printf("Scanner error: %v\n", err)
	}
}

