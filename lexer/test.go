// This file is just for testing purposes.

package lexer

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"quoi/token"
	"strings"
)

const directory = "tests"

// log error, and return an empty string in case of an error.
// this function is not that important, so I am not doing this right.
func readFile(name string) string {
	name = directory + "/" + name
	bx, err := os.ReadFile(name)
	if err != nil {
		fmt.Println("readFile: error: ", err.Error())
		return ""
	}
	return string(bx)
}

func writeFile(name, data string) {
	name = directory + "/" + name
	if err := os.WriteFile(name, []byte(data), 0777); err != nil {
		panic("writeFile: error: " + err.Error())
	}
}

// compare the contents wanted.tokens file to result.tokens file.
// also create a diff.tokens file
//
// these files ending with 'tokens' extensions just contain lines of token represantations
// separated by newlines.
//token represantations are generated by tokRepr function in this file.
func compareResults(want, result, sourceFile string) bool {
	type diff struct {
		wantLine, resultLine string
		lineNumber           uint
	}
	var reprDiff = func(d diff) string {
		return fmt.Sprintf("======\nAT LINE %d\nWANTED: %sGOT: %s\n======", d.lineNumber, d.wantLine, d.resultLine)
	}

	var (
		buf1  bytes.Buffer
		buf2  bytes.Buffer
		diffs = []diff{}
	)
	buf1.WriteString(want)
	buf2.WriteString(result)
	i := uint(0)
	for {
		line1, err := buf1.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic("compareResults: buf1 read string: " + err.Error())
		}
		line2, err := buf2.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic("compareResults: buf2 read string: " + err.Error())
		}
		if line1 != line2 {
			diffs = append(diffs, diff{wantLine: line1, resultLine: line2, lineNumber: i})
		}
		i++
	}
	var eof1, eof2 bool
	if _, err := buf1.ReadByte(); err != nil {
		eof1 = true
	}
	if _, err := buf2.ReadByte(); err != nil {
		eof2 = true
	}
	bothReachedEOF := eof1 && eof2
	if !(bothReachedEOF) {
		log.Println("[!!!] Different number of line counts (it may be because of the newline at the end)")
		return false
	}
	var diffStr strings.Builder
	for _, v := range diffs {
		diffStr.WriteString(reprDiff(v))
	}
	diffFileName := "diff." + sourceFile + ".tokens"
	writeFile(diffFileName, diffStr.String())
	return len(diffs) == 0
}

func tokRepr(t token.Token) string {
	return fmt.Sprintf("(%s, %s)\n", t.Type, t.Literal)
}

// ignore whitespace tokens
func runTest(sourceFile string) {
	l := New(readFile(sourceFile))
	var res strings.Builder
	for {
		t := l.Next()
		if t.Type == token.EOF {
			break
		}
		repr := tokRepr(t)
		res.WriteString(repr)
	}
	resultsFileName := "result." + sourceFile + ".tokens"
	writeFile(resultsFileName, res.String())
	ok := compareResults(readFile("wanted.tokens"), readFile(resultsFileName), sourceFile)
	if ok {
		log.Println("[✊] PASS :)")
	}
}
