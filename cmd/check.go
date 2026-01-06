package cmd

import (
	"fmt"
	"os"

	"github.com/jianyun8023/bookimporter/internal/application/check"
	"github.com/jianyun8023/bookimporter/internal/infrastructure/epub"
	"github.com/jianyun8023/bookimporter/internal/infrastructure/filesystem"
	"github.com/jianyun8023/bookimporter/pkg/ui"
	"github.com/spf13/cobra"
)

// CheckConfig 检测命令配置
// 注意：这里的结构体是为了绑定 cobra flags， application 层有自己的 Config
// 我们可以直接使用 application 层的 Config，或者在这里做一个转换。
// 为了代码清晰，我们直接实例化 application.Config。
var checkCmdConfig = &check.Config{}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "检测 EPUB 文件完整性",
	Long: `检测 EPUB 文件是否损坏，包括 ZIP 结构、必需文件和元数据验证。
可以选择将损坏的文件移动到指定目录或删除。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 依赖注入
		repo := filesystem.NewRepository()
		epubService := epub.NewService()
		svc := check.NewService(repo, epubService)

		if err := validateCheckConfig(checkCmdConfig, repo); err != nil {
			fmt.Fprintf(os.Stderr, "配置错误: %v\n", err)
			os.Exit(1)
		}

		// 打印头部
		fmt.Println(ui.RenderHeader("EPUB 文件检测", "检查文件完整性、ZIP 结构和元数据"))
		fmt.Println()

		if err := svc.Run(checkCmdConfig); err != nil {
			fmt.Fprintf(os.Stderr, "检测失败: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	checkCmd.Flags().StringVarP(&checkCmdConfig.Path, "path", "p", "",
		"文件或目录路径（必需）")
	checkCmd.Flags().BoolVarP(&checkCmdConfig.Recursive, "recursive", "r", false,
		"递归搜索子目录")
	checkCmd.Flags().BoolVar(&checkCmdConfig.OnlyErrors, "only-errors", false,
		"只显示有问题的文件")
	checkCmd.Flags().StringVar(&checkCmdConfig.MoveTo, "move-to", "",
		"移动损坏文件到指定目录")
	checkCmd.Flags().BoolVar(&checkCmdConfig.Delete, "delete", false,
		"删除损坏的文件")
	checkCmd.Flags().BoolVar(&checkCmdConfig.Force, "force", false,
		"删除时不需要确认（与 --delete 配合）")
	checkCmd.Flags().BoolVar(&checkCmdConfig.DoTry, "do-try", false,
		"试运行模式，不实际执行操作")
	checkCmd.Flags().BoolVarP(&checkCmdConfig.Debug, "debug", "d", false,
		"调试模式")

	checkCmd.MarkFlagRequired("path")
}

// validateCheckConfig 验证配置
// 这里需要 repo 来检查文件是否存在，或者使用 os.Stat
func validateCheckConfig(cfg *check.Config, repo *filesystem.Repository) error {
	if cfg.Path == "" {
		return fmt.Errorf("必须指定 --path 参数")
	}

	if !repo.Exists(cfg.Path) {
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
