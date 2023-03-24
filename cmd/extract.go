package cmd

import (
	"fmt"
	"github.com/jianyun8023/bookimporter/internal/util"
	"github.com/jianyun8023/bookimporter/pkg/epub"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Used for downloading books from sanqiu website.
var extractConfig = &ExtractConfig{}

// renameBookCmd used for download books from sanqiu.cc
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract book metadata from epub file",
	Run: func(cmd *cobra.Command, args []string) {
		// Validate config.
		extractConfig.validateConfig()

		err := filepath.Walk(extractConfig.Path, func(path string, info os.FileInfo, err error) error {
			// 如果是子目录，则跳过继续处理
			if info.IsDir() {
				return nil
			}

			// 如果文件的扩展名是".epub"，则打印它的路径
			if filepath.Ext(path) != ".epub" {
				return nil
			}

			err = extract(path, extractConfig)

			if err != nil && !extractConfig.Skip {
				panic(fmt.Errorf("file: %v  %v", path, err))
			} else if err != nil && extractConfig.Skip {
				fmt.Printf("file: %v  %v", path, err)
			}

			return nil
		})

		if err != nil {
			panic(err)
		}

	},
}

func (c ExtractConfig) validateConfig() {
	if !util.Exists(c.Path) {
		fmt.Println("文件路径不存在，请检查")
		os.Exit(1)
	}
	if util.IsFile(c.Path) && !strings.HasSuffix(c.Path, ".epub") {
		fmt.Println("文件格式不存在，请检查")
		os.Exit(1)
	}
}

func extract(file string, c *ExtractConfig) error {

	isbn, err := epub.FindIsbn(file)
	if err != nil {
		return err
	}
	if isbn == "" {
		return nil
	}

	fmt.Printf("路径: \033[1;31;40m%v\033[0m ISBN: \033[1;32;40m%v\033[0m \n", file, isbn)

	if c.DoTry {
		return nil
	}
	//	var resultText []byte

	// 执行单个shell命令时, 直接运行即可
	command := fmt.Sprintf("ebook-meta \"%s\" --isbn=%s", file, isbn)
	cmdx := exec.Command("bash", "-c", command)

	if _, err = cmdx.Output(); err != nil {
		fmt.Println(cmdx.Env)
		fmt.Println(cmdx.Args)
		return err
		//		os.Exit(1)
	}
	//	fmt.Println(strings.Trim(string(resultText), "\n"))
	return nil

}

func init() {
	extractCmd.Flags().StringVarP(&extractConfig.Path, "path", "p", "./",
		"Directory or file")
	extractCmd.Flags().BoolVarP(&extractConfig.DoTry, "do-try", "t", true,
		"Try to run")
	extractCmd.Flags().BoolVarP(&extractConfig.Skip, "skip", "j", false,
		"Skip books that cannot be parsed.")
	extractCmd.Flags().BoolVar(&extractConfig.ISBN, "isbn", false, "Whether to extract ISBN")
	extractCmd.Flags().BoolVarP(&extractConfig.Debug, "debug", "d", false, "Enable debugging information.")
}

type ExtractConfig struct {
	Path  string
	DoTry bool
	Debug bool
	Skip  bool
	ISBN  bool
}
