package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fikryfahrezy/forward/csv-processing/caster"
	"github.com/fikryfahrezy/forward/csv-processing/worker"
)

var usage = func() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "\nA simple CSV Processing with worker pool.\n\n")
	fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
	flag.PrintDefaults()
}

func main() {
	wokerPool := flag.Int("pool", 3, "Number of worker pool")
	csvPath := flag.String("csv", "", "Comma separated CSV path")
	flag.Usage = usage
	flag.Parse()

	if *csvPath == "" {
		usage()
		return
	}

	wrk := worker.New(*wokerPool)
	csvPaths := strings.Split(*csvPath, ",")

	results := make([]<-chan string, 0)
	for _, csvPath := range csvPaths {
		result, err := wrk.Add(func() any {
			return csvPath
		})
		if err != nil {
			fmt.Printf("Failed to process CSV: %s, Error: %s", csvPath, err)
			continue
		}
		results = append(results, caster.ChanType[string](result))
	}

	for i, resultChan := range results {
		result := <-resultChan
		fmt.Printf("CSV %s result: %s\n", csvPaths[i], result)
	}

	wrk.Close()
	fmt.Println("All jobs completed and results collected.")
}
