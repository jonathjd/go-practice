package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode/utf8"

	"github.com/spf13/cobra"
)

var (
	countBytes bool
	countLines bool
	countWords bool
	countChars bool
)

var rootCmd = &cobra.Command{
	Use:   "ccwc [file]",
	Short: "ccwc is a wc cli utility written in Go",
	Long:  `A fast and expressive wc command line utility written in Go`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runCodingChallengeWordCount,
}

func init() {
	rootCmd.Flags().BoolVarP(&countBytes, "bytes", "c", false, "Count bytes")
	rootCmd.Flags().BoolVarP(&countLines, "lines", "l", false, "Count lines")
	rootCmd.Flags().BoolVarP(&countWords, "words", "w", false, "Count words")
	rootCmd.Flags().BoolVarP(&countChars, "characters", "m", false, "Count characters")
}

func runCodingChallengeWordCount(cmd *cobra.Command, args []string) error {
	var (
		r      io.Reader
		name   string
		allOps bool
	)

	// stdin vs file
	if len(args) == 0 || args[0] == "-" {
		r = os.Stdin
		name = ""
	} else {
		f, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer f.Close()
		r = f
		name = args[0]
	}

	// default to lines, words, bytes if no flags given
	if !countBytes && !countLines && !countWords && !countChars {
		allOps = true
	}

	// one pass over input
	var (
		bytesCount int64
		linesCount int64
		wordsCount int64
		charsCount int64
	)

	scanner := bufio.NewScanner(r)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024) // allow larger lines if needed

	for scanner.Scan() {
		line := scanner.Text()

		// bytes: include newline if it was present
		bytesCount += int64(len(line)) + 1

		// lines
		linesCount++

		// words
		inWord := false
		for _, r := range line {
			if r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == '\f' || r == '\v' {
				if inWord {
					wordsCount++
					inWord = false
				}
			} else {
				inWord = true
			}
		}
		if inWord {
			wordsCount++
		}

		// chars (runes)
		charsCount += int64(utf8.RuneCountInString(line) + 1) // +1 for newline rune
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// If input came from a regular file and byte counting is requested or defaulted,
	// prefer Stat for exact bytes including final newline presence.
	if f, ok := r.(*os.File); ok && (countBytes || allOps) {
		if fi, err := f.Stat(); err == nil {
			bytesCount = fi.Size()
		}
	}

	// printing logic
	var out []interface{}
	if countLines || allOps {
		out = append(out, linesCount)
	}
	if countWords || allOps {
		out = append(out, wordsCount)
	}
	if countBytes || allOps {
		out = append(out, bytesCount)
	}
	if countChars {
		out = append(out, charsCount)
	}
	if name != "" {
		out = append(out, name)
	}

	fmt.Println(out...)
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
