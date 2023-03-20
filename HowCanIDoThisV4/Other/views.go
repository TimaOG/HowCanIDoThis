package other

type User struct {
	Id               uint32 `json:"id"`
	Name             string `json:"name"`
	Discribtion      string `json:"discribtion"`
	ProfileImg       string `json:"profileImg"`
	Rating           uint8  `json:"rating"`
	IsPremiumUser    bool   `json:"IsPremiumUser"`
	IsActiveUser     bool   `json:"IsActiveUser"`
	HistotyCount     uint32 `json:"historyCount"`
	Responsibility   uint8  `json:"Responsibility"`
	DoneOnTime       uint8  `json:"DoneOnTime"`
	AnswerSpead      uint8  `json:"AnswerSpead"`
	RegistrationDate string `json:"RegistrationDate"`
}
