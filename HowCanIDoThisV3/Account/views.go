package account

type UserInfo struct {
	Name             string
	Email            string
	Discribtion      string
	Rating           uint8
	ProfileImg       string
	CountFeedbacks   uint32
	IsPremiumUser    bool
	IsActiveUser     bool
	HistotyCount     uint32
	Responsibility   uint8
	DoneOnTime       uint8
	AnswerSpead      uint8
	RegistrationDate string
	Balance          uint32
}
type UserOrders struct {
	Id          uint32
	Name        string
	Discribtion string
	Price       uint32
}
type UserOffers struct {
	Id        uint32
	Name      string
	Price     uint32
	IsActive  bool
	CoverPath string
}

type UserOrdersDoingNow struct {
	Id              uint32
	Name            string
	Discribtion     string
	Price           uint32
	ExecutorName    string
	ExecutorId      uint32
	CustomerName    string
	CustomerId      uint32
	StartTime       string
	ExpectedEndTime string
}
type UserOffersDoingNow struct {
	Id              uint32
	Name            string
	Price           uint32
	ExecutorName    string
	ExecutorId      uint32
	CustomerName    string
	CustomerId      uint32
	StartTime       string
	ExpectedEndTime string
}
type UserChats struct {
	UserChatId      uint32
	UserSecondId    uint32
	UserSecondName  string
	UserSecondImg   string
	LastMessageText string
	LastMessageTime string
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

type ViewData struct {
	User           UserInfo
	Orders         []UserOrders
	Offers         []UserOffers
	OrdersDoingNow []UserOrdersDoingNow
	OffersDoingNow []UserOffersDoingNow
	Chats          []UserChats
	TypeFirst      []TaskTypeFirst
	TypeSecond     []TaskTypeSecond
}

type Order struct {
	Name                string
	FkUserOwner         uint32
	Discribtion         string
	Price               uint32
	Deadline            string
	Urgency             uint8
	WorkType            uint8
	OrderCategory       uint8
	OrderCategorySecond uint8
	Tags                string
	IsActive            bool
	TzPath              string
}
type Offer struct {
	Name                string
	FkUserOwner         uint32
	Discribtion         string
	Price               uint32
	DaysToComplite      uint8
	WorkType            uint8
	OrderCategory       uint8
	OrderCategorySecond uint8
	Tags                string
	IsActive            bool
	CoverPath           string
}
type OfferById struct {
	Status              string `json:"status"`
	Name                string `json:"name"`
	Discribtion         string `json:"discribtion"`
	Price               uint32 `json:"price"`
	DaysToComplite      uint8  `json:"daysToComplite"`
	WorkType            uint8  `json:"workType"`
	OrderCategory       uint8  `json:"orderCategory"`
	OrderCategorySecond uint8  `json:"orderCategorySecond"`
	Tags                string `json:"tags"`
}
