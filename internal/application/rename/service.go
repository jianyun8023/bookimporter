package rename

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jianyun8023/bookimporter/internal/domain/book"
	"github.com/jianyun8023/bookimporter/pkg/ui"
)

type Config struct {
	Debug      bool
	DoTry      bool
	Formats    []string
	Recursive  bool
	SourceDir  string
	OutputDir  string
	Template   string
	StartIndex int
}

type Service struct {
	repo book.Repository
}

func NewService(repo book.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Run(cfg *Config) error {
	files, err := s.findFiles(cfg.SourceDir, cfg.Formats, cfg.Recursive)
	if err != nil {
		return fmt.Errorf("find files failed: %w", err)
	}

	if len(files) == 0 {
		fmt.Println(ui.RenderWarning("未找到匹配的文件"))
		return nil
	}

	fmt.Println(ui.RenderInfo(fmt.Sprintf("找到 %d 个文件", len(files))))
	fmt.Println()

	if cfg.DoTry {
		s.preview(files, cfg)
	} else {
		if err := s.execute(files, cfg); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) findFiles(dir string, formats []string, recursive bool) ([]string, error) {
	var files []string

	matches, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return nil, err
	}

	for _, match := range matches {
		included := false
		// 检查格式
		for _, format := range formats {
			if format == "*" || strings.HasSuffix(strings.ToLower(match), "."+strings.ToLower(format)) || strings.Contains(match, format) {
				if strings.Contains(match, format) {
					included = true
					break
				}
			}
		}

		info, err := os.Stat(match)
		if err != nil {
			return nil, err
		}

		if !info.IsDir() && included {
			files = append(files, match)
		} else if recursive && info.IsDir() {
			subFiles, err := s.findFiles(match, formats, true)
			if err != nil {
				return nil, err
			}
			files = append(files, subFiles...)
		}
	}

	return files, nil
}

func (s *Service) preview(files []string, cfg *Config) {
	fmt.Println(ui.RenderTitle("📋 重命名预览"))
	fmt.Println()

	tableConfig := ui.NewTableConfig()
	tableConfig.Headers = []string{" # ", " 原文件名 ", " → ", " 新文件名 "}
	tableConfig.BorderStyle = "rounded"
	tableConfig.CompactMode = false

	var rows [][]string
	for i, file := range files {
		newName := s.buildNewName(cfg.Template, cfg.StartIndex+i, file)
		var outputPath string
		if cfg.OutputDir != "" {
			outputPath = filepath.Join(cfg.OutputDir, newName)
		} else {
			outputPath = filepath.Join(filepath.Dir(file), newName)
		}

		// UI 渲染逻辑保持
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

		if i >= 19 && len(files) > 20 {
			rows = append(rows, []string{"...", "...", "...", "..."})
			break
		}
	}

	tableConfig.Rows = rows
	table := ui.NewTable(tableConfig)
	fmt.Println(table.Render())
	fmt.Println()

	// 统计预览
	s.printSummary(len(files), 0, cfg)
}

func (s *Service) execute(files []string, cfg *Config) error {
	var (
		renamedCount int
		movedCount   int
	)

	progress := ui.NewCompactProgressTracker(len(files))
	progress.SetShowMessage(true)

	for i, file := range files {
		newName := s.buildNewName(cfg.Template, cfg.StartIndex+i, file)

		progress.SetMessage(filepath.Base(file))
		fmt.Printf("\r%s", progress.RenderCompact())

		if cfg.OutputDir != "" {
			targetPath := filepath.Join(cfg.OutputDir, newName)

			// We can use Rename for moving too, as os.Rename supports paths.
			_, err := s.repo.Rename(file, targetPath)
			if err != nil {
				// clear progress
				fmt.Print("\r" + "                                                                                                                        " + "\r")
				fmt.Println(ui.RenderError(fmt.Sprintf("重命名失败: %s", err)))
				return err // Exit on error as per original logic
			}
			movedCount++
		} else {
			// Rename in place
			targetPath := filepath.Join(filepath.Dir(file), newName)
			_, err := s.repo.Rename(file, targetPath)
			if err != nil {
				fmt.Print("\r" + "                                                                                                                        " + "\r")
				fmt.Println(ui.RenderError(fmt.Sprintf("重命名失败: %s", err)))
				return err
			}
			renamedCount++
		}

		progress.IncrementSuccess()
		fmt.Print("\r" + "                                                                                                                        " + "\r")
		// fmt.Println(ui.FormatRenamePreview(file, ...)) // Optional: print each line
	}

	fmt.Println(progress.RenderWithStats())
	fmt.Println()

	s.printSummary(len(files), movedCount+renamedCount, cfg)
	return nil
}

func (s *Service) buildNewName(template string, index int, file string) string {
	ext := filepath.Ext(file)
	newName := strings.Replace(template, "@n", strconv.Itoa(index), -1)
	newName += ext
	return newName
}

func (s *Service) printSummary(total int, processed int, cfg *Config) {
	fmt.Println()
	fmt.Println(ui.RenderSeparator(60))
	fmt.Println()

	tableConfig := ui.NewTableConfig()
	tableConfig.Headers = []string{"  项目  ", " 值 "}
	tableConfig.BorderStyle = "rounded"
	tableConfig.AlignRight = []int{1}

	var rows [][]string
	rows = append(rows, []string{
		" 文件总数 ",
		fmt.Sprintf(" %d ", total),
	})

	if cfg.OutputDir != "" {
		rows = append(rows, []string{
			" 目标目录 ",
			fmt.Sprintf(" %s ", cfg.OutputDir),
		})
	}

	if cfg.DoTry {
		rows = append(rows, []string{" 模式 ", " 预览模式 "})
	} else {
		rows = append(rows, []string{" 已处理 ", fmt.Sprintf(" %d ", processed)})
	}

	tableConfig.Rows = rows
	table := ui.NewTable(tableConfig)
	fmt.Println(table.Render())
	fmt.Println()

	if cfg.OutputDir != "" {
		if cfg.DoTry {
			fmt.Println(ui.RenderInfo(fmt.Sprintf("📝 [试运行] 将移动 %d 个文件到: %s", total, cfg.OutputDir)))
		} else {
			fmt.Println(ui.RenderSuccess(fmt.Sprintf("✨ 成功移动 %d 个文件到: %s", processed, cfg.OutputDir)))
		}
	} else {
		if cfg.DoTry {
			fmt.Println(ui.RenderInfo(fmt.Sprintf("📝 [试运行] 将重命名 %d 个文件", total)))
		} else {
			fmt.Println(ui.RenderSuccess(fmt.Sprintf("✨ 成功重命名 %d 个文件", processed)))
		}
	}
}
