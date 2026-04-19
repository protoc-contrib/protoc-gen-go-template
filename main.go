package main // import "github.com/protoc-contrib/protoc-gen-go-template"

import (
	"io"
	"log"
	"os"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"

	pgghelpers "github.com/protoc-contrib/protoc-gen-go-template/helpers"
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("reading input: %v", err)
	}

	req := &pluginpb.CodeGeneratorRequest{}
	if err := proto.Unmarshal(data, req); err != nil {
		log.Fatalf("parsing input proto: %v", err)
	}

	if len(req.FileToGenerate) == 0 {
		log.Fatal("no files to generate")
	}

	resp := &pluginpb.CodeGeneratorResponse{}
	if err := pgghelpers.ParseParams(req, resp); err != nil {
		msg := err.Error()
		resp.Error = &msg
	}

	out, err := proto.Marshal(resp)
	if err != nil {
		log.Fatalf("failed to marshal output proto: %v", err)
	}

	if _, err := os.Stdout.Write(out); err != nil {
		log.Fatalf("failed to write output proto: %v", err)
	}
}
