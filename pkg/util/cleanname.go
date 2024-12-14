package util

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
			println(title + "--------" + match[0])
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

func TryCleanTitle(title string) string {

	stack := &Stack{}

	//fmt.Println(title)
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
		//《大江大河》作者阿耐合集（共12册）（含《大江大河》(全4册)《欢乐颂》（全3册）《都挺好》(全2册)《不得往生》《食荤者》《余生》，阿耐出品，必是精品！作品改编影视剧均引起热议！）
		//fmt.Println("【needSymbol】 " + strings.Join(symbol.GetItems(), " | "))
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

		//fmt.Println(word, "----", strings.Join(stack.GetItems(), " | "))
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

func NewCleanTitle(title string) string {
	// 移除方括号 【】和 []
	reSquareBrackets := regexp.MustCompile(`[【\[].*?[】\]]`)
	title = reSquareBrackets.ReplaceAllString(title, "")

	// 定义需要保留的括号内容的正则表达式
	preservePatterns := []string{
		`.{2,6}篇`,
		`修订版`,
		`[上中下、]+[册本卷部]`,
		`套装.*?[册本卷部]`,
		`[全共].*?[册本卷部]`,
		`第.*?[版卷部]`,
		`[\d一二三四五六七八九十百千]+[-~—～][\d一二三四五六七八九十百千]+`,
		`\d{4}[-~—～]\d{4}`,
	}
	// 合并保留模式为一个正则表达式
	preserveRegex := regexp.MustCompile(strings.Join(preservePatterns, "|"))

	// 处理中文括号 （）
	reChineseParentheses := regexp.MustCompile(`（.*?）`)
	title = reChineseParentheses.ReplaceAllStringFunc(title, func(s string) string {
		content := strings.Trim(s, "（）")
		println(content)
		// 如果内容长度小于8，保留
		if len([]rune(content)) < 8 {
			return s
		}
		if preserveRegex.MatchString(content) {
			return s
		}
		return ""
	})

	// 处理英文括号 ()
	reEnglishParentheses := regexp.MustCompile(`\(.*?\)`)
	title = reEnglishParentheses.ReplaceAllStringFunc(title, func(s string) string {
		content := strings.Trim(s, "()")
		println(content)
		// 如果内容长度小于8，保留
		if len([]rune(content)) < 8 {
			return s
		}
		if preserveRegex.MatchString(content) {
			return s
		}
		return ""
	})

	// 处理中文英文括号 （)
	reChinese2Parentheses := regexp.MustCompile(`（.*?\)`)
	title = reChinese2Parentheses.ReplaceAllStringFunc(title, func(s string) string {
		content := strings.Trim(s, "（)")
		println(content)
		// 如果内容长度小于8，保留
		if len([]rune(content)) < 8 {
			return s
		}
		if preserveRegex.MatchString(content) {
			return s
		}
		return ""
	})
	// 处理英文中文括号 (）
	reEnglish2Parentheses := regexp.MustCompile(`\(.*?）`)
	title = reEnglish2Parentheses.ReplaceAllStringFunc(title, func(s string) string {
		content := strings.Trim(s, "(）")
		println(content)
		// 如果内容长度小于8，保留
		if len([]rune(content)) < 8 {
			return s
		}
		if preserveRegex.MatchString(content) {
			return s
		}
		return ""
	})

	// 移除书名号 《和》
	if strings.HasPrefix(title, "《") && strings.HasSuffix(title, "》") {
		title = title[3:]
		title = title[:len(title)-3]
	}
	// 去除首尾空格
	title = strings.ReplaceAll(title, "\"", " ")
	title = strings.TrimSpace(title)
	return title
}
