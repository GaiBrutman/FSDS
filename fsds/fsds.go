// fsds (File Size Distribution System) calculates the size distribution of all directories in a given path
package fsds

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

// A struct type representing the result of a GetDirSize() goroutine.
type dirResult struct {
	Path     string  // File path of the given directory/file.
	Size     int64   //Size of the given directory/file.
	Duration float64 //Time (in seconds) took to calculate the size.
}

// GetDirSizeMap calculates the size of each file/directory in a given directory.
// Returns a map[string]int64 of relative directory/file path to directory/file size (in bytes).
func GetDirSizeMap(dir string) map[string]int64 {
	start := time.Now()

	m := make(map[string]int64)

	f, err := os.Stat(dir)

	if err != nil {
		log.Fatal(err)
	} else if !f.IsDir() {
		fmt.Printf("%q: size = %d\n", f.Name(), f.Size())
	}

	fmt.Printf("Distribution of %q:\n", dir)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan dirResult)

	for _, f := range files {
		go GetDirSize(filepath.Join(dir, f.Name()), ch)
	}

	var sum int64
	for range files {
		result := <-ch
		m[result.Path] = result.Size
		sum += result.Size

		fmt.Printf("\t%.2fs\t%.3fMB\t\t%q\n", result.Duration, float64(result.Size)*9.5367e-7, result.Path)
	}

	fmt.Printf("%.2fs elapsed\t%.3fMB Total\n", time.Since(start).Seconds(), float64(sum)*9.5367e-7)
	return m
}

// GetDirSizeMap calculates recursively the size of a given directory/file.
// Returns the size of the given directory/file and writes to a dirResult channel if given (not nil).
func GetDirSize(dir string, ch chan<- dirResult) int64 {
	var sum int64
	start := time.Now()

	f, err := os.Stat(dir)

	if err != nil {
		log.Println(err)

		if ch != nil {
			ch <- dirResult{dir, 0, time.Since(start).Seconds()}
		}
		return 0
	} else if !f.IsDir() {
		if ch != nil {
			ch <- dirResult{dir, int64(f.Size()), time.Since(start).Seconds()}
		}
		return int64(f.Size())
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)

		if ch != nil {
			ch <- dirResult{dir, 0, time.Since(start).Seconds()}
		}
		return 0
	}

	for _, f := range files {
		sum += GetDirSize(filepath.Join(dir, f.Name()), nil)
	}

	if ch != nil {
		ch <- dirResult{dir, sum, time.Since(start).Seconds()}
	}
	return sum
}
