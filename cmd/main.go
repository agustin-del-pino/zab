package main

import (
	"bytes"
	"fmt"
	"github.com/agustin-del-pino/zab/pkg/zab"
	"go/format"
	"os"
	"path"
	"strings"
)

func main() {
	args := os.Args[1:]
	dir := path.Dir(args[0])
	fn := strings.ReplaceAll(path.Base(args[0]), ".go", "blr.go")
	n, _, err := zab.ReadFile(args[0])

	if err != nil {
		panic(err)
	}

	s := zab.GetStructs(n)
	bf := bytes.NewBufferString("package ")

	bf.WriteString(n.Name.Name)
	bf.WriteString("\n\n")

	for _, strc := range s {
		_, wErr := bf.Write(zab.CreateBuilder(strc))
		if wErr != nil {
			panic(wErr)
		}
	}

	fms, fmsErr := format.Source(bf.Bytes())

	if fmsErr != nil {
		fmt.Println(string(bf.Bytes()))
		panic(fmsErr)
	}

	wErr := os.WriteFile(fmt.Sprintf("%s/%s", dir, fn), fms, 0777)
	if wErr != nil {
		panic(wErr)
	}
}
