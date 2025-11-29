package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bookimporter",
	Short: "📚 书籍导入助手 - 管理和整理电子书库的强大工具",
	Long: `BookImporter 是一个用 Go 语言开发的书籍导入助手工具。

主要功能:
  • 清理书籍标题中的无用描述 (clname)
  • 检测 EPUB 文件完整性 (check)
  • 批量重命名文件 (rename)

使用示例:
  bookimporter check -p /books/     检测目录中的所有 EPUB 文件
  bookimporter clname -p /books/    清理书籍标题
  bookimporter rename . -f txt -t "book-@n"  批量重命名

项目地址: https://github.com/jianyun8023/bookimporter`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(clnameCmd)
	rootCmd.AddCommand(renameCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(checkCmd)
}
