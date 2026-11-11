package main

import (
	"context"
	"fmt"
	"os"

	"omnix/internal/catalog"
	"omnix/internal/config"
	"omnix/internal/keychain"
)

func main() {
	if err := run(context.Background(), config.Load()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, cfg config.Config) error {
	items, err := catalog.Fetch(ctx)
	if err != nil {
		return err
	}
	if err := catalog.Save(cfg.CatalogFile, items); err != nil {
		return err
	}
	fmt.Printf("%d cosmetics -> %s\n", len(items), cfg.CatalogFile)

	keys, err := keychain.Fetch(ctx)
	if err != nil {
		return err
	}
	if err := keychain.Save(cfg.KeychainFile, keys); err != nil {
		return err
	}
	fmt.Printf("%d keys -> %s\n", len(keys), cfg.KeychainFile)
	return nil
}
