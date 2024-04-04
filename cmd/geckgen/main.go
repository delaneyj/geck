package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/delaneyj/geck/generator"
	geckpb "github.com/delaneyj/geck/pb/gen/geck/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	ctx := context.Background()
	if err := run(ctx); err != nil {
		panic(err)
	}
}

func run(ctx context.Context) error {
	opts := &geckpb.GeneratorOptions{}
	b, err := os.ReadFile("./geckgen.json")
	if err != nil {
		return fmt.Errorf("failed to read bundle definition: %w", err)
	}

	jsonOpts := protojson.UnmarshalOptions{DiscardUnknown: true}
	if err = jsonOpts.Unmarshal(b, opts); err != nil {
		return fmt.Errorf("failed to unmarshal bundle definition: %w", err)
	}
	return generator.BuildECS(ctx, opts)
}
