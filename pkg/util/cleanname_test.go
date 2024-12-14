package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

type Item struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type Item2 struct {
	ID     int    `json:"id"`
	Title1 string `json:"title1"`
	Title2 string `json:"title2"`
}

func TestCleanTitle(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		//{"他的秘密【有些秘密注定要永远保守下去，除非你做好了失去一切的准备。《大小谎言》作者、澳洲小说天后莫里亚蒂成名作。被译介为121种版本！出版后盘踞《纽约时报》畅销榜近150周。】", "他的秘密"},
		//{"体坛周报（2024年第91期）", "体坛周报"},
		//{"历史的裂变：中国历史上的十三场政变（畅销书《大唐兴亡三百年》作者王觉仁力作，用小说笔法，讲述中华五千年历史上的13场知名政变，聚焦那些封建王朝中皇权的非正常更迭，还原权力争斗下最真实的人性。）", "历史的裂变：中国历史上的十三场政变"},
		//{"具象之力（世界科幻大师丛书）", "具象之力"},
		//{"幸运儿：晚清留美幼童的故事 (他们是大文豪马克・吐温的朋友。他们曾目睹一个神话般的时代。他们曾亲身经历近代中国的风云激荡；他们的命运，离奇而曲折；他们的故事，美丽而忧伤。他们有一个永远的名字：“留美幼童”。)", "幸运儿：晚清留美幼童的故事"},
		//{"武英殿本四库全书总目·上（1-30册）【电子版独家上线！国家图书馆倾情贡献！豆瓣9.6！】", "武英殿本四库全书总目·上（1-30册）"},
		//{"\"当代中国人文大系“精选（套装共35册）【人大出版社积累多年，集合当代名家著作，收罗中西方政治、哲学、历史研究精粹！】", "当代中国人文大系“精选（套装共35册）"},
		//{"版式设计法则", "版式设计法则"},
		//{"成功企业这样管理（套装12册）", "成功企业这样管理（套装12册）"},
		//{"深入理解Java虚拟机：JVM高级特性与最佳实践（第3版）", "深入理解Java虚拟机：JVM高级特性与最佳实践（第3版）"},
		//{"（第9版）公务员录用考试华图名家讲义系列教材：申论万能宝典", "（第9版）公务员录用考试华图名家讲义系列教材：申论万能宝典"},
		//{"第二座山（第一座山是构建自我、定义自我，其意义在于获取；第二座山是摆脱自我、舍弃自我，其意义在于奉献。《纽约时报》畅销书作者戴维·布鲁克斯全新作品，以新的诠释为人类生命的意义提出省思。）", "第二座山"},
		//{"《大江大河》作者阿耐合集（共12册）（含《大江大河》(全4册)《欢乐颂》（全3册）《都挺好》(全2册)《不得往生》《食荤者》《余生》，阿耐出品，必是精品！作品改编影视剧均引起热议！）", "《大江大河》作者阿耐合集（共12册）"},
		//{"喀耳刻（HBO改编8集系列剧即将上映！由里克·雅法（《侏罗纪世界》）和阿曼达·西尔弗（《猩球崛起》）担任编剧!严歌苓力荐：世上有一种作家，能够化古老的传说为神奇的故事。她让我们不是爱上历史，而是爱上那历史里的人！）", "喀耳刻"},
		{"民国印记（套装3本） 民国风度（回顾一个绝代芳华的时代，怀念一种活色生香的生活。曾经有那样一个时代，曾经有那样一批人物，他们那样地想着，那样地活着，有些清贫，有些窘迫，却又别样鲜活。） 民国印象：唯有时间，懂得爱（一切都会过去，时光会流逝，人会老去，唯有爱，不惧时光。时光沉淀了爱情的深，于是执手成说，于是那些人鲜活了，那些写满爱的旧纸张温软了。） 再见时光里的一瞥惊鸿（风物闲美、旧时掠影、遣怀故友等等，这里有你想看的每个类型的美文，它们承载着民国大家生活的点滴记忆。） ...", "民国印记（套装3本）"},
		{"课外英语-美国总统演讲选萃(上)（双语版）", "课外英语-美国总统演讲选萃(上)（双语版）"},
	}

	parseAll()

	for _, test := range tests {
		result := TryCleanTitle(test.input)
		if result != test.expected {
			t.Errorf("CleanTitle(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}

func parseAll() bool {
	// 打开 JSON 文件
	file, err := os.Open("/Users/zhaojianyun/Downloads/all.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return true
	}
	defer file.Close()

	// 读取文件内容
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return true
	}

	// 解析 JSON 数据
	var items []Item
	if err := json.Unmarshal(bytes, &items); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return true
	}

	var outdata []Item2
	// 打印解析后的数据
	for _, item := range items {

		oldTitle := item.Title
		newItem := Item2{
			ID:     item.ID,
			Title1: item.Title,
			Title2: TryCleanTitle(item.Title),
		}
		if newItem.Title2 == item.Title {
			continue
		}
		outdata = append(outdata, newItem)
		fmt.Printf("ID: %d, Title: %s ,OLD-Title: %s\n", item.ID, item.Title, oldTitle)
	}

	// 写入 JSON 文件
	outfile, err := os.Create("/Users/zhaojianyun/Downloads/out.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return true
	}
	defer outfile.Close()
	outbytes, err := json.MarshalIndent(outdata, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return true
	}
	if _, err := outfile.Write(outbytes); err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return true
	}
	return false
}
