package bot

import (
	"fmt"
	"strconv"
	"net/url"
)



const BASE_URL = "https://api.vk.com/method/"

type Bot struct
{
	token string // access_token
	names []interface{} // bot mentions
	version string // vk api version
	is_group bool // is bot group?
	poll LPData // data of LongPoll
	id int // bot vk id
	bot_id int // bot system id
	main int // link to another bot
	config map[string]interface{} // map of config

	_bots *[]Bot // another system bots
}


func (bot *Bot) SetBots(bots *[]Bot) {
	bot._bots = bots
}

func (bot *Bot) GetBots(bottype string) []Bot { 
	bots := make([]Bot, 0)

	for i:=0; i < len((*bot._bots)); i++ {

		if(bottype == "all") {
			if((*bot._bots)[i].token != bot.token) {
				bots = append(bots, (*bot._bots)[i])
			}
				
		} else {
			if((*bot._bots)[i].token != bot.token) {
				if(bottype == "group") {
					if ((*bot._bots)[i].is_group) {
						bots = append(bots, (*bot._bots)[i])
					}
				} else {
					if (!(*bot._bots)[i].is_group) {
						bots = append(bots, (*bot._bots)[i])
					}
				}
			}
		}
	}

	if len(bots) > 0 { return bots } else { return nil }
}

func (bot *Bot) Config(cfg map[string]interface{}) {
	bot.names = List(cfg["names"])

	if(bot.names[0] == nil) {
		return
	}

	bot.token = cfg["token"].(string)
	bot.version = cfg["v"].(string)
	bot.is_group = int(cfg["is_group"].(float64)) == 1
	bot.bot_id = int(cfg["bot_id"].(float64))
	bot.config = cfg

	if(cfg["main"] != nil) {
		bot.main = int(cfg["main"].(float64))
	} else {
		bot.main = -1
	}
}

func (bot *Bot) method(method string, params map[string]interface{}) map[string]interface{} {
	if(params == nil) {
		params = make(map[string]interface {})
	} 

	_params := url.Values{}

	if(params["access_token"] == nil) {
		params["access_token"] = bot.token
	}

	if(params["v"] == nil) {
		params["v"] = bot.version
	}

	for k, v := range params {
		_params.Add(k, v.(string))
	}

	return request(BASE_URL+method, _params)
}

func (bot *Bot) send_files(files []File, text string, peer_id string, params map[string]interface{}) map[string]interface{} {
	attachment := ""
	UploadParams :=  make(map[string] interface{})
	servers := make(map[string] interface{})

	if(params == nil) {
		params = make(map[string]interface {})
	}

	for i:=0; i < len(files); i++ {
		if(files[i].filetype == "photo" && servers["photo"] == nil) { 
			servers["photo"] = Dict(bot.method("photos.getMessagesUploadServer", nil))

			if(Dict(servers["photo"])["error"] != nil) {
				fmt.Println(Dict(servers["photo"])["error"])
				break
			}

		} else if(files[i].filetype == "doc" && servers["doc"] == nil) {
			UploadParams["peer_id"] = peer_id
			servers["doc"] = Dict(bot.method("docs.getMessagesUploadServer", UploadParams))

			if(Dict(servers["doc"])["error"] != nil) {
				fmt.Println(Dict(servers["doc"])["error"])
				break
			}
		}

		if(files[i].filetype == "photo") {
			response := UploadFile(string(Dict(Dict(servers[files[i].filetype])["response"])["upload_url"].(string)), files[i].path, "file1")

			if(response["server"] != nil) {
				UploadParams["album_id"] = "-3"
				UploadParams["server"] = strconv.Itoa(int(response["server"].(float64)))
				UploadParams["photo"] = string(response["photo"].(string))
				UploadParams["hash"] = string(response["hash"].(string))
				img := bot.method("photos.saveMessagesPhoto", UploadParams)

				if(img["error"] != nil) {
					continue
				}
				attachment+=fmt.Sprintf("%s%d_%d,", files[i].filetype, int(Dict(List(img["response"])[0])["owner_id"].(float64)), int(Dict(List(img["response"])[0])["id"].(float64)))
			}
			
		} else {
			response := UploadFile(string(Dict(Dict(servers[files[i].filetype])["response"])["upload_url"].(string)), files[i].path, "file")

			if(response["file"] != nil) {
				UploadParams["file"] = string(response["file"].(string))
				img := bot.method("docs.save", UploadParams)

				if(img["error"] != nil) {
					continue
				}

				attachment+=fmt.Sprintf("%s%d_%d,", files[i].filetype, int(Dict(Dict(img["response"])["doc"])["owner_id"].(float64)), int(Dict(Dict(img["response"])["doc"])["id"].(float64)))
			}
			
		}
	}

	params["attachment"] = attachment
	return bot.message_send("", peer_id, params)
}

func (bot *Bot) message_send(text string, peer_id string, params map[string]interface{}) map[string]interface{} {
	if(params == nil) {
		params = make(map[string]interface {})
	} 

	params["message"] = text
	params["peer_id"] = peer_id

	if(params["random_id"] == nil) {
		params["random_id"] = "0"
	}

	return bot.method("messages.send", params)
}