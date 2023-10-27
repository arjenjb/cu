package terminal

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_overflow(t *testing.T) {
	s := NewScreen(Point{5, 5}, nil)
	s.Write([]byte("123456\nAap\nBanaan"))

	fmt.Print(s.Buffer())
}

func Test_writeText_1(t *testing.T) {
	var runs []Run
	runs = writeText(runs, "Hello Bert", 0, Style{})
	runs = writeText(runs, "World!", 6, Style{})
	fmt.Print((&Line{
		width: 80,
		runs:  runs,
	}).String())
}

func Test_writeText_2(t *testing.T) {
	var runs []Run
	runs = writeText(runs, "Hello Bert", 0, Style{})
	runs = writeText(runs, "u", 7, Style{})

	assert.Equal(t, "Hello Burt", (&Line{runs: runs}).String())
}

func Test_writeText_3(t *testing.T) {
	var runs []Run
	runs = writeText(runs, "Hello Bert", 0, Style{})
	runs = writeText(runs, "u", 7, Style{})
	runs = writeText(runs, "Henk!", 6, Style{})

	assert.Equal(t, "Hello Henk!", (&Line{runs: runs}).String())
}

func Test_writeText_4(t *testing.T) {
	var runs []Run
	runs = writeText(runs, "0", 0, Style{})
	runs = writeText(runs, "1", 1, Style{})
	runs = writeText(runs, "2", 2, Style{})
	runs = writeText(runs, "3", 3, Style{})
	runs = writeText(runs, "4", 4, Style{})
	runs = writeText(runs, "abc", 0, Style{})

	assert.Equal(t, "abc34", (&Line{runs: runs}).String())
}

func TestNewScreen(t *testing.T) {
	s := NewScreen(Point{20, 40}, nil)
	s.WriteString("Hello world\n")
	s.CursorUp(1)
	s.CursorRight(6)
	s.WriteString("Bert")

	fmt.Printf("%s\n", s.Buffer())
}

func TestMovementCenter(t *testing.T) {
	input := "Positioning the cursor the end of the line:\n\033[5;5H1. At the center\u001B[2;1H2. At the second line"

	s := NewScreen(Point{20, 40}, nil)
	s.WriteString(input)

	fmt.Printf("%s\n", s.Buffer())
}

//func TestProgressBar(t *testing.T) {
//	s := NewScreen(Point{20, 40}, nil)
//
//	bar := progressbar.NewOptions(1000,
//		progressbar.OptionSetWriter(s),
//		progressbar.OptionEnableColorCodes(true),
//		progressbar.OptionShowBytes(true),
//		progressbar.OptionSetWidth(15),
//		progressbar.OptionSetDescription("[cyan][1/3][reset] Writing moshable file..."),
//		progressbar.OptionSetTheme(progressbar.Theme{
//			Saucer:        "[green]=[reset]",
//			SaucerHead:    "[green]>[reset]",
//			SaucerPadding: " ",
//			BarStart:      "[",
//			BarEnd:        "]",
//		}))
//
//	for i := 0; i < 100; i++ {
//		bar.Add(1)
//		time.Sleep(5 * time.Millisecond)
//	}
//
//	fmt.Printf("%s\n", s.Buffer())
//}
