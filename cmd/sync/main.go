package main

import (
	"context"
	"fmt"
	"os"

	"omnix/internal/catalog"
	"omnix/internal/config"
)

func main() {
	cfg := config.Load()
	items, err := catalog.Fetch(context.Background())
	if err == nil {
		err = catalog.Save(cfg.CatalogFile, items)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("%d cosmetics -> %s\n", len(items), cfg.CatalogFile)
}
