package main

import (
	"context"
	_ "embed"
	"io"
	"log"
	"os"

	"github.com/FlowingSPDG/streamdeck"

	sdnewtek "github.com/FlowingSPDG/streamdeck-newtek/Source/code"
)

func main() {
	logfile, err := os.Create("./streamdeck-newtek.log")
	if err != nil {
		panic("cannnot open log:" + err.Error())
	}
	defer logfile.Close()
	log.SetOutput(io.MultiWriter(logfile, os.Stdout))
	log.SetFlags(log.Ldate | log.Ltime)

	ctx := context.Background()
	log.Println("Starting...")
	if err := run(ctx); err != nil {
		log.Fatalf("%v\n", err)
	}
}

func run(ctx context.Context) error {
	params, err := streamdeck.ParseRegistrationParams(os.Args)
	if err != nil {
		return err
	}

	client := sdnewtek.NewSDNewTek(ctx, params)

	return client.Run()
}
