package cmd

import (
	"bufio"
	"fmt"
	"io"

	"github.com/w-h-a/interpreter/internal/lexer"
	"github.com/w-h-a/interpreter/internal/token"
)

const PROMPT = ">> "

func StartRepl(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)

		if !scanner.Scan() {
			return scanner.Err()
		}

		line := scanner.Text()

		tks := lexer.Lex(line)

		for tk := range tks {
			if tk.Type == token.EOF {
				break
			}
			fmt.Fprintf(out, "%+v\n", tk)
		}
	}
}
