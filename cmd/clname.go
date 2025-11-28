package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

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

		if util.IsDir(c.Path) {
			m, _ := filepath.Glob(path.Join(c.Path, "*.epub"))
			for _, val := range m {
				//				fmt.Println(val)
				epubpath := val
				err := ParseEpub(epubpath, c)
				if err != nil && !c.Skip {
					panic(fmt.Errorf("file %v  %v", epubpath, err))
				} else if err != nil && c.Skip {
					fmt.Printf("file %v  %v\n", epubpath, err)
				}
			}

		} else {
			epubpath := c.Path
			err := ParseEpub(epubpath, c)
			if err != nil && !c.Skip {
				panic(fmt.Errorf("file %v  %v", epubpath, err))
			} else if err != nil && c.Skip {
				fmt.Printf("file %v  %v\n", epubpath, err)
			}
		}
	},
}

func ValidateConfig(c *ClnameConfig) {
	if !util.Exists(c.Path) {
		fmt.Println("文件路径不存在，请检查")
		os.Exit(1)
	}
	if util.IsFile(c.Path) && !strings.HasSuffix(c.Path, ".epub") {
		fmt.Println("文件格式不存在，请检查")
		os.Exit(1)
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

func ParseEpub(file string, c *ClnameConfig) error {
	// 预先检测 EPUB 文件完整性
	if err := util.ValidateEpubFile(file); err != nil {
		if c.Debug {
			fmt.Printf("EPUB 文件检测失败: %v\n", err)
		}

		// 处理损坏的文件
		if c.MoveCorruptedTo != "" {
			if c.DoTry {
				fmt.Printf("  → [试运行] 将移动损坏文件到: %s\n", c.MoveCorruptedTo)
			} else {
				newPath, moveErr := util.MoveFileWithConflictHandling(file, c.MoveCorruptedTo)
				if moveErr != nil {
					fmt.Printf("  → 移动损坏文件失败: %v\n", moveErr)
				} else {
					fmt.Printf("  → 已移动损坏文件到: %s\n", newPath)
				}
			}
		} else if c.DeleteCorrupted {
			if c.DoTry {
				fmt.Printf("  → [试运行] 将删除损坏文件\n")
			} else {
				needConfirm := !c.ForceDelete
				deleteErr := util.SafeDeleteFile(file, needConfirm)
				if deleteErr != nil {
					fmt.Printf("  → 删除损坏文件失败: %v\n", deleteErr)
				} else {
					fmt.Printf("  → 已删除损坏文件\n")
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
		return nil
	}

	fmt.Printf("路径: \033[1;31;40m%v\033[0m 新名称: \033[1;32;40m%v\033[0m 旧名称: \033[1;33;40m%v\033[0m\n", file, newTitle, title)

	if c.DoTry {
		return nil
	}
	//	var resultText []byte

	// 执行单个shell命令时, 直接运行即可
	command := fmt.Sprintf("ebook-meta \"%s\" -t \"%s\"", file, newTitle)
	cmdx := exec.Command("bash", "-c", command)

	if _, err = cmdx.Output(); err != nil {
		fmt.Println(cmdx.Env)
		fmt.Println(cmdx.Args)
		return err
		//		os.Exit(1)
	}
	//	fmt.Println(strings.Trim(string(resultText), "\n"))
	return err
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
