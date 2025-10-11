package cmd

import (
	"bufio"
	"fmt"
	"io"

	"github.com/w-h-a/interpreter/internal/lexer"
	"github.com/w-h-a/interpreter/internal/parser"
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
		p := parser.New(tks)

		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			if err := printParserErrors(out, p.Errors()); err != nil {
				return err
			}
			continue
		}

		if _, err := io.WriteString(out, program.String()); err != nil {
			return err
		}

		if _, err := io.WriteString(out, "\n"); err != nil {
			return err
		}
	}
}

func printParserErrors(out io.Writer, errors []string) error {
	if _, err := io.WriteString(out, "Whoops! We ran into some monkey business here!\n"); err != nil {
		return err
	}

	if _, err := io.WriteString(out, " parser errors:\n"); err != nil {
		return err
	}

	for _, msg := range errors {
		if _, err := fmt.Fprintf(out, "\t%s\n", msg); err != nil {
			return err
		}
	}

	return nil
}
