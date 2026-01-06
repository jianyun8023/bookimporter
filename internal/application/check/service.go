package check

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jianyun8023/bookimporter/internal/domain/book"
	"github.com/jianyun8023/bookimporter/internal/infrastructure/epub"
	"github.com/jianyun8023/bookimporter/pkg/ui"
)

// CheckStats 统计检测结果
type CheckStats struct {
	Total   int
	Passed  int
	Failed  int
	Handled int
}

// Config 检测配置
type Config struct {
	Path       string
	Recursive  bool
	OnlyErrors bool
	MoveTo     string
	Delete     bool
	Force      bool
	DoTry      bool
	Debug      bool
}

// Service 负责 EPUB 检测业务流程
type Service struct {
	repo        book.Repository
	epubService *epub.Service
}

// NewService 创建检测服务
func NewService(repo book.Repository, epubService *epub.Service) *Service {
	return &Service{
		repo:        repo,
		epubService: epubService,
	}
}

// Run 执行检测流程
func (s *Service) Run(cfg *Config) error {
	// 1. 查找所有待检测文件
	var files []string
	if s.repo.IsFile(cfg.Path) {
		files = []string{cfg.Path}
	} else {
		// 这里我们复用 book.Repository 来查找文件
		// 实际上 Repository 返回的是 Book 对象列表，我们只需要路径
		books, err := s.repo.FindAll(cfg.Path, cfg.Recursive)
		if err != nil {
			return err
		}
		for _, b := range books {
			files = append(files, b.Path)
		}
	}

	if len(files) == 0 {
		fmt.Println(ui.RenderWarning("未找到 EPUB 文件"))
		return nil
	}

	fmt.Println(ui.RenderInfo(fmt.Sprintf("找到 %d 个 EPUB 文件", len(files))))
	fmt.Println()

	stats := &CheckStats{Total: len(files)}
	progress := ui.NewCompactProgressTracker(len(files))
	progress.SetShowMessage(true)

	for i, file := range files {
		progress.SetMessage(filepath.Base(file))
		if len(files) > 1 {
			fmt.Printf("\r%s", progress.RenderCompact())
		}

		err := s.checkSingleFile(file, cfg, stats)

		if err == nil {
			progress.IncrementSuccess()
		} else {
			progress.IncrementFailure()
		}

		if len(files) > 1 && i < len(files)-1 {
			fmt.Print("\r" + "                                                                                                                        " + "\r")
		}
	}

	if len(files) > 1 {
		fmt.Print("\r" + "                                                                                                                        " + "\r")
		fmt.Println(progress.RenderWithStats())
		fmt.Println()
	}

	s.printStats(stats)
	return nil
}

func (s *Service) checkSingleFile(file string, cfg *Config, stats *CheckStats) error {
	err := s.epubService.Validate(file)

	if err == nil {
		stats.Passed++
		if !cfg.OnlyErrors {
			fmt.Println(ui.FormatFilePath("检查", file))
			fmt.Println(ui.RenderSuccess("通过"))
			fmt.Println()
		}
		return nil
	}

	stats.Failed++
	fmt.Println(ui.FormatFilePath("检查", file))
	fmt.Println(ui.RenderError(fmt.Sprintf("失败: %v", err)))

	// 处理损坏文件
	if cfg.MoveTo != "" {
		if cfg.DoTry {
			fmt.Println(ui.RenderInfo(fmt.Sprintf("[试运行] 将移动到: %s", filepath.Join(cfg.MoveTo, filepath.Base(file)))))
		} else {
			if _, err := s.repo.Move(file, cfg.MoveTo); err != nil {
				fmt.Println(ui.RenderError(fmt.Sprintf("移动失败: %v", err)))
			} else {
				stats.Handled++
				fmt.Println(ui.RenderInfo("已移动损坏文件"))
			}
		}
	} else if cfg.Delete {
		if cfg.DoTry {
			fmt.Println(ui.RenderInfo("[试运行] 将删除"))
		} else {
			// 如果需要确认，这里可能需要 UI 交互，或者假设 force 已经处理
			// DDD 中 Service 层不应直接处理 stdio 交互，最好是由 UI 层传入确认回调，或者在 Config 中明确
			// 为了简化，我们假设 force 参数已经在上层处理过确认逻辑（在 CheckConfig 解析时？）
			// 原 cmd/check.go 里的逻辑是：--delete 需要确认，除非 --force。
			// 这里我们为了保持行为，可能需要在 repo 中实现，或者在这里实现简单的确认
			// 但最好的方式是：Service 只管执行。Confirm 的责任在于上层。
			// 鉴于这是一个 CLI 工具，我们在 Repository 实现 Delete 时是直接删除。
			// 这里的确认逻辑在 Service 内部似乎也不太合适。
			// 让我们简单点：如果 cfg.Force 为 false 且 cfg.Delete 为 true，
			// 我们应该在 Service 里调用 UI 确认吗？
			// 实际上，原代码是在 SafeDeleteFile 工具函数里做确认。
			// 为了不破坏现有体验，我们在 Infrastructure 层实现 Delete 时并没有带确认。
			// 我们可以在这里简单实现一个确认逻辑，或者假设调用者已经确认。

			// 考虑到 CMD 层调用 Service，Service 处理业务逻辑。
			// 删除是一个业务动作。确认是一个 UI 动作。
			// 理想情况下：Service 定义一个 Delegate 接口用于 Confirm，CMD 实现它。
			// 简单起见，我们把 confirm 逻辑直接内联在这里，使用 pkg/ui

			needConfirm := !cfg.Force
			confirmed := true
			if needConfirm {
				fmt.Println()
				warningMsg := ui.RenderMessageBox(
					"删除确认",
					fmt.Sprintf("即将删除文件:\n%s\n\n此操作无法撤销！", ui.RenderPath(file)),
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
					fmt.Println(ui.RenderError(fmt.Sprintf("删除失败: %v", err)))
				} else {
					stats.Handled++
					fmt.Println(ui.RenderInfo("已删除"))
				}
			}
		}
	}

	fmt.Println()
	return nil
}

func (s *Service) printStats(stats *CheckStats) {
	fmt.Println(ui.RenderSeparator(60))
	fmt.Println()

	tableConfig := ui.NewTableConfig()
	tableConfig.Headers = []string{"  状态  ", " 数量 ", " 百分比 "}
	tableConfig.BorderStyle = "rounded"
	tableConfig.AlignRight = []int{1, 2}

	var rows [][]string
	if stats.Passed > 0 {
		percentage := float64(stats.Passed) / float64(stats.Total) * 100
		rows = append(rows, []string{
			ui.IconSuccess + " 通过  ",
			fmt.Sprintf(" %d ", stats.Passed),
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
	if stats.Handled > 0 {
		rows = append(rows, []string{
			ui.IconInfo + " 已处理 ",
			fmt.Sprintf(" %d ", stats.Handled),
			fmt.Sprintf(" - "),
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

	if stats.Failed == 0 {
		fmt.Println(ui.RenderSuccess(fmt.Sprintf("✨ 所有 %d 个文件检测通过！", stats.Total)))
	} else {
		fmt.Println(ui.RenderWarning(fmt.Sprintf("⚠️  发现 %d 个问题文件", stats.Failed)))
		if stats.Handled > 0 {
			fmt.Println(ui.RenderInfo(fmt.Sprintf("已处理 %d 个问题文件", stats.Handled)))
		}
	}
}
