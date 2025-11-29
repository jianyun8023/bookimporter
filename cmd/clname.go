package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jianyun8023/bookimporter/pkg/ui"
	"github.com/jianyun8023/bookimporter/pkg/util"
	"github.com/kapmahc/epub"
	"github.com/spf13/cobra"
)

// Used for downloading books from sanqiu website.
var c = &ClnameConfig{}

// renameBookCmd used for download books from sanqiu.cc
var clnameCmd = &cobra.Command{
	Use:   "clname",
	Short: "清理书籍标题中的无用描述",
	Long: `清理 EPUB 书籍标题中的无用描述符和标记

自动移除书籍标题中的各种括号标记，如：（）【】()[]
使用 Calibre 的 ebook-meta 工具修改元数据。

支持：
  • 单个文件或批量目录处理
  • 递归搜索子目录
  • 预览模式（不实际修改）
  • 自动处理损坏的 EPUB 文件`,
	Example: `  # 清理单个文件
  bookimporter clname -p book.epub

  # 批量清理当前目录（不包含子目录）
  bookimporter clname -p /path/to/books/

  # 递归清理所有子目录
  bookimporter clname -p /path/to/books/ -r

  # 预览模式（不实际修改）
  bookimporter clname -p /path/to/books/ -r -t

  # 自动移动损坏文件
  bookimporter clname -p /path/to/books/ -r --move-corrupted-to /path/to/corrupted/`,
	Run: func(cmd *cobra.Command, args []string) {
		// Validate config.
		ValidateConfig(c)

		// 打印头部
		fmt.Println(ui.RenderHeader("清理书籍标题", "移除标题中的无用描述符和标记"))
		fmt.Println()

		// 统计信息
		stats := &ClnameStats{
			Total:   0,
			Updated: 0,
			Skipped: 0,
			Failed:  0,
		}

		if util.IsDir(c.Path) {
			// 根据 Recursive 参数决定是否递归搜索 EPUB 文件
			var m []string
			var err error

			if c.Recursive {
				// 递归搜索所有子目录
				err = filepath.Walk(c.Path, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".epub") {
						m = append(m, path)
					}
					return nil
				})
			} else {
				// 仅搜索当前目录
				pattern := filepath.Join(c.Path, "*.epub")
				m, err = filepath.Glob(pattern)
			}

			if err != nil {
				fmt.Println(ui.RenderError(fmt.Sprintf("扫描目录失败: %v", err)))
				return
			}

			stats.Total = len(m)

			if stats.Total == 0 {
				fmt.Println(ui.RenderWarning("未找到 EPUB 文件"))
				return
			}

			fmt.Println(ui.RenderInfo(fmt.Sprintf("找到 %d 个 EPUB 文件", stats.Total)))
			fmt.Println()

			// 创建增强的进度跟踪器
			progress := ui.NewCompactProgressTracker(stats.Total)
			progress.SetShowMessage(true)

			for i, val := range m {
				epubpath := val
				progress.SetMessage(filepath.Base(epubpath))

				// 显示进度
				if stats.Total > 1 {
					fmt.Printf("\r%s", progress.RenderCompact())
				}

				err := ParseEpub(epubpath, c, stats, progress)
				if err != nil {
					// 清除进度行
					if stats.Total > 1 {
						fmt.Print("\r" + strings.Repeat(" ", 120) + "\r")
					}

					// 记录错误并继续处理下一个文件
					fmt.Println(ui.FormatFilePath("文件", epubpath))
					fmt.Println(ui.RenderWarning(fmt.Sprintf("跳过: %v", err)))
					fmt.Println()
					stats.Failed++
					if progress != nil {
						progress.IncrementFailure()
					}
				}

				// 清除进度行，为文件详情腾出空间
				if stats.Total > 1 && i < stats.Total-1 {
					fmt.Print("\r" + strings.Repeat(" ", 120) + "\r")
				}
			}

			// 清除最后的进度行并显示最终统计
			if stats.Total > 1 {
				fmt.Print("\r" + strings.Repeat(" ", 120) + "\r")
				fmt.Println(progress.RenderWithStats())
				fmt.Println()
			}

		} else {
			stats.Total = 1
			epubpath := c.Path
			err := ParseEpub(epubpath, c, stats, nil)
			if err != nil {
				// 记录错误
				fmt.Println(ui.FormatFilePath("文件", epubpath))
				fmt.Println(ui.RenderWarning(fmt.Sprintf("处理失败: %v", err)))
				stats.Failed++
			}
		}

		// 打印统计信息
		printClnameStats(stats)

		// 如果有失败且未设置忽略错误，设置退出码为 1（便于脚本检测）
		if stats.Failed > 0 && !c.IgnoreErrors {
			os.Exit(1)
		}
	},
}

func ValidateConfig(c *ClnameConfig) {
	// 验证路径是否存在
	if !util.Exists(c.Path) {
		fmt.Println(ui.RenderError("文件路径不存在，请检查"))
		os.Exit(1)
	}

	// 验证文件格式
	if util.IsFile(c.Path) && !strings.HasSuffix(c.Path, ".epub") {
		fmt.Println(ui.RenderError("文件格式不正确，必须是 .epub 文件"))
		os.Exit(1)
	}

	// 验证互斥参数
	if c.MoveCorruptedTo != "" && c.DeleteCorrupted {
		fmt.Println(ui.RenderError("--move-corrupted-to 和 --delete-corrupted 参数不能同时使用"))
		fmt.Println(ui.RenderInfo("提示: 请选择移动或删除损坏文件，不能同时进行"))
		os.Exit(1)
	}

	// 验证 force-delete 的使用
	if c.ForceDelete && !c.DeleteCorrupted {
		fmt.Println(ui.RenderWarning("警告: --force-delete 参数需要配合 --delete-corrupted 使用"))
		fmt.Println(ui.RenderInfo("提示: --force-delete 用于在删除损坏文件时跳过确认步骤"))
	}

	// 验证 move-corrupted-to 目标目录
	if c.MoveCorruptedTo != "" {
		// 确保目标目录不是源目录的子目录
		absPath, err := filepath.Abs(c.Path)
		if err == nil {
			absDst, err := filepath.Abs(c.MoveCorruptedTo)
			if err == nil && strings.HasPrefix(absDst, absPath) {
				fmt.Println(ui.RenderError("错误: 目标目录不能是源目录的子目录"))
				os.Exit(1)
			}
		}
	}
}

// ClnameStats 清理标题统计
type ClnameStats struct {
	Total   int
	Updated int
	Skipped int
	Failed  int
}

// printClnameStats 打印统计信息
func printClnameStats(stats *ClnameStats) {
	fmt.Println()
	fmt.Println(ui.RenderSeparator(60))
	fmt.Println()

	// 使用新的表格组件
	tableConfig := ui.NewTableConfig()
	tableConfig.Headers = []string{"  状态  ", " 数量 ", " 百分比 "}
	tableConfig.BorderStyle = "rounded"
	tableConfig.AlignRight = []int{1, 2}

	var rows [][]string

	// 已更新
	if stats.Updated > 0 {
		percentage := float64(stats.Updated) / float64(stats.Total) * 100
		rows = append(rows, []string{
			ui.IconSuccess + " 已更新 ",
			fmt.Sprintf(" %d ", stats.Updated),
			fmt.Sprintf(" %.1f%% ", percentage),
		})
	}

	// 跳过（无需修改）
	if stats.Skipped > 0 {
		percentage := float64(stats.Skipped) / float64(stats.Total) * 100
		rows = append(rows, []string{
			ui.IconSkip + " 跳过  ",
			fmt.Sprintf(" %d ", stats.Skipped),
			fmt.Sprintf(" %.1f%% ", percentage),
		})
	}

	// 失败
	if stats.Failed > 0 {
		percentage := float64(stats.Failed) / float64(stats.Total) * 100
		rows = append(rows, []string{
			ui.IconError + " 失败  ",
			fmt.Sprintf(" %d ", stats.Failed),
			fmt.Sprintf(" %.1f%% ", percentage),
		})
	}

	// 总计
	rows = append(rows, []string{
		"  总计  ",
		fmt.Sprintf(" %d ", stats.Total),
		" 100% ",
	})

	tableConfig.Rows = rows
	table := ui.NewTable(tableConfig)
	fmt.Println(table.Render())
	fmt.Println()

	// 添加总结信息
	if stats.Updated == 0 && stats.Failed == 0 {
		fmt.Println(ui.RenderInfo("✨ 所有文件标题都已是最佳状态"))
	} else if stats.Updated > 0 {
		fmt.Println(ui.RenderSuccess(fmt.Sprintf("✨ 成功更新 %d 个文件的标题", stats.Updated)))
	}
	if stats.Failed > 0 {
		fmt.Println(ui.RenderWarning(fmt.Sprintf("⚠️  %d 个文件处理失败", stats.Failed)))
	}
}

func init() {
	// 基础参数
	clnameCmd.Flags().StringVarP(&c.Path, "path", "p", "./",
		"指定要处理的 EPUB 文件或目录路径")
	clnameCmd.Flags().BoolVarP(&c.Recursive, "recursive", "r", false,
		"递归搜索子目录中的所有 EPUB 文件")

	// 运行模式
	clnameCmd.Flags().BoolVarP(&c.DoTry, "dotry", "t", false,
		"预览模式，显示将要进行的修改但不实际执行")
	clnameCmd.Flags().BoolVarP(&c.IgnoreErrors, "ignore-errors", "i", false,
		"忽略错误，即使有失败也返回退出码 0")

	// 损坏文件处理（互斥选项）
	clnameCmd.Flags().StringVar(&c.MoveCorruptedTo, "move-corrupted-to", "",
		"将损坏的 EPUB 文件移动到指定目录（与 --delete-corrupted 互斥）")
	clnameCmd.Flags().BoolVar(&c.DeleteCorrupted, "delete-corrupted", false,
		"直接删除损坏的 EPUB 文件（与 --move-corrupted-to 互斥）")
	clnameCmd.Flags().BoolVar(&c.ForceDelete, "force-delete", false,
		"删除损坏文件时不需要用户确认（需配合 --delete-corrupted 使用）")

	// 调试选项
	clnameCmd.Flags().BoolVarP(&c.Debug, "debug", "d", false,
		"启用调试模式，显示详细的执行信息")
}

func ParseEpub(file string, c *ClnameConfig, stats *ClnameStats, progress *ui.ProgressTracker) error {
	// 预先检测 EPUB 文件完整性
	if err := util.ValidateEpubFile(file); err != nil {
		if c.Debug {
			fmt.Println(ui.RenderError(fmt.Sprintf("EPUB 文件检测失败: %v", err)))
		}

		// 处理损坏的文件
		if c.MoveCorruptedTo != "" {
			if c.DoTry {
				fmt.Println(ui.RenderInfo(fmt.Sprintf("[试运行] 将移动损坏文件到: %s", c.MoveCorruptedTo)))
			} else {
				newPath, moveErr := util.MoveFileWithConflictHandling(file, c.MoveCorruptedTo)
				if moveErr != nil {
					fmt.Println(ui.RenderError(fmt.Sprintf("移动损坏文件失败: %v", moveErr)))
				} else {
					fmt.Println(ui.RenderInfo(fmt.Sprintf("已移动损坏文件到: %s", newPath)))
				}
			}
		} else if c.DeleteCorrupted {
			if c.DoTry {
				fmt.Println(ui.RenderInfo("[试运行] 将删除损坏文件"))
			} else {
				needConfirm := !c.ForceDelete
				deleteErr := util.SafeDeleteFile(file, needConfirm)
				if deleteErr != nil {
					fmt.Println(ui.RenderError(fmt.Sprintf("删除损坏文件失败: %v", deleteErr)))
				} else {
					fmt.Println(ui.RenderInfo("已删除损坏文件"))
				}
			}
		}

		return fmt.Errorf("EPUB 文件检测失败: %w", err)
	}

	book, err := epub.Open(file)
	if err != nil {
		return err
	}
	if book == nil || len(book.Opf.Metadata.Title) == 0 {
		return fmt.Errorf("无法获得书籍标题")
	}
	title := book.Opf.Metadata.Title[0]
	newTitle := util.TryCleanTitle(title)
	if title == newTitle {
		stats.Skipped++
		if progress != nil {
			progress.IncrementSkipped()
		}
		return nil
	}

	// 美化输出
	fmt.Println(ui.FormatFilePath("路径", file))
	fmt.Println(ui.FormatFileOperation("标题", title, newTitle))

	if c.DoTry {
		fmt.Println(ui.RenderInfo("[试运行] 将更新标题"))
		fmt.Println()
		stats.Skipped++
		if progress != nil {
			progress.IncrementSkipped()
		}
		return nil
	}

	// 执行单个shell命令时, 直接运行即可
	command := fmt.Sprintf("ebook-meta \"%s\" -t \"%s\"", file, newTitle)
	cmdx := exec.Command("bash", "-c", command)

	if _, err = cmdx.Output(); err != nil {
		if c.Debug {
			fmt.Println(cmdx.Env)
			fmt.Println(cmdx.Args)
		}
		return err
	}

	fmt.Println(ui.RenderSuccess("已更新"))
	fmt.Println()
	stats.Updated++
	if progress != nil {
		progress.IncrementSuccess()
	}
	return nil
}

type ClnameConfig struct {
	Path            string
	Recursive       bool // 是否递归搜索子目录
	DoTry           bool
	Debug           bool
	IgnoreErrors    bool   // 忽略错误，有失败也返回 0
	MoveCorruptedTo string // 损坏文件移动目标目录
	DeleteCorrupted bool   // 是否删除损坏文件
	ForceDelete     bool   // 删除时不需要确认
}
