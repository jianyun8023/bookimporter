package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jianyun8023/bookimporter/internal/application/rename"
	"github.com/jianyun8023/bookimporter/internal/infrastructure/filesystem"
	"github.com/jianyun8023/bookimporter/pkg/ui"
	"github.com/spf13/cobra"
)

var renameCmd = &cobra.Command{
	Use:   "rename <目录路径>",
	Short: "按自定义模板批量重命名或移动文件",
	Long: `按自定义模板批量重命名或移动文件

rename 命令可以根据指定的模板批量重命名文件，支持序列号自动编号。
常用于整理大量文件，使其具有统一的命名规则。

核心功能：
  • 支持自定义文件名模板（使用 @n 作为序号占位符）
  • 支持多种文件格式过滤
  • 支持递归搜索子目录
  • 支持移动文件到指定目录
  • 提供预览模式，查看重命名结果
  • 自动保留原始文件扩展名

序号占位符：
  @n  - 会被替换为实际的序列号
  
重要说明：
  • 模板中必须包含 @n 占位符
  • 文件扩展名会自动添加，无需在模板中指定
  • 序列号默认从 1 开始，可通过 --start-num 自定义
  • 使用 --do-try 可以先预览结果，确认无误后再执行`,
	Example: `  # 基础用法：重命名当前目录下的 txt 文件
  bookimporter rename . -f txt -t "book-@n"
  结果: file1.txt → book-1.txt, file2.txt → book-2.txt

  # 递归搜索子目录中的 EPUB 文件并重命名
  bookimporter rename /path/to/books -f epub -t "novel-@n" -r
  
  # 重命名多种格式的文件
  bookimporter rename . -f epub -f pdf -f mobi -t "ebook-@n"
  
  # 重命名并移动到新目录（整理文件）
  bookimporter rename /source -f jpg -t "photo-@n" -o /photos
  
  # 从指定序号开始
  bookimporter rename . -f txt -t "doc-@n" --start-num 100
  结果: 从 doc-100.txt 开始编号
  
  # 预览模式：先查看将要进行的操作
  bookimporter rename . -f pdf -t "file-@n" --do-try`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			fmt.Println(ui.RenderError("错误: 缺少必需的目录路径参数"))
			fmt.Println()
			fmt.Println(ui.RenderInfo("用法示例:"))
			fmt.Println("  bookimporter rename . -f txt -t \"book-@n\"")
			fmt.Println("  bookimporter rename /path/to/dir -f epub -t \"novel-@n\" -r")
			fmt.Println()
			fmt.Println(ui.RenderInfo("使用 'bookimporter rename --help' 查看详细帮助"))
			os.Exit(1)
		}
		array, err := cmd.Flags().GetStringArray("format")
		if err != nil {
			panic(err)
		}
		config := &rename.Config{
			Debug:      cmd.Flag("debug").Value.String() == "true",
			DoTry:      cmd.Flag("do-try").Value.String() == "true",
			Formats:    array,
			Recursive:  cmd.Flag("recursive").Value.String() == "true",
			SourceDir:  args[0],
			OutputDir:  cmd.Flag("output").Value.String(),
			Template:   cmd.Flag("template").Value.String(),
			StartIndex: parseIntFlag(cmd, "start-num"),
		}

		validateRenameConfig(config)

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

		repo := filesystem.NewRepository()
		svc := rename.NewService(repo)
		if err := svc.Run(config); err != nil {
			os.Exit(1)
		}
	},
}

func validateRenameConfig(config *rename.Config) {
	// 检查模板是否包含序号占位符
	if !strings.Contains(config.Template, "@n") {
		fmt.Println(ui.RenderError(fmt.Sprintf("错误: 模板 '%s' 中缺少序号占位符 @n", config.Template)))
		fmt.Println()
		fmt.Println(ui.RenderInfo("模板必须包含 @n 作为序号占位符，例如:"))
		fmt.Println("  ✓ 正确: -t \"book-@n\"    → book-1.epub, book-2.epub")
		fmt.Println("  ✓ 正确: -t \"file_@n\"    → file_1.txt, file_2.txt")
		fmt.Println("  ✓ 正确: -t \"doc-@n-new\" → doc-1-new.pdf, doc-2-new.pdf")
		fmt.Println("  ✗ 错误: -t \"book\"       (没有 @n)")
		os.Exit(1)
	}

	// 检查源目录是否存在
	if _, err := os.Stat(config.SourceDir); os.IsNotExist(err) {
		fmt.Println(ui.RenderError(fmt.Sprintf("错误: 目录不存在: %s", config.SourceDir)))
		os.Exit(1)
	}

	// 如果指定了输出目录，确保其存在或可创建
	if config.OutputDir != "" {
		if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
			fmt.Println(ui.RenderError(fmt.Sprintf("错误: 无法创建输出目录: %s", err)))
			os.Exit(1)
		}
	}
}

func init() {
	renameCmd.Flags().StringArrayP("format", "f", []string{"*"},
		"指定要处理的文件格式（如 'txt', 'epub'），可多次使用以匹配多种格式")
	renameCmd.Flags().StringP("template", "t", "file-@n",
		"文件名模板，@n 为序号占位符（如 'book-@n' → book-1.epub）")
	renameCmd.Flags().BoolP("recursive", "r", false,
		"递归搜索子目录中的所有匹配文件")
	renameCmd.Flags().StringP("output", "o", "",
		"输出目录路径，指定后会将文件移动到此目录（不指定则在原位置重命名）")
	renameCmd.Flags().Int("start-num", 1,
		"序列号起始值（默认为 1）")
	renameCmd.Flags().Bool("do-try", false,
		"预览模式，仅显示将要执行的操作，不实际修改文件")
	renameCmd.Flags().Bool("debug", false,
		"启用调试模式，显示详细的配置信息")
	_ = rootCmd.MarkFlagRequired("format")
	_ = rootCmd.MarkFlagRequired("template")
}

func parseIntFlag(cmd *cobra.Command, name string) int {
	val, _ := cmd.Flags().GetInt(name)
	return val
}
