package offers

type SortFields struct {
	TaskType       string `json:"taskType"`
	SecondType     string `json:"secondType"`
	DaysToComplite string `json:"daysToComplite"`
	WorkType       string `json:"workType"`
	PriceDown      string `json:"priceDown"`
	PriceUp        string `json:"priceUp"`
	Rating         string `json:"rating"`
	SellerRating   string `json:"sellerRating"`
	Tags           string `json:"tags"`
	Offset         string `json:"offset"`
	OrderBy        string `json:"orderBy"`
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

type Offer struct {
	Id                 uint32
	Name               string
	FkUserOwner        uint32
	UserOwnerName      string
	IsPremiumUserOwner bool
	Price              uint32
	DaysToComplite     uint8
	WorkType           uint8
	Tags               string
	Rating             uint8
	HistoryCount       uint32
	CoverPath          string
}
type OfferById struct {
	Id                 uint32 `json:"id"`
	Name               string `json:"name"`
	Discribtion        string `json:"discribtion"`
	FkUserOwner        uint32 `json:"fkUserOwner"`
	UserOwnerName      string `json:"userOwnerName"`
	IsPremiumUserOwner bool   `json:"isPremiumUserOwner"`
	Price              uint32 `json:"price"`
	DaysToComplite     uint8  `json:"daysToComplite"`
	WorkType           uint8  `json:"workType"`
	Tags               string `json:"tags"`
	Rating             uint8  `json:"rating"`
	HistoryCount       uint32 `json:"historyCount"`
	CoverPath          string `json:"coverPath"`
}

type ViewData struct {
	Offers     []Offer
	TypeFirst  []TaskTypeFirst
	TypeSecond []TaskTypeSecond
}
