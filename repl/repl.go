package repl

import (
	"bufio"
	"fmt"
	"io"

	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

const prompt = ">> "
const elephant = `
		( ͡° ͜ʖ ͡°)
`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		_, err := fmt.Fprintf(out, prompt)
		if err != nil {
			return
		}
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			_, err = io.WriteString(out, evaluated.Inspect())
			if err != nil {
				return
			}
			_, err = io.WriteString(out, "\n")
			if err != nil {
				return
			}
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	_, err := io.WriteString(out, elephant+"\n")
	if err != nil {
		return
	}
	_, err = io.WriteString(out, "parser errors:\n")
	if err != nil {
		return
	}
	for _, msg := range errors {
		_, err := io.WriteString(out, "\t"+msg+"\n")
		if err != nil {
			return
		}
	}
}
