package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fikryfahrezy/forward/csv-processing/bankstatement"
	"github.com/fikryfahrezy/forward/csv-processing/caster"
	"github.com/fikryfahrezy/forward/csv-processing/worker"
)

type csvResult struct {
	transactions []bankstatement.Transaction
	errMessages  map[string]any
}

var usage = func() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "\nA simple Bank Statement CSV processing with worker pool.\n\n")
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
		os.Exit(1)
	}

	wrk := worker.New(*wokerPool, worker.WithLogger(os.Stdout))
	csvPaths := strings.Split(*csvPath, ",")

	results := make([]<-chan csvResult, 0)
	for _, csvPath := range csvPaths {
		result, err := wrk.Add(func() any {
			f, err := os.Open(strings.TrimSpace(csvPath))
			if err != nil {
				fmt.Printf("Failed to open CSV: %s, Error: %s", csvPath, err)
				return csvResult{}
			}

			defer func() {
				if err := f.Close(); err != nil {
					fmt.Printf("Failed to close CSV: %s, Error: %s", csvPath, err)
				}
			}()

			transactions, errMessages := bankstatement.ParseCSV(f)
			return csvResult{
				transactions: transactions,
				errMessages:  errMessages,
			}
		})
		if err != nil {
			fmt.Printf("Failed to process CSV: %s, Error: %s", csvPath, err)
			continue
		}
		results = append(results, caster.ChanType[csvResult](result))
	}

	for i, resultChan := range results {
		result := <-resultChan
		fmt.Printf("CSV %s, transactions: %v, error: %v \n", csvPaths[i], result.transactions, result.errMessages)
	}

	wrk.Close()
	fmt.Println("All jobs completed and results collected.")
}
