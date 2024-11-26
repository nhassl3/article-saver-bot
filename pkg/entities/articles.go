package entities

type Updates struct {
	UpdateId int      `json:"update_id"`
	Message  *Message `json:"message"`
}

type UpdatesResponse struct {
	Ok     bool      `json:"ok"`
	Result []Updates `json:"result"`
}

type SendMessage struct {
	Offset int `json:"offset"`
}

type Message struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ChatId int `json:"id"`
}
