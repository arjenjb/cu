package terminal

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_LineSplit(t *testing.T) {
	l := VirtualLine{}
	l.Write("123451234512345_", 5, Style{})

	lines := l.Split(5)
	assert.Len(t, lines, 5)
}

func Test_Resize(t *testing.T) {
	s := NewScreen(Point{5, 5}, nil)
	s.Write([]byte("123456\nAap\nBanaan"))
	assert.Equal(t, "12345\n6\nAap\nBanaa\nn\n", s.Buffer())

	s.updateWidth(10)
	assert.Equal(t, "123456\nAap\nBanaan\n", s.Buffer())

	s.updateWidth(5)
	assert.Equal(t, "12345\n6\nAap\nBanaa\nn\n", s.Buffer())
}

func Test_overflow(t *testing.T) {
	s := NewScreen(Point{5, 5}, nil)
	s.Write([]byte("123456\nAap\nBanaan"))

	fmt.Println("--- Small")
	fmt.Print(s.Buffer())
	fmt.Println("")

	fmt.Println("--- Big")
	s.updateWidth(10)

	fmt.Print(s.Buffer())
	fmt.Println("")

	s.updateWidth(5)

	fmt.Print(s.Buffer())

	s.updateWidth(10)

	fmt.Print(s.Buffer())

}

func Test_writeText_1(t *testing.T) {
	var runs []Run
	runs = writeText(runs, "Hello Bert", 0, Style{})
	runs = writeText(runs, "World!", 6, Style{})
	fmt.Print((&Line{
		runs: runs,
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

func TestVirtualLines(t *testing.T) {
	s := NewScreen(Point{2, 10}, nil)
	s.WriteString("abc")
	s.WriteNewLine()
	s.WriteString("1234567")

	assert.Len(t, s.Lines(), 6)

	l := s.virtualLines()
	assert.Len(t, l, 2)
}

func TestResizeGreaterWidth(t *testing.T) {
	s := NewScreen(Point{2, 10}, nil)
	s.WriteString("abc")
	s.WriteNewLine()
	s.WriteString("1234567")

	assert.Len(t, s.Lines(), 6)

	s.updateWidth(10)
	assert.Len(t, s.Lines(), 2)
}

func TestResizeHeight(t *testing.T) {
	s := NewScreen(Point{3, 3}, nil)
	s.WriteString("abcdef")
	s.WriteNewLine()
	s.WriteString("b")
	s.WriteNewLine()
	s.WriteString("c")
	s.WriteNewLine()
	s.WriteString("d")

	require.Len(t, s.Lines(), 5)
	require.Equal(t, s.scrollTop, 2)

	s.scrollTop = 1
	s.updateHeight(2)

	require.Equal(t, s.scrollTop, 2)
}

func TestResizeBehaviorHeight(t *testing.T) {
	s := NewScreen(Point{1, 1}, nil)
	s.WriteString("a")
	s.WriteNewLine()
	s.WriteString("b")

	require.Len(t, s.Lines(), 2)
	require.Equal(t, 1, s.scrollTop)

	s.updateHeight(2)
	require.Equal(t, 0, s.scrollTop)

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
