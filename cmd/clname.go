package cmd

import (
	"fmt"
	"github.com/jianyun8023/bookimporter/internal/util"
	"github.com/jianyun8023/bookimporter/pkg/epub"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// Used for downloading books from sanqiu website.
var c = &ClnameConfig{
	ReNameReg: regexp.MustCompile(`(?m)(\s?[(（【]([^)）】(（【册卷套版辑]|出版){4,}[)）】])`),
}

// renameBookCmd used for download books from sanqiu.cc
var clnameCmd = &cobra.Command{
	Use:   "clname",
	Short: "Clean up useless descriptions in book titles",
	Run: func(cmd *cobra.Command, args []string) {
		// Validate config.

		ValidateConfig(c)

		err := filepath.Walk(c.Path, func(path string, info os.FileInfo, err error) error {
			// 如果是子目录，则跳过继续处理
			if info.IsDir() {
				return nil
			}

			// 如果文件的扩展名是".epub"，则打印它的路径
			if filepath.Ext(path) != ".epub" {
				return nil
			}

			err = ParseEpub(path, c)

			if err != nil && !c.Skip {
				panic(fmt.Errorf("file: %v  %v", path, err))
			} else if err != nil && c.Skip {
				fmt.Printf("file: %v  %v", path, err)
			}
			return nil
		})

		if err != nil {
			panic(err)
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
		"Directory or file")
	clnameCmd.Flags().BoolVarP(&c.DoTry, "do-try", "t", false,
		"Try to run")
	clnameCmd.Flags().BoolVarP(&c.Skip, "skip", "j", false,
		"Skip books that cannot be parsed.")
	clnameCmd.Flags().BoolVarP(&c.Debug, "debug", "d", false, "Enable debugging information")
}

func ParseEpub(file string, c *ClnameConfig) error {
	metadata, err := epub.ReadMetadata(file)
	if err != nil {
		return err
	}
	if len(metadata.Title) == 0 {
		return fmt.Errorf("unable to obtain book title")
	}
	title := metadata.Title

	if len(c.ReNameReg.FindAllString(title, -1)) == 0 {
		return nil
	}
	newTitle := c.ReNameReg.ReplaceAllString(title, "")
	newTitle = strings.TrimSpace(strings.ReplaceAll(newTitle, "\"", " "))
	if len(newTitle) == 0 {
		return nil
	}

	fmt.Printf("Path: \033[1;31;40m%v\033[0m 新名称: \033[1;32;40m%v\033[0m 旧名称: \033[1;33;40m%v\033[0m\n", file, newTitle, title)

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
	Path  string
	DoTry bool
	Debug bool
	Skip  bool

	ReNameReg *regexp.Regexp
}
