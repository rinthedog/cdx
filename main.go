package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

type Aliases struct {
	Paths map[string]string `json:"paths"`
}

var (
	aliasFile string
	aliases   Aliases
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}
	aliasFile = filepath.Join(home, ".cdx_aliases.json")
	aliases.Paths = make(map[string]string)
	loadAliases()
	cobra.EnableCommandSorting = false
}

func loadAliases() {
	data, err := os.ReadFile(aliasFile)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("Error reading aliases file:", err)
		}
		return
	}

	if err := json.Unmarshal(data, &aliases); err != nil {
		fmt.Println("Error parsing aliases file:", err)
	}
}

func saveAliases() error {
	data, err := json.MarshalIndent(aliases, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(aliasFile, data, 0644)
}

var rootCmd = &cobra.Command{
	Use:   "cdx [alias]",
	Short: "CDX - Directory alias manager",
	Long:  `CDX helps you manage directory aliases for quick navigation.`,
	Args:  cobra.ArbitraryArgs,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		alias := args[0]
		if path, exists := aliases.Paths[alias]; exists {
			explorer := exec.Command("explorer.exe", path)
			if err := explorer.Start(); err != nil {
				fmt.Printf("Error opening directory: %v\n", err)
				os.Exit(1)
			}
			return
		}
		fmt.Printf("Alias '%s' not found\n", alias)
		os.Exit(1)
	},
}

var setCmd = &cobra.Command{
	Use:   "set [alias] [path]",
	Short: "Set an alias for a directory path",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		alias := args[0]
		path := args[1]

		absPath, err := filepath.Abs(path)
		if err != nil {
			fmt.Println("Error resolving path:", err)
			return
		}

		if info, err := os.Stat(absPath); err != nil || !info.IsDir() {
			fmt.Println("Invalid directory path:", path)
			return
		}

		aliases.Paths[alias] = absPath
		if err := saveAliases(); err != nil {
			fmt.Println("Error saving alias:", err)
			return
		}
		fmt.Printf("Set alias '%s' to '%s'\n", alias, absPath)
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove [alias]",
	Short: "Remove a directory alias",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		alias := args[0]
		if _, exists := aliases.Paths[alias]; !exists {
			fmt.Printf("Alias '%s' does not exist\n", alias)
			return
		}

		delete(aliases.Paths, alias)
		if err := saveAliases(); err != nil {
			fmt.Println("Error removing alias:", err)
			return
		}
		fmt.Printf("Removed alias '%s'\n", alias)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all directory aliases",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(aliases.Paths) == 0 {
			fmt.Println("No aliases set")
			return
		}

		fmt.Println("Directory Aliases:")
		for alias, path := range aliases.Paths {
			fmt.Printf("  %s -> %s\n", alias, path)
		}
	},
}

func main() {
	rootCmd.AddCommand(setCmd, removeCmd, listCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
