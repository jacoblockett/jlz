package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	jlz "github.com/jacoblockett/jlz/engine"
)

func main() {
	if len(os.Args) < 3 || len(os.Args) > 4 {
		usage()
	}
	for _, name := range os.Args[2:] {
		if strings.HasPrefix(name, "-") {
			usage() // no flags or optional algorithms
		}
	}
	action, path := os.Args[1], os.Args[2]
	if action != "compress" && action != "decompress" {
		usage()
	}
	info, err := os.Stat(path)
	if err != nil {
		die(err)
	}
	if info.Size() > jlz.MaxInputSize {
		die(fmt.Errorf("input exceeds %d-byte safety limit", jlz.MaxInputSize))
	}
	input, err := os.ReadFile(path)
	if err != nil {
		die(err)
	}
	start := time.Now()
	var output []byte
	var report jlz.Report
	var target string
	if action == "compress" {
		output, report, err = jlz.Compress(input)
		target = path + ".jlz"
	} else {
		output, err = jlz.Decompress(input, jlz.MaxInputSize)
		target = strings.TrimSuffix(path, ".jlz") + ".restored"
	}
	if err != nil {
		die(err)
	}
	if len(os.Args) == 4 {
		target = os.Args[3]
	}
	if err = save(target, output, path); err != nil {
		die(err)
	}
	if action == "compress" {
		report.Print(target, time.Since(start))
	} else {
		fmt.Printf("Restored %d bytes -> %s (SHA-256 verified, %s)\n", len(output), target, time.Since(start).Round(time.Millisecond))
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: jlz <compress|decompress> <input> [output]")
	fmt.Fprintln(os.Stderr, "JLZ has one automatic compression algorithm. No levels or experimental switches.")
	os.Exit(2)
}
func die(err error) { fmt.Fprintln(os.Stderr, "Error:", err); os.Exit(1) }
func save(path string, b []byte, source string) error {
	a, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	s, err := filepath.Abs(source)
	if err != nil {
		return err
	}
	if a == s {
		return errors.New("refusing to overwrite input file")
	}
	return os.WriteFile(path, b, 0644)
}
