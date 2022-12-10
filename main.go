package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/neovim/go-client/nvim"
)

func main() {
	if len(os.Args) < 3 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		usage()
		os.Exit(0)
	}

	// Get address from environment variable set by Nvim.
	addr := os.Getenv("NVIM")
	if addr == "" {
		log.Fatal("NVIM not set")
	}

	// Dial with default options.
	v, err := nvim.Dial(addr)
	if err != nil {
		log.Fatal(err)
	}

	// Cleanup on return.
	defer v.Close()

	switch os.Args[1] {
	case "-c":
		command(v)
	case "-cwait":
		command(v)
		wait_current_buffer(v)
	case "-e":
		eval(v)
	default:
		fmt.Printf("Unkown option '%s'\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Printf(`
Usage: %s [opts] arguments

Opts:
	-h, --help
		Show this message
	-c ...
		Run a command
	-cwait ...
		Run a command and wait for the 'bdelete'
	-e ...
		Evaluate an expression
`, os.Args[0])
}

func command(v *nvim.Nvim) {
	cmd := strings.Join(os.Args[2:], " ")
	if err := v.Command(cmd); err != nil {
		log.Fatal(err)
	}
}

func eval(v *nvim.Nvim) {
	expr := strings.Join(os.Args[2:], " ")
	var result string
	if err := v.Eval(expr, &result); err != nil {
		log.Fatal(err)
	}
	fmt.Print(result)
}

func wait_current_buffer(v *nvim.Nvim) {
	chainID := v.ChannelID()
	b := v.NewBatch()
	b.Command("augroup term-util")
	b.Command(fmt.Sprintf("autocmd BufDelete <buffer> silent! call rpcnotify(%d, 'BufDelete')", chainID))
	b.Command(fmt.Sprintf("autocmd VimLeave * if exists('v:exiting') && v:exiting > 0 | silent! call rpcnotify(%d, 'Exit') | endif", chainID))
	b.Command("augroup END")
	if err := b.Execute(); err != nil {
		log.Fatal(err)
	}

	finishCh := make(chan struct{}, 1)
	finish := func() {
		select {
		case finishCh <- struct{}{}:
		default:
		}
	}
	if err := v.RegisterHandler("BufDelete", finish); err != nil {
		log.Fatal(err)
	}
	if err := v.RegisterHandler("Exit", finish); err != nil {
		log.Fatal(err)
	}

	<-finishCh
}
