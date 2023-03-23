package cmd

import (
	"fmt"
	"github.com/kapmahc/epub"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// Used for downloading books from sanqiu website.
var c = &ClnameConfig{
	ReNameReg: regexp.MustCompile(`(?m)(\s?[(（【][^)）】(（【册卷套版]{4,}[)）】])`),
}

// renameBookCmd used for download books from sanqiu.cc
var clnameCmd = &cobra.Command{
	Use:   "clname",
	Short: "清理书籍标题中的无用描述",
	Run: func(cmd *cobra.Command, args []string) {
		// Validate config.

		ValidateConfig(c)

		if IsDir(c.Path) {
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
	if !Exists(c.Path) {
		fmt.Println("文件路径不存在，请检查")
		os.Exit(1)
	}
	if IsFile(c.Path) && !strings.HasSuffix(c.Path, ".epub") {
		fmt.Println("文件格式不存在，请检查")
		os.Exit(1)
	}
}

func init() {
	clnameCmd.Flags().StringVarP(&c.Path, "path", "p", "./",
		"目录或者文件")
	clnameCmd.Flags().BoolVarP(&c.Debug, "dotry", "t", false,
		"尝试运行")
	clnameCmd.Flags().BoolVarP(&c.Skip, "skip", "j", false,
		"跳过无法解析的书籍")
	clnameCmd.Flags().BoolVarP(&c.Debug, "debug", "d", false, "The number of download threads.")
}

func ParseEpub(file string, c *ClnameConfig) error {
	book, err := epub.Open(file)
	if err != nil {
		return err
	}
	if book == nil || len(book.Opf.Metadata.Title) == 0 {
		return fmt.Errorf("无法获得书籍标题")
	}
	title := book.Opf.Metadata.Title[0]

	if len(c.ReNameReg.FindAllString(title, -1)) == 0 {
		return nil
	}
	newTitle := strings.TrimSpace(c.ReNameReg.ReplaceAllString(title, ""))
	newTitle = strings.ReplaceAll(newTitle, "\"", " ")
	if len(newTitle) == 0 {
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
	Path  string
	DoTry bool
	Debug bool
	Skip  bool

	ReNameReg *regexp.Regexp
}

// Exists 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}
