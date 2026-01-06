package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jianyun8023/bookimporter/internal/application/clean"
	"github.com/jianyun8023/bookimporter/internal/domain/naming"
	"github.com/jianyun8023/bookimporter/internal/infrastructure/epub"
	"github.com/jianyun8023/bookimporter/internal/infrastructure/filesystem"
	"github.com/jianyun8023/bookimporter/pkg/ui"
	"github.com/spf13/cobra"
)

// cleanCmdConfig 用于绑定 flags
var cleanCmdConfig = &clean.Config{}

// clnameCmd used for download books from sanqiu.cc
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
		// 依赖注入
		repo := filesystem.NewRepository()
		cleaner := naming.NewCleaner()
		epubService := epub.NewService()
		svc := clean.NewService(repo, cleaner, epubService)

		// Validate config.
		validateCleanConfig(cleanCmdConfig, repo)

		// 打印头部
		fmt.Println(ui.RenderHeader("清理书籍标题", "移除标题中的无用描述符和标记"))
		fmt.Println()

		if err := svc.Run(cleanCmdConfig); err != nil {
			// application service 已经打印了错误信息，这里只需要退出
			os.Exit(1)
		}
	},
}

func validateCleanConfig(c *clean.Config, repo *filesystem.Repository) {
	// 验证路径是否存在
	if !repo.Exists(c.Path) {
		fmt.Println(ui.RenderError("文件路径不存在，请检查"))
		os.Exit(1)
	}

	// 验证文件格式
	if repo.IsFile(c.Path) && !strings.HasSuffix(c.Path, ".epub") {
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

func init() {
	// 基础参数
	clnameCmd.Flags().StringVarP(&cleanCmdConfig.Path, "path", "p", "./",
		"指定要处理的 EPUB 文件或目录路径")
	clnameCmd.Flags().BoolVarP(&cleanCmdConfig.Recursive, "recursive", "r", false,
		"递归搜索子目录中的所有 EPUB 文件")

	// 运行模式
	clnameCmd.Flags().BoolVarP(&cleanCmdConfig.DoTry, "dotry", "t", false,
		"预览模式，显示将要进行的修改但不实际执行")
	clnameCmd.Flags().BoolVarP(&cleanCmdConfig.IgnoreErrors, "ignore-errors", "i", false,
		"忽略错误，即使有失败也返回退出码 0")

	// 损坏文件处理（互斥选项）
	clnameCmd.Flags().StringVar(&cleanCmdConfig.MoveCorruptedTo, "move-corrupted-to", "",
		"将损坏的 EPUB 文件移动到指定目录（与 --delete-corrupted 互斥）")
	clnameCmd.Flags().BoolVar(&cleanCmdConfig.DeleteCorrupted, "delete-corrupted", false,
		"直接删除损坏的 EPUB 文件（与 --move-corrupted-to 互斥）")
	clnameCmd.Flags().BoolVar(&cleanCmdConfig.ForceDelete, "force-delete", false,
		"删除损坏文件时不需要用户确认（需配合 --delete-corrupted 使用）")

	// 调试选项
	clnameCmd.Flags().BoolVarP(&cleanCmdConfig.Debug, "debug", "d", false,
		"启用调试模式，显示详细的执行信息")
}
