package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/styling"
	"github.com/gotd/td/telegram/uploader"
	"github.com/sgorblex/tgsend/tdwrap"
)

func main() {
	file := flag.String("file", "", "file to be sent (optional)")
	init := flag.Bool("init", false, "initialize client. Overrides everything else")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s: easily send Telegram messages or files.\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "USAGE: %s [OPTIONS] RECIPIENT TEXT\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "OPTIONS:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "When sending a file, TEXT is its caption and can be omitted.\n")
	}
	flag.Parse()
	args := flag.Args()

	// client init mode
	if *init {
		err := tdwrap.SaveClientCreds()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	sendFile := *file != ""
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Insert recipient as first argument")
		os.Exit(1)
	}
	if !sendFile && len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Insert message text as second argument")
		os.Exit(1)
	}
	recipient := args[0]
	text := strings.Join(args[1:], " ")

	flow := tdwrap.GetFlow()
	client, err := tdwrap.GetClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	job := func(ctx context.Context) error {
		// auth if needed
		if err := client.Auth().IfNecessary(ctx, flow); err != nil {
			return fmt.Errorf("auth: %w", err)
		}

		api := client.API()
		sender := message.NewSender(api)

		// send text (no file)
		if !sendFile {
			_, err := sender.Resolve(recipient).Text(ctx, text)
			if err != nil {
				return fmt.Errorf("send text: %w", err)
			}
			return nil
		}

		// upload file
		u := uploader.NewUploader(api)
		sender.WithUploader(u)
		upload, err := u.FromPath(ctx, *file)
		if err != nil {
			return fmt.Errorf("upload %q: %w", *file, err)
		}

		// send uploaded file
		document := message.UploadedDocument(upload, styling.Plain(text))
		document.Filename(filepath.Base(*file))
		target := sender.Resolve(recipient)
		if _, err := target.Media(ctx, document); err != nil {
			return fmt.Errorf("send: %w", err)
		}

		return nil
	}

	err = client.Run(ctx, job)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
