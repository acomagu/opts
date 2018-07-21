package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/build"
	"go/format"
	"io"
	"os"
	"strings"
	"text/scanner"

	"github.com/pkg/errors"
)

type importValue []string

func (v *importValue) String() string {
	return "import paths list"
}

func (v *importValue) Set(e string) error {
	*v = append(*v, e)
	return nil
}

var (
	typ    = flag.String("type", "", "type name; must be set")
	opt    = flag.String("name", "", "option name; if omitted, used type name")
	output = flag.String("output", "opts.go", "output file name; default srcdir/opts.go")
	apnd   = flag.Bool("append", false, "append to the output file.")
	imp = &importValue{}
)

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: opts <flags> [<directory>]")
	flag.PrintDefaults()
}

func main() {
	dir := "."

	flag.Var(imp, "import", "import path; can be specified multiple times")
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if len(args) >= 1 {
		dir = args[0]
	}
	if len(*typ) == 0 {
		fmt.Fprintln(os.Stderr, "invalid type name")
		os.Exit(1)
	}
	if len(*opt) == 0 {
		sep := strings.Split(*typ, ".")
		*opt = sep[len(sep)-1]
		if (*opt)[0] == '*' {
			*opt = (*opt)[1:]
		}
	}

	if err := run(dir, *apnd, *imp, *typ, *opt, *output); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(dir string, apnd bool, imports []string, typeName, optName, output string) error {
	pkgName, err := getPkgName(".", dir)
	if err != nil {
		return errors.Wrap(err, "could not get the dir's package info")
	}

	f, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return errors.Wrap(err, "could not open file")
	}

	g := &Generator{
		f:       f,
		pkgName: pkgName,
		imports: imports,
		typName: typeName,
		optName: optName,
		append:  apnd,
	}

	if err := g.generate(); err != nil {
		return err
	}

	return nil
}

func getPkgName(path, dir string) (string, error) {
	pkg, err := build.Import(path, dir, 0)
	if err != nil {
		return "", err
	}

	return pkg.Name, nil
}

type Generator struct {
	f       *os.File
	pkgName string
	imports []string
	typName string
	optName string
	append  bool
}

func (g *Generator) generate() error {
	buf := bytes.NewBuffer(nil)

	if g.append {
		if _, err := io.Copy(buf, g.f); err != nil {
			return err
		}

		if _, err := g.f.Seek(0, 0); err != nil {
			return err
		}
	} else {
		g.f.Truncate(0)

		if err := writeHeadTmpl(buf, g.pkgName); err != nil {
			return err
		}
	}

	buf = bytes.NewBufferString(g.addImports(buf.String(), g.imports))

	if err := writeTmpl(buf, g.typName, g.optName); err != nil {
		return err
	}

	bs, err := format.Source(buf.Bytes())
	if err != nil {
		return errors.Wrapf(err, "could not format the result source:\n%s", buf.String())
	}

	if _, err := g.f.Write(bs); err != nil {
		return errors.Wrap(err, "could not write to file")
	}

	return nil
}

func (g *Generator) addImports(orig string, paths []string) string {
	r := bytes.NewReader([]byte(orig))
	s := new(scanner.Scanner).Init(r)
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		if s.TokenText() == "package" {
			s.Scan()
			break
		}
	}

	buf := bytes.NewBufferString(orig[0:s.Pos().Offset])
	for _, path := range paths {
		fmt.Fprintf(buf, "\nimport \"%s\"\n", path)
	}

	buf.WriteString(orig[s.Pos().Offset:len(orig)])

	return buf.String()
}
