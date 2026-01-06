package naming

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	ReNameReg = regexp.MustCompile(`(?m)(\s?[(（【][^)）】(（【册卷套]{4,}[)）】])`)
)

func CleanTitle(title string) string {
	if len(ReNameReg.FindAllString(title, -1)) == 0 {
		return title
	}

	for _, match := range ReNameReg.FindAllStringSubmatch(title, -1) {
		if len(match[0]) < 10 {
			// Keeping println for debugging logic as in original,
			// though logging interface would be better in DDD.
			// println(title + "--------" + match[0])
			return title
		}
	}
	newTitle := ReNameReg.ReplaceAllString(title, "")
	newTitle = strings.TrimSpace(strings.ReplaceAll(newTitle, "\"", " "))
	return newTitle
}

type Stack struct {
	items []string
}

func (s *Stack) Push(item string) {
	s.items = append(s.items, item)
}

func (s *Stack) Pop() string {
	if len(s.items) == 0 {
		return ""
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack) Peek() string {
	if len(s.items) == 0 {
		return ""
	}
	return s.items[len(s.items)-1]
}

func (s *Stack) GetItems() []string {
	return s.items
}

// Cleaner 封装标题清洗逻辑
type Cleaner struct{}

// NewCleaner 创建一个新的清洗器
func NewCleaner() *Cleaner {
	return &Cleaner{}
}

// Clean 清洗标题
func (c *Cleaner) Clean(title string) string {
	return TryCleanTitle(title)
}

func TryCleanTitle(title string) string {

	stack := &Stack{}

	reader := strings.NewReader(title)

	word := ""
	symbol := &Stack{}
	symbolPair := map[string]string{
		"【": "】",
		"[": "]",
		"（": "）",
		"(": ")",
		"】": "【",
		"]": "[",
		"）": "（",
		")": "(",
	}
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			break
		}
		char := string(r)

		switch char {
		case "【", "[", "（", "(":
			if symbol.IsEmpty() && utf8.RuneCountInString(word) > 0 {
				stack.Push(word)
				word = char
			} else {
				word += char
			}
			symbol.Push(char)
		case "】", "]", "）", ")":
			word += char
			if symbol.Peek() == symbolPair[char] {
				symbol.Pop()
				stack.Push(word)
				word = ""
			}
		default:
			word += char
		}
	}

	if symbol.IsEmpty() || len(word) != 0 {
		stack.Push(word)
	}

	outTitle := ""
	for i, v := range stack.GetItems() {
		if i == 0 {
			outTitle = v
		} else if i <= 2 && preserve(v) {
			outTitle += v
		}
	}

	// 去除首尾空格
	outTitle = strings.ReplaceAll(outTitle, "\"", " ")
	outTitle = strings.TrimSpace(outTitle)
	if utf8.RuneCountInString(outTitle) == 0 {
		return title
	}
	return outTitle
}

func preserve(content string) bool {

	c := strings.TrimPrefix(content, "【")
	c = strings.TrimPrefix(c, "[")
	c = strings.TrimPrefix(c, "（")
	c = strings.TrimPrefix(c, "(")

	c = strings.TrimSuffix(c, "】")
	c = strings.TrimSuffix(c, "]")
	c = strings.TrimSuffix(c, "）")
	c = strings.TrimSuffix(c, ")")

	if utf8.RuneCountInString(c) <= 3 {
		return true
	}
	// 定义需要保留的括号内容的正则表达式
	preservePatterns := []string{
		`.{2,6}篇`,
		`[上中下+]`,
		`[上中下、]+[册本卷部辑]`,
		`套装.*?[册本卷部辑]`,
		`[全共].*?[册本卷部辑]`,
		`\d+[册本卷部辑]`,
		`第.*?[版卷部辑]`,
		`[\d一二三四五六七八九十百千]+[-~—～][\d一二三四五六七八九十百千]+`,
		`\d{4}[-~—～]\d{4}`,
	}
	// 合并保留模式为一个正则表达式
	preserveRegex := regexp.MustCompile(strings.Join(preservePatterns, "|"))
	b := preserveRegex.MatchString(c) && utf8.RuneCountInString(c) < 20
	if !b {
		return strings.HasSuffix(c, "版") && utf8.RuneCountInString(c) < 10
	}
	return b
}
