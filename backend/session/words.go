package session

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var wordList []string = []string{"tiger", "world", "hello", "china", "apple", "child"}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GetWord() string {
	index := rand.Intn(len(wordList))
	fmt.Println(index)
	return strings.ToUpper(wordList[index])
}
