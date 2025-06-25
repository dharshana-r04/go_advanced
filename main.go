package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "sync"
)
func ProcessLogs(inputFiles []string, outputFile string) error {
    var wg sync.WaitGroup
    channel := make(chan string, 100)

    for _, file := range inputFiles {
        wg.Add(1)
        go func(file string) {
            defer wg.Done()
            read(file, channel)
        }(file)
    }

    go func() {
        wg.Wait()
        close(channel)
    }()

    return write(outputFile, channel)
}
func read(filename string, channel chan<- string) {
    file, err := os.Open(filename)
    if err != nil {
        fmt.Printf("Failed to open file %s: %v\n", filename, err)
        return
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, "ERROR") {
            channel <- line
        }
    }
    if err := scanner.Err(); err != nil {
        fmt.Printf("Error reading file %s: %v\n", filename, err)
    }
}

func write(outputFile string, channel <-chan string) error {
    file, err := os.Create(outputFile)
    if err != nil {
        return fmt.Errorf("failed to create output file: %v", err)
    }
    defer file.Close()

    writer := bufio.NewWriter(file)
    defer writer.Flush()

    for line := range channel {
        _, err := writer.WriteString(line + "\n")
        if err != nil {
            return fmt.Errorf("failed to write to output file: %v", err)
        }
    }

    return nil
}

func main() {
    inputFiles := []string{"server1.log", "server2.log", "server3.log"}
    err := ProcessLogs(inputFiles, "errors.log")
    if err != nil {
        fmt.Println("Error processing logs:", err)
    } else {
        fmt.Println("Errors written to errors.log file")
    }
}
