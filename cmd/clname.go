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

			// 创建进度跟踪器
			progress := ui.NewProgressTracker(stats.Total)

			for i, val := range m {
				epubpath := val
				progress.SetMessage(fmt.Sprintf("正在处理: %s", filepath.Base(epubpath)))

				// 显示进度
				if stats.Total > 1 {
					fmt.Printf("\r%s", progress.RenderSimple())
				}

				err := ParseEpub(epubpath, c, stats)
				if err != nil && !c.Skip {
					// 清除进度行
					if stats.Total > 1 {
						fmt.Print("\r" + strings.Repeat(" ", 80) + "\r")
					}
					fmt.Println(ui.RenderError(fmt.Sprintf("文件 %v: %v", epubpath, err)))
					panic(fmt.Errorf("file %v  %v", epubpath, err))
				} else if err != nil && c.Skip {
					// 清除进度行
					if stats.Total > 1 {
						fmt.Print("\r" + strings.Repeat(" ", 80) + "\r")
					}
					fmt.Println(ui.FormatFilePath("文件", epubpath))
					fmt.Println(ui.RenderWarning(fmt.Sprintf("跳过: %v", err)))
					fmt.Println()
					stats.Failed++
				}

				progress.Increment()

				// 清除进度行，为文件详情腾出空间
				if stats.Total > 1 && i < stats.Total-1 {
					fmt.Print("\r" + strings.Repeat(" ", 80) + "\r")
				}
			}

			// 清除最后的进度行
			if stats.Total > 1 {
				fmt.Print("\r" + strings.Repeat(" ", 80) + "\r")
			}

		} else {
			stats.Total = 1
			epubpath := c.Path
			err := ParseEpub(epubpath, c, stats)
			if err != nil && !c.Skip {
				fmt.Println(ui.RenderError(fmt.Sprintf("文件 %v: %v", epubpath, err)))
				panic(fmt.Errorf("file %v  %v", epubpath, err))
			} else if err != nil && c.Skip {
				fmt.Println(ui.FormatFilePath("文件", epubpath))
				fmt.Println(ui.RenderWarning(fmt.Sprintf("跳过: %v", err)))
				stats.Failed++
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
	fmt.Println(ui.RenderSeparator(50))
	fmt.Println()

	// 准备统计数据
	statsMap := map[string]int{
		"total": stats.Total,
	}

	if stats.Updated > 0 {
		statsMap["updated"] = stats.Updated
	}

	if stats.Skipped > 0 {
		statsMap["skipped"] = stats.Skipped
	}

	if stats.Failed > 0 {
		statsMap["failed"] = stats.Failed
	}

	// 渲染统计表格
	fmt.Println(ui.RenderStatsSummary(statsMap))
	fmt.Println()

	// 添加总结信息
	if stats.Updated == 0 && stats.Failed == 0 {
		fmt.Println(ui.RenderInfo("所有文件标题都已是最佳状态"))
	} else if stats.Updated > 0 {
		fmt.Println(ui.RenderSuccess(fmt.Sprintf("成功更新 %d 个文件的标题", stats.Updated)))
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

func ParseEpub(file string, c *ClnameConfig, stats *ClnameStats) error {
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
		return nil
	}

	// 美化输出
	fmt.Println(ui.FormatFilePath("路径", file))
	fmt.Println(ui.FormatFileOperation("标题", title, newTitle))

	if c.DoTry {
		fmt.Println(ui.RenderInfo("[试运行] 将更新标题"))
		fmt.Println()
		stats.Skipped++
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
