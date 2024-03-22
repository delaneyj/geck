package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/delaneyj/geck/generator"
	geckpb "github.com/delaneyj/geck/pb/gen/geck/v1"
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
	b, err := os.ReadFile("./ecs.gen.json")
	if err != nil {
		return fmt.Errorf("failed to read bundle definition: %w", err)
	}
	if err := opts.UnmarshalJSON(b); err != nil {
		return fmt.Errorf("failed to unmarshal bundle definition: %w", err)
	}
	return generator.BuildECS(ctx, opts)
}
