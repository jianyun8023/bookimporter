package cmd

import (
	"fmt"
	"runtime"

	"github.com/jianyun8023/bookimporter/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	gitVersion = ""
	gitCommit  = "" // sha1 from git, output of $(git rev-parse HEAD)
	buildDate  = "" // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	goVersion  = runtime.Version()
	platform   = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
)

// versionCmd represents the version command.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示 BookImporter 版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		// 打印美化的版本信息
		fmt.Println(ui.RenderHeader("📚 BookImporter", "书籍导入助手工具"))
		fmt.Println()

		// 创建版本信息表格
		tableConfig := ui.NewTableConfig()
		tableConfig.Headers = []string{" 项目 ", " 值 "}
		tableConfig.BorderStyle = "rounded"
		tableConfig.CompactMode = false

		var rows [][]string

		// 版本号
		version := gitVersion
		if version == "" {
			version = "dev"
		}
		rows = append(rows, []string{
			" 版本 ",
			fmt.Sprintf(" %s ", version),
		})

		// Commit
		if gitCommit != "" {
			commit := gitCommit
			if len(commit) > 12 {
				commit = commit[:12]
			}
			rows = append(rows, []string{
				" Commit ",
				fmt.Sprintf(" %s ", commit),
			})
		}

		// 构建日期
		if buildDate != "" {
			rows = append(rows, []string{
				" 构建日期 ",
				fmt.Sprintf(" %s ", buildDate),
			})
		}

		// Go 版本
		rows = append(rows, []string{
			" Go 版本 ",
			fmt.Sprintf(" %s ", goVersion),
		})

		// 平台
		rows = append(rows, []string{
			" 平台 ",
			fmt.Sprintf(" %s ", platform),
		})

		tableConfig.Rows = rows
		table := ui.NewTable(tableConfig)
		fmt.Println(table.Render())
		fmt.Println()

		// 项目信息
		fmt.Println(ui.RenderInfo("项目地址: https://github.com/jianyun8023/bookimporter"))
		fmt.Println(ui.RenderInfo("使用 'bookimporter --help' 查看帮助信息"))
	},
}
