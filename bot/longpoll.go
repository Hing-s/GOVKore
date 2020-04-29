package bot

import (
	"fmt"
	"strconv"
	"time"
	"strings"
	"net/url"
)


type LPData struct {
	server string // LongPoll sever url
	key string // LongPoll access_key
	ts string // LongPoll ts

	updates []interface{} // updates map
}

type Event struct {
	peer_id string // peer_id of update
	text string // text of update
	userid int // from_id of update
	splited []string // splited text with " "
	args string // cleaned message of mention and command
	bot *Bot // bot what recievs update
	etype string // type of event
	other map[string] string // optional data
}

func (event *Event) group(object map[string]interface{}, bot *Bot) *Event {
	event.etype = object["type"].(string)

	if(event.etype == "message_new") {
		object := Dict(object["object"])

		event.text = object["text"].(string)
		event.peer_id = strconv.Itoa(int(object["peer_id"].(float64)))
		event.userid = int(object["from_id"].(float64))
		event.bot = bot
		event.splited = strings.Split(strings.ToLower(event.text), " ")

		event.etype = "MSG"
	}

	return event
}

func (event *Event) page(object []interface{}, bot *Bot) *Event {
	event.etype = strconv.Itoa(int(object[0].(float64)))

	if(int(object[0].(float64)) == 4) {
		peer_id := int(object[3].(float64))

		event.text = string(object[5].(string))
		event.bot = bot

		if(peer_id < 2000000000) {
			event.userid = peer_id
		} else {
			userid, err := strconv.Atoi(string(Dict(object[6])["from"].(string)))
			
			if(err != nil) {
				return nil
			}

			event.userid = userid
		}

		event.etype = "MSG"
		event.peer_id = strconv.Itoa(peer_id)
		event.splited = strings.Split(strings.ToLower(event.text), " ")
	}

	return event
}


func (poll *LPData) GetLongPoll(bot *Bot) {
	var server map[string] interface{}

	params := make(map[string]interface {})

	if(bot.is_group) {
		id := Dict(List(bot.method("groups.getById", nil)["response"])[0])["id"]

		if(id != nil) {
			bot.id = int(id.(float64))
			params["group_id"] = strconv.Itoa(bot.id)
			server = Dict(bot.method("groups.getLongPollServer", params)["response"])
		} else {
			bot.poll.server = ""
			return
		}

	} else {
		id := Dict(List(bot.method("users.get", nil)["response"])[0])["id"]

		if(id != nil) {
			bot.id = int(id.(float64))
			server = Dict(bot.method("messages.getLongPollServer", params)["response"])
		} else {
			bot.poll.server = ""
			return
		}
		
	}

	if(bot.is_group) {
		fmt.Println("LongPoll -"+strconv.Itoa(bot.id)+" was successfully getted")
	} else {
		fmt.Println("LongPoll "+strconv.Itoa(bot.id)+" was successfully getted")
	}
	
	if(bot.is_group) {
		bot.poll.server = server["server"].(string)
		bot.poll.ts = server["ts"].(string)
		bot.poll.key = server["key"].(string)
	} else {
		bot.poll.server = "https://"+string(server["server"].(string))
		bot.poll.ts = strconv.Itoa(int(server["ts"].(float64)))
		bot.poll.key = string(server["key"].(string))
	}	
}

func (poll *LPData) Updates(bot *Bot) {
	var event *Event

	if(bot.poll.server == "") {
		return
	}

	for true {
		updates := request(poll.server+"?act=a_check", url.Values{"key": {poll.key}, "ts": {poll.ts}, "wait": {"30"}, "version": {"3"}, "mode": {"2"}})

		if(updates["ts"] == nil) {
			poll.GetLongPoll(bot)
		}

		if(bot.is_group) {
			poll.ts = updates["ts"].(string)
		} else {
			poll.ts = strconv.Itoa(int(updates["ts"].(float64)))
		}
		

		if(updates["updates"] == nil) {
			continue
		}

		for i:=0; i < len(List(updates["updates"])); i++ {
			if(bot.is_group) {
				object := Dict(List(updates["updates"])[i])
				event = new(Event).group(object, bot)
			} else {
				object := List(List(updates["updates"])[i])
				event = new(Event).page(object, bot)
			}

			if(event == nil) {
				continue
			}

			RunHandlers(event.etype, event, updates)
			
			if(event.etype == "MSG") {
				if(len(event.splited) > 1) {

					if(len(event.splited) > 2) {
						event.args = event.text[len(event.splited[0])+len(event.splited[1])+2:]
					} else {
						event.args = event.text[len(event.splited[0])+1:]
					}

					if(bot.main != -1 && bot.bot_id != bot.main) {
						continue
					}

					if(HasString(bot.names, strings.ToLower(event.splited[0]))) {
						cmd := COMMANDS[event.splited[1]].cmd
						
						if(cmd != nil) {
							cmd(event)
						} else {
							bot.message_send("Команда не найдена!", event.peer_id, nil) // cmd not found
						}
					}
				}
			}
		}
	}
}

type LongPoll struct {
	bots []Bot; // all system bots
}

func (lp *LongPoll) Init(bots []Bot) {
	lp.bots = bots
}

func (lp *LongPoll) Listen() {
	errs := make([]Bot, 0)

	for i:=0; i < len(lp.bots); i++ {
		lp.bots[i].poll.GetLongPoll(&lp.bots[i])

		if(lp.bots[i].poll.server == "") {
			errs = append(errs, lp.bots[i])
			fmt.Println("Bot with system id "+strconv.Itoa(lp.bots[i].bot_id)+" doesn't launched")
			continue
		}

		lp.bots[i].SetBots(&lp.bots)

		go lp.bots[i].poll.Updates(&lp.bots[i])
	}

	if(len(errs) == len(lp.bots)) {
		fmt.Println("No one bot was launched successfully")
		return
	}

	// ¯\_(ツ)_/¯
	for true {
		time.Sleep(1000000)
	}
}
