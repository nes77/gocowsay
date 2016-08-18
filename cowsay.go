package main

import "flag"
import "fmt"
import "strings"
import "bytes"
import "unicode/utf8"
import "os"

const noSource string = "\x00\x01\x02\x03\x04"

type cowSayOptions struct {
	source string
}

func main() {
	var opts cowSayOptions
	flag.StringVar(&opts.source,
		"from",
		noSource,
		"A text source file containing lines to randomly select from. Use -- for stdin.")

	var wrapLen int
	flag.IntVar(&wrapLen, "wrap", 55, "number of characters from string per line.")
	flag.Parse()

	if noSource != opts.source {

	} else if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(-1)
	} else {
		// Regular arg
		txt := flag.Args()[0]
		out, err := encowseString(&txt, wrapLen)
		if err != nil {
			panic(err.Error())
		} else {

			fmt.Println(out)
		}
	}
}

const cow string = `
                \   ^__^
                 \  (oo)\_______
                    (__)\       )\/\
                         ||----w |
                         ||     ||
`

func encowseString(text *string, wrapLen int) (string, error) {
	textBox, err := makeTextBox(text, wrapLen)
	if err != nil {
		return "", err
	}
	textOut := make([]string, 2, 2)
	textOut[0] = textBox
	textOut[1] = cow
	return strings.Join(textOut, "\n"), nil
}

func makeTextBox(text *string, wrapLen int) (string, error) {
	var lines []string
	wordLines, err := wrap(text, wrapLen)
	if err != nil {
		return strings.Join(wordLines, "\n"), err
	}

	var topLineBytes = make([]rune, wrapLen+2, wrapLen+2)
	var bottomLineBytes = make([]rune, wrapLen+2, wrapLen+2)

	for i := 0; i < wrapLen+2; i++ {
		if i == 0 {
			topLineBytes[i] = '/'
			bottomLineBytes[i] = '\\'
		} else {
			topLineBytes[i] = '\u00AF'
			bottomLineBytes[i] = '_'
		}
	}

	topLineBytes[wrapLen+1] = '\\'
	bottomLineBytes[wrapLen+1] = '/'

	lines = append(lines, string(topLineBytes))
	buf := new(bytes.Buffer)

	for _, line := range wordLines {
		buf.WriteRune('|')
		buf.WriteString(line)
		for i := 0; i < wrapLen-utf8.RuneCountInString(line); i++ {
			buf.WriteRune(' ')
		}
		buf.WriteRune('|')

		lines = append(lines, buf.String())
		buf.Reset()
	}

	lines = append(lines, string(bottomLineBytes))

	return strings.Join(lines, "\n"), nil
}

func wrap(text *string, wrapLen int) ([]string, error) {
	out := make([]string, 0, strings.Count(*text, " "))
	words := strings.Split(*text, " ")
	currentLine := ""

	for _, word := range words {
		wordLen := utf8.RuneCountInString(word)
		if wordLen > wrapLen {
			return out, fmt.Errorf("%s too long for line wrap of %d chars", word, wrapLen)
		} else if utf8.RuneCountInString(currentLine)+wordLen > wrapLen {
			// Start new line
			out = append(out, currentLine)
			currentLine = word
		} else {
			currentLine = currentLine + " " + word
		}
	}

	if currentLine != "" {
		out = append(out, currentLine)
	}

	return out, nil

}

func longestString(strs []string) int {
	longestSeen := int(0)
	lll := int(-1)
	for index, s := range strs {
		cur := utf8.RuneCountInString(s)
		if cur > lll {
			lll = cur
			longestSeen = index
		}
	}

	return longestSeen
}
