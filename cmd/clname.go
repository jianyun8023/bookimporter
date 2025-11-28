package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
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
			m, _ := filepath.Glob(path.Join(c.Path, "*.epub"))
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
				if err != nil && !c.Skip {
					// 清除进度行
					if stats.Total > 1 {
						fmt.Print("\r" + strings.Repeat(" ", 120) + "\r")
					}
					fmt.Println(ui.RenderError(fmt.Sprintf("文件 %v: %v", epubpath, err)))
					panic(fmt.Errorf("file %v  %v", epubpath, err))
				} else if err != nil && c.Skip {
					// 清除进度行
					if stats.Total > 1 {
						fmt.Print("\r" + strings.Repeat(" ", 120) + "\r")
					}
					fmt.Println(ui.FormatFilePath("文件", epubpath))
					fmt.Println(ui.RenderWarning(fmt.Sprintf("跳过: %v", err)))
					fmt.Println()
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
			if err != nil && !c.Skip {
				fmt.Println(ui.RenderError(fmt.Sprintf("文件 %v: %v", epubpath, err)))
				panic(fmt.Errorf("file %v  %v", epubpath, err))
			} else if err != nil && c.Skip {
				fmt.Println(ui.FormatFilePath("文件", epubpath))
				fmt.Println(ui.RenderWarning(fmt.Sprintf("跳过: %v", err)))
			}
		}

		// 打印统计信息
		printClnameStats(stats)
	},
}

func ValidateConfig(c *ClnameConfig) {
	if !util.Exists(c.Path) {
		fmt.Println(ui.RenderError("文件路径不存在，请检查"))
		os.Exit(1)
	}
	if util.IsFile(c.Path) && !strings.HasSuffix(c.Path, ".epub") {
		fmt.Println(ui.RenderError("文件格式不正确，必须是 .epub 文件"))
		os.Exit(1)
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
	clnameCmd.Flags().StringVarP(&c.Path, "path", "p", "./",
		"目录或者文件")
	clnameCmd.Flags().BoolVarP(&c.DoTry, "dotry", "t", false,
		"尝试运行")
	clnameCmd.Flags().BoolVarP(&c.Skip, "skip", "j", false,
		"跳过无法解析的书籍")
	clnameCmd.Flags().BoolVarP(&c.Debug, "debug", "d", false, "调试模式")
	clnameCmd.Flags().StringVar(&c.MoveCorruptedTo, "move-corrupted-to", "",
		"将损坏的文件移动到指定目录")
	clnameCmd.Flags().BoolVar(&c.DeleteCorrupted, "delete-corrupted", false,
		"删除损坏的文件")
	clnameCmd.Flags().BoolVar(&c.ForceDelete, "force-delete", false,
		"删除损坏的文件时不需要确认")
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
	DoTry           bool
	Debug           bool
	Skip            bool
	MoveCorruptedTo string // 损坏文件移动目标目录
	DeleteCorrupted bool   // 是否删除损坏文件
	ForceDelete     bool   // 删除时不需要确认
}
