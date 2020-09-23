package Config


type LoginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type Comment struct {
	Id 		string
	Bid     string
	Uid     string
	Content string
	Time 	string
}

type InvalidComment struct {
	Id 		string
	Bid     string
	Uid     string
	Content string
	Time 	string
	Sensitive string
}

type UserCookie struct {
	BI_UUID string
	UserName string
}




