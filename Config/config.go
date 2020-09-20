package Config


type LoginForm struct {
	Username string `form:username`
	Password string `from:password`
}

type Comment struct {
	Bid string
	Uid string
	Content string
}

type InvalidComment struct {
	Comment
	ValidKey []string
}

type LoginErrorStruct struct {
	err string
}