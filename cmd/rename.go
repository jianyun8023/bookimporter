package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename or move files according to a template",
	Long: `Rename or move files according to a template.

The rename command will rename or move files according to a template that you
specify. The template can include the file extension, so if you want to keep
the original extension, you can include it in the template. You can also use
a sequence number in the template to number the files sequentially.`,
	Run: func(cmd *cobra.Command, args []string) {

		config := RenameConfig{
			Debug:      cmd.Flag("debug").Value.String() == "true",
			DoTry:      cmd.Flag("do-try").Value.String() == "true",
			Format:     cmd.Flag("format").Value.String(),
			Recursive:  cmd.Flag("recursive").Value.String() == "true",
			Move:       cmd.Flag("move").Value.String() == "true",
			OutputDir:  cmd.Flag("output").Value.String(),
			Template:   cmd.Flag("template").Value.String(),
			StartIndex: parseIntFlag(cmd, "start-num"),
		}

		if config.Debug {
			fmt.Println("Debugging information:")
			fmt.Printf("  - Format: %s\n", config.Format)
			fmt.Printf("  - Recursive: %v\n", config.Recursive)
			fmt.Printf("  - Move: %v\n", config.Move)
			fmt.Printf("  - OutputDir: %s\n", config.OutputDir)
			fmt.Printf("  - Template: %s\n", config.Template)
			fmt.Printf("  - StartIndex: %d\n", config.StartIndex)
		}

		rename(config)

	},
}

func rename(config RenameConfig) {

	files, err := findFiles(config.Format, config.Recursive)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No files found.")
		return
	}

	var (
		renamedFiles []string
		movedFiles   []string
	)

	for i, file := range files {
		newName := buildNewName(config.Template, config.StartIndex+i, file)

		if config.Move {
			outputDir := config.OutputDir
			if outputDir == "" {
				outputDir = filepath.Dir(file)
			}

			outputPath := filepath.Join(outputDir, newName)

			if !config.DoTry {
				err = os.Rename(file, outputPath)
				if err != nil {
					fmt.Printf("Error: %s\n", err)
					os.Exit(1)
				}
			}

			movedFiles = append(movedFiles, file+" -> "+outputPath)
		} else {
			outputPath := filepath.Join(filepath.Dir(file), newName)

			if !config.DoTry {
				err = os.Rename(file, outputPath)
				if err != nil {
					fmt.Printf("Error: %s\n", err)
					os.Exit(1)
				}
			}

			renamedFiles = append(renamedFiles, file+" -> "+outputPath)
		}
	}
	fmt.Printf("%d files found.\n", len(files))
	if config.Move {
		fmt.Printf("%d files moved:\n", len(movedFiles))
		for _, file := range movedFiles {
			fmt.Printf("  - %s\n", file)
		}
	} else {
		fmt.Printf("%d files renamed:\n", len(renamedFiles))
		for _, file := range renamedFiles {
			fmt.Printf("  - %s\n", file)
		}
	}
}

func findFiles(format string, recursive bool) ([]string, error) {
	var files []string

	matches, err := filepath.Glob(format)
	if err != nil {
		return nil, err
	}

	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			return nil, err
		}

		if !info.IsDir() {
			files = append(files, match)
		} else if recursive {
			subFiles, err := findFiles(filepath.Join(match, format), true)
			if err != nil {
				return nil, err
			}

			files = append(files, subFiles...)
		}
	}

	return files, nil

}

func init() {
	renameCmd.Flags().Bool("debug", false, "Enable debugging information")
	renameCmd.Flags().Bool("do-try", false, "Only print out actions that would be performed")
	renameCmd.Flags().StringP("format", "f", "*", "File format to match (e.g. '*.txt')")
	renameCmd.Flags().BoolP("recursive", "r", false, "Recursively search for files")
	renameCmd.Flags().BoolP("move", "m", false, "Move files instead of renaming them")
	renameCmd.Flags().String("output", "", "Output directory for moved files")
	renameCmd.Flags().String("template", "file-$n", "Template for new filename")
	renameCmd.Flags().Int("start-num", 1, "Starting number for sequence")
}

func buildNewName(template string, index int, file string) string {
	ext := filepath.Ext(file)
	newName := strings.Replace(template, "$n", strconv.Itoa(index), -1)
	newName = strings.Replace(newName, "$e", ext, -1)

	if ext != "" && !strings.Contains(newName, "$e") {
		newName += ext
	}

	return newName
}

func parseIntFlag(cmd *cobra.Command, name string) int {
	val, _ := cmd.Flags().GetInt(name)
	return val
}

type RenameConfig struct {
	Debug      bool
	DoTry      bool
	Format     string
	Recursive  bool
	Move       bool
	OutputDir  string
	Template   string
	StartIndex int
}
