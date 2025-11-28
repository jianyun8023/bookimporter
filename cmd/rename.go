package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jianyun8023/bookimporter/pkg/ui"
	"github.com/spf13/cobra"
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
			fmt.Println(ui.RenderError("需要指定扫描路径"))
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

		// 打印头部
		fmt.Println(ui.RenderHeader("批量重命名文件", "按模板批量重命名或移动文件"))
		fmt.Println()

		if config.Debug {
			fmt.Println(ui.RenderInfo("调试信息:"))
			fmt.Printf("  - 试运行: %v\n", config.DoTry)
			fmt.Printf("  - 文件格式: %v\n", config.Formats)
			fmt.Printf("  - 递归搜索: %v\n", config.Recursive)
			fmt.Printf("  - 源目录: %s\n", config.SourceDir)
			fmt.Printf("  - 输出目录: %s\n", config.OutputDir)
			fmt.Printf("  - 模板: %s\n", config.Template)
			fmt.Printf("  - 起始序号: %d\n", config.StartIndex)
			fmt.Println()
		}
		rename(config)
	},
}

func validateConfig(config *RenameConfig) {
	if !strings.Contains(config.Template, "@n") {
		fmt.Println(ui.RenderError(fmt.Sprintf("模板 [%s] 中不存在占位符 @n", config.Template)))
		os.Exit(1)
	}
}

func rename(config *RenameConfig) {

	files, err := findFiles(config.SourceDir, config.Formats, config.Recursive)
	if err != nil {
		fmt.Println(ui.RenderError(fmt.Sprintf("查找文件失败: %s", err)))
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println(ui.RenderWarning("未找到匹配的文件"))
		return
	}

	fmt.Println(ui.RenderInfo(fmt.Sprintf("找到 %d 个文件", len(files))))
	fmt.Println()

	var (
		renamedFiles []string
		movedFiles   []string
	)

	// 如果是试运行模式，显示预览表格
	if config.DoTry {
		fmt.Println(ui.RenderTitle("重命名预览"))
		fmt.Println()
	}

	// 创建进度跟踪器
	var progress *ui.ProgressTracker
	if !config.DoTry && len(files) > 1 {
		progress = ui.NewProgressTracker(len(files))
	}

	for i, file := range files {
		newName := buildNewName(config.Template, config.StartIndex+i, file)

		if config.OutputDir != "" {
			outputPath := filepath.Join(config.OutputDir, newName)

			if config.DoTry {
				// 预览模式
				fmt.Println(ui.FormatRenamePreview(file, outputPath))
			} else {
				// 显示进度
				if progress != nil {
					progress.SetMessage(fmt.Sprintf("重命名: %s", filepath.Base(file)))
					fmt.Printf("\r%s", progress.RenderSimple())
				}

				err = os.Rename(file, outputPath)
				if err != nil {
					if progress != nil {
						fmt.Print("\r" + strings.Repeat(" ", 80) + "\r")
					}
					fmt.Println(ui.RenderError(fmt.Sprintf("重命名失败: %s", err)))
					os.Exit(1)
				}

				if progress != nil {
					fmt.Print("\r" + strings.Repeat(" ", 80) + "\r")
				}
				fmt.Println(ui.FormatRenamePreview(file, outputPath))

				if progress != nil {
					progress.Increment()
				}
			}

			movedFiles = append(movedFiles, file+" -> "+outputPath)
		} else {
			outputPath := filepath.Join(filepath.Dir(file), newName)

			if config.DoTry {
				// 预览模式
				fmt.Println(ui.FormatRenamePreview(file, outputPath))
			} else {
				// 显示进度
				if progress != nil {
					progress.SetMessage(fmt.Sprintf("重命名: %s", filepath.Base(file)))
					fmt.Printf("\r%s", progress.RenderSimple())
				}

				err = os.Rename(file, outputPath)
				if err != nil {
					if progress != nil {
						fmt.Print("\r" + strings.Repeat(" ", 80) + "\r")
					}
					fmt.Println(ui.RenderError(fmt.Sprintf("重命名失败: %s", err)))
					os.Exit(1)
				}

				if progress != nil {
					fmt.Print("\r" + strings.Repeat(" ", 80) + "\r")
				}
				fmt.Println(ui.FormatRenamePreview(file, outputPath))

				if progress != nil {
					progress.Increment()
				}
			}

			renamedFiles = append(renamedFiles, file+" -> "+outputPath)
		}
	}

	// 清除进度行
	if progress != nil {
		fmt.Print("\r" + strings.Repeat(" ", 80) + "\r")
	}

	// 打印统计信息
	fmt.Println()
	fmt.Println(ui.RenderSeparator(50))
	fmt.Println()

	statsMap := map[string]int{
		"total": len(files),
	}

	if config.OutputDir != "" {
		fmt.Println(ui.RenderStatsSummary(statsMap))
		fmt.Println()
		if config.DoTry {
			fmt.Println(ui.RenderInfo(fmt.Sprintf("[试运行] 将移动 %d 个文件到: %s", len(movedFiles), config.OutputDir)))
		} else {
			fmt.Println(ui.RenderSuccess(fmt.Sprintf("成功移动 %d 个文件到: %s", len(movedFiles), config.OutputDir)))
		}
	} else {
		fmt.Println(ui.RenderStatsSummary(statsMap))
		fmt.Println()
		if config.DoTry {
			fmt.Println(ui.RenderInfo(fmt.Sprintf("[试运行] 将重命名 %d 个文件", len(renamedFiles))))
		} else {
			fmt.Println(ui.RenderSuccess(fmt.Sprintf("成功重命名 %d 个文件", len(renamedFiles))))
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
