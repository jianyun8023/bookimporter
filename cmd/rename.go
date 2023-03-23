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
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			fmt.Printf("Error: 需要指定扫描路径\n")
			os.Exit(1)
		}
		array, err := cmd.Flags().GetStringArray("format")
		if err != nil {
			panic(err)
		}
		config := &RenameConfig{
			Debug:      cmd.Flag("debug").Value.String() == "true",
			DoTry:      cmd.Flag("do-try").Value.String() == "true",
			Formats:    array,
			Recursive:  cmd.Flag("recursive").Value.String() == "true",
			SourceDir:  args[0],
			OutputDir:  cmd.Flag("output").Value.String(),
			Template:   cmd.Flag("template").Value.String(),
			StartIndex: parseIntFlag(cmd, "start-num"),
		}

		validateConfig(config)

		if config.Debug {
			fmt.Println("Debugging information:")
			fmt.Printf("  - DoTry: %v\n", config.DoTry)
			fmt.Printf("  - Format: %v\n", config.Formats)
			fmt.Printf("  - Recursive: %v\n", config.Recursive)
			fmt.Printf("  - SourceDir: %s\n", config.SourceDir)
			fmt.Printf("  - OutputDir: %s\n", config.OutputDir)
			fmt.Printf("  - Template: %s\n", config.Template)
			fmt.Printf("  - StartIndex: %d\n", config.StartIndex)
		}
		rename(config)
	},
}

func validateConfig(config *RenameConfig) {
	if !strings.Contains(config.Template, "@n") {
		panic(fmt.Errorf("模板[ %s ]中不存在占位符@n", config.Template))
	}

}

func rename(config *RenameConfig) {

	files, err := findFiles(config.SourceDir, config.Formats, config.Recursive)
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

		if config.OutputDir != "" {
			outputPath := filepath.Join(config.OutputDir, newName)

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
	if config.OutputDir != "" {
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

func findFiles(dir string, formats []string, recursive bool) ([]string, error) {
	var files []string

	matches, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return nil, err
	}

	for _, match := range matches {

		included := false
		for _, includedStr := range formats {
			if strings.Contains(match, includedStr) {
				included = true
				break
			}
		}

		info, err := os.Stat(match)
		if err != nil {
			return nil, err
		}

		if !info.IsDir() && included {
			files = append(files, match)
		} else if recursive {
			subFiles, err := findFiles(match, formats, true)
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
	renameCmd.Flags().StringArrayP("format", "f", []string{"*"}, "File format to match (e.g. 'txt')")
	renameCmd.Flags().BoolP("recursive", "r", false, "Recursively search for files")
	renameCmd.Flags().StringP("output", "o", "", "Output directory for moved files")
	renameCmd.Flags().StringP("template", "t", "file-@n", "Template for new filename (e.g. 'file-@n')")
	renameCmd.Flags().Int("start-num", 1, "Starting number for sequence")
	_ = rootCmd.MarkFlagRequired("format")
	_ = rootCmd.MarkFlagRequired("template")
}

func buildNewName(template string, index int, file string) string {
	ext := filepath.Ext(file)
	newName := strings.Replace(template, "@n", strconv.Itoa(index), -1)
	newName += ext
	return newName
}

func parseIntFlag(cmd *cobra.Command, name string) int {
	val, _ := cmd.Flags().GetInt(name)
	return val
}

type RenameConfig struct {
	Debug      bool
	DoTry      bool
	Formats    []string
	Recursive  bool
	SourceDir  string
	OutputDir  string
	Template   string
	StartIndex int
}
