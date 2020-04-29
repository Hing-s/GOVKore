package bot

import (
	"math/rand"
	"time"
	"strconv"
	"fmt"
)

type COMMAND struct {
	access int
	cmd func(*Event)
	bot_id []int
}

var COMMANDS map[string] COMMAND
var HANDLERS map[string] []func(*Event, map[string]interface{})
var hadle bool


func Init(call func(*Event), command string, access int, bot_ids []int) {
	fmt.Println("Команда "+command+" иниализирована")
	cmd := new(COMMAND)

	cmd.cmd = call
	cmd.access = access
	cmd.bot_id = bot_ids

	COMMANDS[command] = *cmd
}

func HandleEvent(eventtype string, call func(*Event, map[string]interface{})) {
	if(!hadle) {
		HANDLERS = make(map[string] []func(*Event, map[string]interface{}))
		hadle = true
	}

	if(HANDLERS[eventtype] == nil) {
		HANDLERS[eventtype] = make([]func(*Event, map[string]interface{}), 0)
	}

	HANDLERS[eventtype] = append(HANDLERS[eventtype], call)
}

func RunHandlers(eventtype string, event *Event, updates map[string]interface{}) {
	if(!hadle) {
		return
	}

	handlers := HANDLERS[eventtype]

	if(handlers != nil) {
		for i:=0; i < len(handlers); i++ {
			handlers[i](event, updates)
		}
	}
}

/* 			INIT 			*/

func CmdsInit() {
	COMMANDS = make(map[string]COMMAND)

	Init(infa, "инфа", 0, append(make([]int, 0), 1, 2))
	Init(date, "дата", 0, append(make([]int, 0), 1, 2))
	Init(bottle, "бутылка", 0, append(make([]int, 0), 1, 2))
	Init(F, "f", 0, append(make([]int, 0), 1, 2))
}

func infa(event *Event) {
	rand.Seed(time.Now().UnixNano())
	percent := strconv.Itoa(rand.Intn(146))

	event.bot.message_send("Вероятность того, что "+event.args+" равна "+percent+"%", event.peer_id, nil)
}

func date(event *Event) {
	dt := time.Now()

	event.bot.message_send(dt.Format("Mon 15:04 01.02.2006"), event.peer_id, nil)
}

func bottle(event *Event) {
	rand.Seed(time.Now().UnixNano())

	message := "На бутылке у нас "
	params := make(map[string]interface{})

	params["peer_id"] = event.peer_id

	members := event.bot.method("messages.getConversationMembers", params)["response"]

	if(members == nil) {
		event.bot.message_send("Ошибочка вышла", event.peer_id, nil)
		return
	}

	members = Dict(members)["profiles"]
	user := Dict(List(members)[rand.Intn(len(List(members)))])

	params["attachment"] = "photo353166779_456284068"
	event.bot.message_send(message+"[id"+strconv.Itoa(int(user["id"].(float64)))+"|"+user["first_name"].(string)+ " "+user["last_name"].(string)+"]", event.peer_id, params)
}

func F(event *Event) {
	files := append(make([]File, 0), *new(File).Init("test.jpg", "doc"))
	event.bot.send_files(files, "тест", event.peer_id, nil)
}
