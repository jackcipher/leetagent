package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"leetagent/internal/webapp"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "serve":
		if err := serve(os.Args[2:]); err != nil {
			log.Fatal(err)
		}
	case "help", "-h", "--help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", os.Args[1])
		usage()
		os.Exit(2)
	}
}

func serve(args []string) error {
	flags := flag.NewFlagSet("serve", flag.ContinueOnError)
	addr := flags.String("addr", "127.0.0.1:8080", "HTTP listen address")
	root := flags.String("root", ".", "leetagent repository root")
	if err := flags.Parse(args); err != nil {
		return err
	}

	absRoot, err := filepath.Abs(*root)
	if err != nil {
		return fmt.Errorf("resolve repository root: %w", err)
	}
	if _, err := os.Stat(filepath.Join(absRoot, "go.mod")); err != nil {
		return fmt.Errorf("%s does not look like the leetagent repository: %w", absRoot, err)
	}

	app, err := webapp.New(absRoot)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:              *addr,
		Handler:           app,
		ReadHeaderTimeout: webapp.ReadHeaderTimeout,
		IdleTimeout:       webapp.IdleTimeout,
	}
	log.Printf("leetagent web is available at http://%s", *addr)
	log.Printf("repository: %s", absRoot)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  leetagent serve [-addr 127.0.0.1:8080] [-root .]")
}
