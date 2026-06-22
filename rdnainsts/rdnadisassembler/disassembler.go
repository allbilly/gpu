package main

import (
	"debug/elf"
	"fmt"
	"os"
	"path/filepath"

	"gitlab.com/akita/navisim/rdnainsts"
)

func main() {
	path := os.Args[1]
	elfFile, err := elf.Open(path)
	if err != nil {
		_ = fmt.Errorf("failed to open file %v", path)
	}
	defer func() {
		err := elfFile.Close()
		if err != nil {
			panic(err)
		}
	}()

	_, filename := filepath.Split(path)

	disasm := rdnainsts.NewDisassembler()

	disasm.Disassemble(elfFile, filename, os.Stdout)
}
