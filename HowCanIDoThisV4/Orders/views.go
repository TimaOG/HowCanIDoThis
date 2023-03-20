package orders

type SortFields struct {
	TaskType   string `json:"taskType"`
	SecondType string `json:"secondType"`
	Deadline   string `json:"deadline"`
	WorkType   string `json:"workType"`
	Urgency    string `json:"Urgency"`
	PriceDown  string `json:"priceDown"`
	PriceUp    string `json:"priceUp"`
	Tags       string `json:"tags"`
	Offset     string `json:"offset"`
	OrderBy    string `json:"orderBy"`
}

type TaskTypeFirst struct {
	Id   uint8
	Name string
}
type TaskTypeSecond struct {
	Id          uint8
	Name        string
	FkFirstType uint8
}

type Order struct {
	Id              uint32
	Name            string
	FkUserOwner     uint32
	FkUserOwnerName string
	Discribtion     string
	Price           uint32
	Deadline        string
	Urgency         uint8
	WorkType        uint8
	Tags            string
	TzPath          string
}
type OrderById struct {
	Id              uint32 `json:"id"`
	Name            string `json:"name"`
	Discribtion     string `json:"discribtion"`
	FkUserOwner     uint32 `json:"fkUserOwner"`
	FkUserOwnerName string `json:"fkUserOwnerName"`
	Price           uint32 `json:"price"`
	Deadline        string `json:"deadline"`
	Urgency         uint8  `json:"urgency"`
	WorkType        uint8  `json:"workType"`
	Tags            string `json:"tags"`
	TzPath          string `json:"tzPath"`
}

type ViewData struct {
	Orders     []Order
	TypeFirst  []TaskTypeFirst
	TypeSecond []TaskTypeSecond
}
