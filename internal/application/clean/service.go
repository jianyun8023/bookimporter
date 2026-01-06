package clean

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jianyun8023/bookimporter/internal/domain/book"
	"github.com/jianyun8023/bookimporter/internal/domain/naming"
	"github.com/jianyun8023/bookimporter/internal/infrastructure/epub"
	"github.com/jianyun8023/bookimporter/pkg/ui"
)

type Config struct {
	Path            string
	Recursive       bool
	DoTry           bool
	Debug           bool
	IgnoreErrors    bool
	MoveCorruptedTo string
	DeleteCorrupted bool
	ForceDelete     bool
}

type CleanStats struct {
	Total   int
	Updated int
	Skipped int
	Failed  int
}

type Service struct {
	repo        book.Repository
	cleaner     *naming.Cleaner
	epubService *epub.Service
}

func NewService(repo book.Repository, cleaner *naming.Cleaner, epubService *epub.Service) *Service {
	return &Service{
		repo:        repo,
		cleaner:     cleaner,
		epubService: epubService,
	}
}

func (s *Service) Run(cfg *Config) error {
	stats := &CleanStats{}

	// 查找文件
	var files []string
	if s.repo.IsFile(cfg.Path) {
		files = []string{cfg.Path}
	} else {
		books, err := s.repo.FindAll(cfg.Path, cfg.Recursive)
		if err != nil {
			return fmt.Errorf("scan failed: %w", err)
		}
		for _, b := range books {
			files = append(files, b.Path)
		}
	}
	stats.Total = len(files)

	if stats.Total == 0 {
		fmt.Println(ui.RenderWarning("未找到 EPUB 文件"))
		return nil
	}

	fmt.Println(ui.RenderInfo(fmt.Sprintf("找到 %d 个 EPUB 文件", stats.Total)))
	fmt.Println()

	progress := ui.NewCompactProgressTracker(stats.Total)
	progress.SetShowMessage(true)

	for i, file := range files {
		progress.SetMessage(filepath.Base(file))
		if stats.Total > 1 {
			fmt.Printf("\r%s", progress.RenderCompact())
		}

		err := s.processFile(file, cfg, stats, progress)
		if err != nil {
			if stats.Total > 1 {
				fmt.Print("\r" + "                                                                                                                        " + "\r")
			}
			fmt.Println(ui.FormatFilePath("文件", file))
			fmt.Println(ui.RenderWarning(fmt.Sprintf("跳过: %v", err)))
			fmt.Println()
			stats.Failed++
			progress.IncrementFailure()
		}

		if stats.Total > 1 && i < stats.Total-1 {
			fmt.Print("\r" + "                                                                                                                        " + "\r")
		}
	}

	if stats.Total > 1 {
		fmt.Print("\r" + "                                                                                                                        " + "\r")
		fmt.Println(progress.RenderWithStats())
		fmt.Println()
	}

	s.printStats(stats)

	if stats.Failed > 0 && !cfg.IgnoreErrors {
		// 返回错误以便 exit code 为 1
		return fmt.Errorf("process completed with %d failures", stats.Failed)
	}
	return nil
}

func (s *Service) processFile(file string, cfg *Config, stats *CleanStats, progress *ui.ProgressTracker) error {
	// 1. 验证完整性
	if err := s.epubService.Validate(file); err != nil {
		if cfg.Debug {
			fmt.Println(ui.RenderError(fmt.Sprintf("EPUB 文件检测失败: %v", err)))
		}

		// 处理损坏文件
		if cfg.MoveCorruptedTo != "" {
			if cfg.DoTry {
				fmt.Println(ui.RenderInfo(fmt.Sprintf("[试运行] 将移动损坏文件到: %s", cfg.MoveCorruptedTo)))
			} else {
				if _, moveErr := s.repo.Move(file, cfg.MoveCorruptedTo); moveErr != nil {
					fmt.Println(ui.RenderError(fmt.Sprintf("移动损坏文件失败: %v", moveErr)))
				} else {
					fmt.Println(ui.RenderInfo("已移动损坏文件"))
				}
			}
		} else if cfg.DeleteCorrupted {
			if cfg.DoTry {
				fmt.Println(ui.RenderInfo("[试运行] 将删除损坏文件"))
			} else {
				// 需要确认删除
				needConfirm := !cfg.ForceDelete
				confirmed := true

				if needConfirm {
					fmt.Println()
					warningMsg := ui.RenderMessageBox(
						"删除确认",
						fmt.Sprintf("即将删除损坏文件:\n%s\n\n此操作无法撤销！", ui.RenderPath(file)),
						"warning",
					)
					fmt.Println(warningMsg)
					fmt.Println()
					fmt.Print(ui.StyleWarning.Render("确认删除? (y/N): "))

					reader := bufio.NewReader(os.Stdin)
					input, err := reader.ReadString('\n')
					if err != nil {
						fmt.Println(ui.RenderError(fmt.Sprintf("读取输入失败: %v", err)))
						confirmed = false
					} else {
						input = strings.TrimSpace(strings.ToLower(input))
						if input != "y" && input != "yes" {
							fmt.Println(ui.RenderInfo("已取消删除操作"))
							confirmed = false
						}
					}
				}

				if confirmed {
					if err := s.repo.Delete(file); err != nil {
						fmt.Println(ui.RenderError(fmt.Sprintf("删除损坏文件失败: %v", err)))
					} else {
						fmt.Println(ui.RenderInfo("已删除损坏文件"))
					}
				}
			}
		}
		return err
	}

	// 2. 获取标题
	title, err := s.epubService.GetTitle(file)
	if err != nil {
		return err
	}

	// 3. 清洗标题
	newTitle := s.cleaner.Clean(title)
	if title == newTitle {
		stats.Skipped++
		progress.IncrementSkipped()
		return nil
	}

	// 4. 更新元数据
	fmt.Println(ui.FormatFilePath("路径", file))
	fmt.Println(ui.FormatFileOperation("标题", title, newTitle))

	if cfg.DoTry {
		fmt.Println(ui.RenderInfo("[试运行] 将更新标题"))
		fmt.Println()
		stats.Skipped++ // 试运行不算更新成功，算跳过？或者保持原逻辑
		// 原逻辑是 Skipped++
		progress.IncrementSkipped()
		return nil
	}

	if err := s.repo.UpdateMetadata(file, newTitle); err != nil {
		return err
	}

	fmt.Println(ui.RenderSuccess("已更新"))
	fmt.Println()
	stats.Updated++
	progress.IncrementSuccess()
	return nil
}

func (s *Service) printStats(stats *CleanStats) {
	fmt.Println()
	fmt.Println(ui.RenderSeparator(60))
	fmt.Println()

	tableConfig := ui.NewTableConfig()
	tableConfig.Headers = []string{"  状态  ", " 数量 ", " 百分比 "}
	tableConfig.BorderStyle = "rounded"
	tableConfig.AlignRight = []int{1, 2}

	var rows [][]string
	if stats.Updated > 0 {
		percentage := float64(stats.Updated) / float64(stats.Total) * 100
		rows = append(rows, []string{
			ui.IconSuccess + " 已更新 ",
			fmt.Sprintf(" %d ", stats.Updated),
			fmt.Sprintf(" %.1f%% ", percentage),
		})
	}
	if stats.Skipped > 0 {
		percentage := float64(stats.Skipped) / float64(stats.Total) * 100
		rows = append(rows, []string{
			ui.IconSkip + " 跳过  ",
			fmt.Sprintf(" %d ", stats.Skipped),
			fmt.Sprintf(" %.1f%% ", percentage),
		})
	}
	if stats.Failed > 0 {
		percentage := float64(stats.Failed) / float64(stats.Total) * 100
		rows = append(rows, []string{
			ui.IconError + " 失败  ",
			fmt.Sprintf(" %d ", stats.Failed),
			fmt.Sprintf(" %.1f%% ", percentage),
		})
	}
	rows = append(rows, []string{
		"  总计  ",
		fmt.Sprintf(" %d ", stats.Total),
		" 100% ",
	})

	tableConfig.Rows = rows
	table := ui.NewTable(tableConfig)
	fmt.Println(table.Render())
	fmt.Println()

	if stats.Updated == 0 && stats.Failed == 0 {
		fmt.Println(ui.RenderInfo("✨ 所有文件标题都已是最佳状态"))
	} else if stats.Updated > 0 {
		fmt.Println(ui.RenderSuccess(fmt.Sprintf("✨ 成功更新 %d 个文件的标题", stats.Updated)))
	}
	if stats.Failed > 0 {
		fmt.Println(ui.RenderWarning(fmt.Sprintf("⚠️  %d 个文件处理失败", stats.Failed)))
	}
}
