package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"
	"text/template"
)

const (
	generatedFileMode = 0644
	generatedFileFlag = os.O_WRONLY | os.O_TRUNC | os.O_CREATE
)

const templateMain = `// Code generated; DO NOT EDIT.

package {{ .GoPackage }}

func getInputMonkeys() []Monkey {
	return []Monkey{
		{{- range .Monkeys }}
		{
			Items: []int{{ "{" }}{{ .Items }}{{ "}" }},
			Op: func(old int) int {
				return {{ .Op }}
			},
			Test: {{ .Test }},
			JT: {{ .JT }},
			JF: {{ .JF }},
		},
		{{- end }}
	}
}
`

type (
	Monkey struct {
		Items string
		Op    string
		Test  string
		JT    string
		JF    string
	}

	TplData struct {
		GoPackage string
		Monkeys   []Monkey
	}
)

func main() {
	gopackage := os.Getenv("GOPACKAGE")
	if len(gopackage) == 0 {
		log.Fatalln("No GOPACKAGE env var provided by go generate")
	}

	var outfile string
	flag.StringVar(&outfile, "o", "", "output file")
	flag.Parse()

	if len(outfile) == 0 {
		log.Fatalln("No output file (-o) provided")
	}

	tplmain, err := template.New("main").Parse(templateMain)
	if err != nil {
		log.Fatalln(err)
	}

	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	var monkeys []Monkey

	var monkey Monkey

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "Monkey"):
			monkey = Monkey{}
		case line == "":
			monkeys = append(monkeys, monkey)
		case strings.HasPrefix(line, "  Starting items: "):
			monkey.Items = strings.TrimPrefix(line, "  Starting items: ")
		case strings.HasPrefix(line, "  Operation: "):
			monkey.Op = strings.Replace(strings.TrimPrefix(line, "  Operation: "), "new = ", "", 1)
		case strings.HasPrefix(line, "  Test: divisible by "):
			monkey.Test = strings.TrimPrefix(line, "  Test: divisible by ")
		case strings.HasPrefix(line, "    If true: throw to monkey "):
			monkey.JT = strings.TrimPrefix(line, "    If true: throw to monkey ")
		case strings.HasPrefix(line, "    If false: throw to monkey "):
			monkey.JF = strings.TrimPrefix(line, "    If false: throw to monkey ")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	if monkey != (Monkey{}) {
		monkeys = append(monkeys, monkey)
		monkey = Monkey{}
	}

	genfile, err := os.OpenFile(outfile, generatedFileFlag, generatedFileMode)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := genfile.Close(); err != nil {
			log.Println("Failed closing file", err)
		}
	}()
	if err := tplmain.Execute(genfile, TplData{
		GoPackage: gopackage,
		Monkeys:   monkeys,
	}); err != nil {
		log.Fatalln(err)
	}
}
