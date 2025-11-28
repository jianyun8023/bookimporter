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

	// 如果是试运行模式，使用表格显示预览
	if config.DoTry {
		fmt.Println(ui.RenderTitle("📋 重命名预览"))
		fmt.Println()

		// 创建预览表格
		tableConfig := ui.NewTableConfig()
		tableConfig.Headers = []string{" # ", " 原文件名 ", " → ", " 新文件名 "}
		tableConfig.BorderStyle = "rounded"
		tableConfig.CompactMode = false

		var rows [][]string
		for i, file := range files {
			newName := buildNewName(config.Template, config.StartIndex+i, file)
			var outputPath string
			if config.OutputDir != "" {
				outputPath = filepath.Join(config.OutputDir, newName)
			} else {
				outputPath = filepath.Join(filepath.Dir(file), newName)
			}

			// 截断长文件名
			oldName := file
			if len(oldName) > 40 {
				oldName = "..." + oldName[len(oldName)-37:]
			}
			newPath := outputPath
			if len(newPath) > 40 {
				newPath = "..." + newPath[len(newPath)-37:]
			}

			rows = append(rows, []string{
				fmt.Sprintf(" %d ", i+1),
				fmt.Sprintf(" %s ", oldName),
				" → ",
				fmt.Sprintf(" %s ", newPath),
			})

			// 限制预览显示的行数
			if i >= 19 && len(files) > 20 {
				rows = append(rows, []string{
					" ... ",
					fmt.Sprintf(" ... 还有 %d 个文件 ... ", len(files)-20),
					"   ",
					" ... ",
				})
				break
			}
		}

		tableConfig.Rows = rows
		table := ui.NewTable(tableConfig)
		fmt.Println(table.Render())
		fmt.Println()
	}

	// 创建进度跟踪器
	var progress *ui.ProgressTracker
	if !config.DoTry && len(files) > 1 {
		progress = ui.NewCompactProgressTracker(len(files))
		progress.SetShowMessage(true)
	}

	// 如果不是预览模式，执行重命名
	if !config.DoTry {
		for i, file := range files {
			newName := buildNewName(config.Template, config.StartIndex+i, file)

			if config.OutputDir != "" {
				outputPath := filepath.Join(config.OutputDir, newName)

				// 显示进度
				if progress != nil {
					progress.SetMessage(filepath.Base(file))
					fmt.Printf("\r%s", progress.RenderCompact())
				}

				err = os.Rename(file, outputPath)
				if err != nil {
					if progress != nil {
						fmt.Print("\r" + strings.Repeat(" ", 120) + "\r")
					}
					fmt.Println(ui.RenderError(fmt.Sprintf("重命名失败: %s", err)))
					os.Exit(1)
				}

				if progress != nil {
					fmt.Print("\r" + strings.Repeat(" ", 120) + "\r")
					progress.IncrementSuccess()
				}
				fmt.Println(ui.FormatRenamePreview(file, outputPath))

				movedFiles = append(movedFiles, file+" -> "+outputPath)
			} else {
				outputPath := filepath.Join(filepath.Dir(file), newName)

				// 显示进度
				if progress != nil {
					progress.SetMessage(filepath.Base(file))
					fmt.Printf("\r%s", progress.RenderCompact())
				}

				err = os.Rename(file, outputPath)
				if err != nil {
					if progress != nil {
						fmt.Print("\r" + strings.Repeat(" ", 120) + "\r")
					}
					fmt.Println(ui.RenderError(fmt.Sprintf("重命名失败: %s", err)))
					os.Exit(1)
				}

				if progress != nil {
					fmt.Print("\r" + strings.Repeat(" ", 120) + "\r")
					progress.IncrementSuccess()
				}
				fmt.Println(ui.FormatRenamePreview(file, outputPath))

				renamedFiles = append(renamedFiles, file+" -> "+outputPath)
			}
		}

		// 清除进度行并显示最终统计
		if progress != nil {
			fmt.Print("\r" + strings.Repeat(" ", 120) + "\r")
			fmt.Println(progress.RenderWithStats())
			fmt.Println()
		}
	}

	// 打印统计信息
	fmt.Println()
	fmt.Println(ui.RenderSeparator(60))
	fmt.Println()

	// 使用表格展示统计
	tableConfig := ui.NewTableConfig()
	tableConfig.Headers = []string{"  项目  ", " 值 "}
	tableConfig.BorderStyle = "rounded"
	tableConfig.AlignRight = []int{1}

	var rows [][]string
	rows = append(rows, []string{
		" 文件总数 ",
		fmt.Sprintf(" %d ", len(files)),
	})

	if config.OutputDir != "" {
		rows = append(rows, []string{
			" 目标目录 ",
			fmt.Sprintf(" %s ", config.OutputDir),
		})
	}

	if config.DoTry {
		rows = append(rows, []string{
			" 模式 ",
			" 预览模式 ",
		})
	} else {
		rows = append(rows, []string{
			" 已处理 ",
			fmt.Sprintf(" %d ", len(files)),
		})
	}

	tableConfig.Rows = rows
	table := ui.NewTable(tableConfig)
	fmt.Println(table.Render())
	fmt.Println()

	if config.OutputDir != "" {
		if config.DoTry {
			fmt.Println(ui.RenderInfo(fmt.Sprintf("📝 [试运行] 将移动 %d 个文件到: %s", len(files), config.OutputDir)))
		} else {
			fmt.Println(ui.RenderSuccess(fmt.Sprintf("✨ 成功移动 %d 个文件到: %s", len(movedFiles), config.OutputDir)))
		}
	} else {
		if config.DoTry {
			fmt.Println(ui.RenderInfo(fmt.Sprintf("📝 [试运行] 将重命名 %d 个文件", len(files))))
		} else {
			fmt.Println(ui.RenderSuccess(fmt.Sprintf("✨ 成功重命名 %d 个文件", len(renamedFiles))))
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
