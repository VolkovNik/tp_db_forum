package domain

type Vote struct {
	Nickname string `json:"nickname,omitempty"`
	Voice 	 int 	`json:"voice,omitempty"`
}
