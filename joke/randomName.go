package joke

import (
	"fmt"
	"math/rand"
	"time"
)

// 随机生成一个元素
func randomElement(elements []string) string {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(elements))
	return elements[index]
}

// 生成神话人物名字的函数
func GenerateMythicalName() string {
	// 姓氏
	surnames := []string{"伏羲", "女娲", "夸父", "后羿", "嫦娥", "黄帝", "蚩尤", "大禹", "嬴政"}

	// 名字
	names := []string{"天明", "明珠", "飞龙", "云霞", "瑶池", "紫霞", "明月", "丽华", "紫宸", "长春"}

	// 生成神话人物名字
	mythicalName := fmt.Sprintf("%s%s", randomElement(surnames), randomElement(names))

	return mythicalName
}
