package main

import (
	"fmt"
	"encoding/json"
	. "./bot"
)


func main() {
	fmt.Println("Бот запущен!")
	var BOTS []Bot = make([]Bot, 0)

	var CFG map[string]interface{}
	json.Unmarshal([]byte(ReadFile("bots.json")), &CFG)

	for i:=0; i < len(List(CFG["bots"])); i++ {
		BOT := new(Bot)
		BOT.Config(Dict(List(CFG["bots"])[i]))
		BOTS = append(BOTS, *BOT)
	}

	CmdsInit()
	LP := new(LongPoll)
	LP.Init(BOTS)
	LP.Listen()
}