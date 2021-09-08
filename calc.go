package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

const (
	mul  = 25214903917
	mask = (1 << 48) - 1
)

func fetchLinkIJ(seed int64) (int64, int64) {
	result := seed ^ mul&mask
	link := make([]int64, 4)
	for i := range link {
		result = (result*mul + 11) & mask
		if i&1 == 0 {
			link[i] = result >> 16 << 32
		} else {
			link[i] = result << 16 >> 32
		}
	}
	linkI := (link[0] + link[1]) | 1
	linkJ := (link[2] + link[3]) | 1
	return linkI, linkJ
}

func calc(seed int64, magic int64, linkI int64, linkJ int64, xBlock int64, zBlock int64) (int64, int64) {
	result := (16*(xBlock*linkI+zBlock*linkJ) ^ seed + magic) ^ mul&mask
	link := make([]int64, 2)
	for i := range link {
		result = (result*mul + 11) & mask
		link[i] = result >> 44
	}
	return link[0] + 16*xBlock, link[1] + 16*zBlock
}

// Examples Multi-example
func Examples(values ...string) string {
	// Add 2 spaces
	for i, value := range values {
		values[i] = "  " + value
	}
	return strings.Join(values, "\n")
}

type (
	Config struct {
		// Version in (1.16, 1.17)
		Version string
		// Seed of world
		Seed int64
	}
)

func main() {
	cfg := new(Config)
	cmd := &cobra.Command{
		Use: "mcalc [flags] command",
	}
	diamond := &cobra.Command{
		Use:   "diamond [flags] x_block z_block",
		Short: "Calc diamond position",
		Example: Examples(
			"mcalc diamond --seed=123 2 2",
			"mcalc diamond --version=1.17 --seed=123 2 3",
			"mcalc diamond --version=1.16 --seed=123 1 1",
		),
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			xBlock, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				panic(err)
			}
			zBlock, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				panic(err)
			}
			linkI, linkJ := fetchLinkIJ(cfg.Seed)
			var magic int64 = 60011
			if cfg.Version == "1.16" {
				magic = 60009
			}
			x, z := calc(cfg.Seed, magic, linkI, linkJ, xBlock, zBlock)
			fmt.Printf("Diamond at x:%d, z:%d\n", x, z)
		},
	}
	diamond.Flags().StringVar(&cfg.Version, "version", "1.17", "MC version in (1.16, 1.17)")
	diamond.Flags().Int64Var(&cfg.Seed, "seed", 123, "Seed of world")
	err := diamond.MarkFlagRequired("seed")
	if err != nil {
		panic(err)
	}
	cmd.AddCommand(diamond)
	_ = cmd.Execute()
}
