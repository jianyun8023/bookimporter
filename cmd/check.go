package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jianyun8023/bookimporter/pkg/util"
	"github.com/spf13/cobra"
)

// CheckConfig 检测命令配置
type CheckConfig struct {
	Path       string // 文件或目录路径
	Recursive  bool   // 是否递归搜索
	OnlyErrors bool   // 只显示有问题的文件
	MoveTo     string // 移动损坏文件到指定目录
	Delete     bool   // 删除损坏的文件
	Force      bool   // 删除时不需要确认
	DoTry      bool   // 试运行模式
	Debug      bool   // 调试模式
}

var checkConfig = &CheckConfig{}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "检测 EPUB 文件完整性",
	Long: `检测 EPUB 文件是否损坏，包括 ZIP 结构、必需文件和元数据验证。
可以选择将损坏的文件移动到指定目录或删除。`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := validateCheckConfig(checkConfig); err != nil {
			fmt.Fprintf(os.Stderr, "配置错误: %v\n", err)
			os.Exit(1)
		}

		if err := runCheck(checkConfig); err != nil {
			fmt.Fprintf(os.Stderr, "检测失败: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	checkCmd.Flags().StringVarP(&checkConfig.Path, "path", "p", "",
		"文件或目录路径（必需）")
	checkCmd.Flags().BoolVarP(&checkConfig.Recursive, "recursive", "r", false,
		"递归搜索子目录")
	checkCmd.Flags().BoolVar(&checkConfig.OnlyErrors, "only-errors", false,
		"只显示有问题的文件")
	checkCmd.Flags().StringVar(&checkConfig.MoveTo, "move-to", "",
		"移动损坏文件到指定目录")
	checkCmd.Flags().BoolVar(&checkConfig.Delete, "delete", false,
		"删除损坏的文件")
	checkCmd.Flags().BoolVar(&checkConfig.Force, "force", false,
		"删除时不需要确认（与 --delete 配合）")
	checkCmd.Flags().BoolVar(&checkConfig.DoTry, "do-try", false,
		"试运行模式，不实际执行操作")
	checkCmd.Flags().BoolVarP(&checkConfig.Debug, "debug", "d", false,
		"调试模式")

	checkCmd.MarkFlagRequired("path")
}

// validateCheckConfig 验证配置
func validateCheckConfig(cfg *CheckConfig) error {
	if cfg.Path == "" {
		return fmt.Errorf("必须指定 --path 参数")
	}

	if !util.Exists(cfg.Path) {
		return fmt.Errorf("路径不存在: %s", cfg.Path)
	}

	// --move-to 和 --delete 互斥
	if cfg.MoveTo != "" && cfg.Delete {
		return fmt.Errorf("--move-to 和 --delete 不能同时使用")
	}

	// --force 只能与 --delete 一起使用
	if cfg.Force && !cfg.Delete {
		return fmt.Errorf("--force 只能与 --delete 一起使用")
	}

	return nil
}

// runCheck 执行检测
func runCheck(cfg *CheckConfig) error {
	var files []string

	// 收集要检测的文件
	if util.IsFile(cfg.Path) {
		if !strings.HasSuffix(strings.ToLower(cfg.Path), ".epub") {
			return fmt.Errorf("不是 EPUB 文件: %s", cfg.Path)
		}
		files = append(files, cfg.Path)
	} else {
		var err error
		files, err = collectEpubFiles(cfg.Path, cfg.Recursive)
		if err != nil {
			return fmt.Errorf("收集文件失败: %w", err)
		}
	}

	if len(files) == 0 {
		fmt.Println("未找到 EPUB 文件")
		return nil
	}

	// 统计信息
	stats := &CheckStats{
		Total:   len(files),
		Passed:  0,
		Failed:  0,
		Handled: 0,
	}

	// 检测每个文件
	for _, file := range files {
		if err := checkSingleFile(file, cfg, stats); err != nil {
			if cfg.Debug {
				fmt.Fprintf(os.Stderr, "处理文件 %s 时发生错误: %v\n", file, err)
			}
		}
	}

	// 打印统计信息
	printStats(stats)

	return nil
}

// CheckStats 检测统计
type CheckStats struct {
	Total   int // 总文件数
	Passed  int // 通过数
	Failed  int // 失败数
	Handled int // 已处理数（移动或删除）
}

// collectEpubFiles 收集 EPUB 文件
func collectEpubFiles(dir string, recursive bool) ([]string, error) {
	var files []string

	if recursive {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".epub") {
				files = append(files, path)
			}
			return nil
		})
		return files, err
	}

	// 非递归模式，只查找当前目录
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".epub") {
			files = append(files, filepath.Join(dir, entry.Name()))
		}
	}

	return files, nil
}

// checkSingleFile 检测单个文件
func checkSingleFile(file string, cfg *CheckConfig, stats *CheckStats) error {
	err := util.ValidateEpubFile(file)

	if err == nil {
		// 文件正常
		stats.Passed++
		if !cfg.OnlyErrors {
			fmt.Printf("检查: %s\n", file)
			fmt.Println("✓ 通过")
			fmt.Println()
		}
		return nil
	}

	// 文件有问题
	stats.Failed++
	fmt.Printf("检查: %s\n", file)
	fmt.Printf("✗ 失败: %v\n", err)

	// 处理损坏的文件
	if cfg.MoveTo != "" {
		if err := handleMoveFile(file, cfg.MoveTo, cfg.DoTry); err != nil {
			fmt.Printf("  → 移动失败: %v\n", err)
		} else {
			stats.Handled++
		}
	} else if cfg.Delete {
		if err := handleDeleteFile(file, cfg.Force, cfg.DoTry); err != nil {
			fmt.Printf("  → 删除失败: %v\n", err)
		} else {
			stats.Handled++
		}
	}

	fmt.Println()
	return nil
}

// handleMoveFile 处理移动文件
func handleMoveFile(srcPath, dstDir string, doTry bool) error {
	if doTry {
		// 试运行模式，只显示将要执行的操作
		fileName := filepath.Base(srcPath)
		expectedPath := filepath.Join(dstDir, fileName)
		fmt.Printf("  → [试运行] 将移动到: %s\n", expectedPath)
		return nil
	}

	newPath, err := util.MoveFileWithConflictHandling(srcPath, dstDir)
	if err != nil {
		return err
	}

	fmt.Printf("  → 已移动到: %s\n", newPath)
	return nil
}

// handleDeleteFile 处理删除文件
func handleDeleteFile(filePath string, force, doTry bool) error {
	if doTry {
		// 试运行模式，只显示将要执行的操作
		fmt.Printf("  → [试运行] 将删除\n")
		return nil
	}

	needConfirm := !force
	if err := util.SafeDeleteFile(filePath, needConfirm); err != nil {
		return err
	}

	fmt.Printf("  → 已删除\n")
	return nil
}

// printStats 打印统计信息
func printStats(stats *CheckStats) {
	fmt.Println(strings.Repeat("-", 50))
	if stats.Handled > 0 {
		fmt.Printf("统计: 总计 %d 个文件，通过 %d 个，失败 %d 个，已处理 %d 个\n",
			stats.Total, stats.Passed, stats.Failed, stats.Handled)
	} else {
		fmt.Printf("统计: 总计 %d 个文件，通过 %d 个，失败 %d 个\n",
			stats.Total, stats.Passed, stats.Failed)
	}
}
