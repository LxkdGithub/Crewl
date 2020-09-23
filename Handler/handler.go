package Handler

import "C"
import (
	"Crewl/Config"
	"Crewl/Model"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/satori/uuid"
	"html/template"
	"log"
	"net/http"
	"strconv"
)


// middleware
func Start(c *gin.Context) bool {
	BiUUid, err := c.Cookie("BI_UUID")
	if err != nil {
		// 查找数据库uid是否过期 if (没有过期)
		return false
	}
	UserName, err := c.Cookie("USERNAME")
	if err != nil {
		return false
	}
	var userCookie Config.UserCookie
	userCookie.BI_UUID = BiUUid
	userCookie.UserName = UserName
	if Model.RedisAuth(&userCookie) == true {
		return true
	} else {
		return false
	}



}

// 从Model获取数据 渲染返回
func GetContentRet(c *gin.Context, page int) {
	rows, count, allCount := Model.GetComments(page)

	// ---------- Debug ---------
	//fmt.Println("allCount: ", allCount)

	div := Render(rows, count)
	// 返回
	UserName, err := c.Cookie("USERNAME")
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/Login")
	}
	var had bool = true
	if count <= 0 {
		had = false
	}
	pageCount := int(allCount / 10)
	if pageCount * 10 < allCount {
		pageCount += 1
	}
	var (
		HasNext bool = true
		HasPre bool = true
		PreLink string = "/Normal?page=" + strconv.Itoa(page-1)
		NextLink string = "/Normal?page=" + strconv.Itoa(page+1)
	)
	if pageCount == page {
		HasNext = false
		NextLink = ""
	}
	if page == 1 {
		HasPre = false
		PreLink = ""
	}

	// -------- Debug -----------
	//fmt.Println(page, pageCount, allCount, HasNext, NextLink, HasPre, PreLink)

	c.HTML(200, "show.html", gin.H{
		"Content": div,
		"Username": UserName,
		"had": had,

		"IsExcept": false,

		"page": page,
		"pageCount": pageCount,
		"HasPre": HasPre,
		"PreLink": PreLink,
		"HasNext": HasNext,
		"NextLink": NextLink,
	})	// 直接返回
}

func GetExceptRet(c *gin.Context, page int) {
	rows, count, allCount := Model.GetValidComments(page)
	fmt.Println("allCount: ", allCount)
	div := RenderExcept(rows, count)
	// 返回
	UserName, err := c.Cookie("USERNAME")
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/Login")
	}
	var had bool = true
	if count <= 0 {
		had = false
	}
	pageCount := int(allCount / 10)
	if pageCount * 10 < allCount {
		pageCount += 1
	}
	var (
		HasNext bool = true
		HasPre bool = true
		PreLink string = "/Execption?page=" + strconv.Itoa(page-1)
		NextLink string = "/Execption?page=" + strconv.Itoa(page+1)
	)
	if pageCount == page {
		HasNext = false
		NextLink = ""
	}
	if page == 1 {
		HasPre = false
		PreLink = ""
	}
	fmt.Println(page, pageCount, allCount, HasNext, NextLink, HasPre, PreLink)
	c.HTML(200, "show.html", gin.H{
		"Content": div,
		"Username": UserName,
		"had": had,

		"IsExcept": true,

		"page": page,
		"pageCount": pageCount,
		"HasPre": HasPre,
		"PreLink": PreLink,
		"HasNext": HasNext,
		"NextLink": NextLink,
	})	// 直接返回
}

func Login(c *gin.Context) {
	c.HTML(200, "Login.html", gin.H{
		"error": "",
	})
}

//登录处理
func LoginAction(c *gin.Context) {
	fmt.Println("LoginAction")
	var data Config.LoginForm
	var err error
	err = c.Bind(&data)
	if err != nil {
		log.Fatal("Param Bind Failed")
	}
	fmt.Println(data)
	if Model.Authorized(&data) == true {
		u1 := uuid.Must(uuid.NewV4(), err)  //Must edal with only the error auto...
		// 加入redis
		_, err := Model.Client.Set(u1.String(), data.Username, 0).Result()
		if err != nil {
			log.Fatal("Redis Set Failed")
		}
		// 设置cookie
		c.SetCookie("BI_UUID", u1.String(), 3600, "/", "", false, true)
		c.SetCookie("USERNAME", data.Username, 3600, "/", "", false, true)
		// 查表显示 Model
		c.Redirect(http.StatusMovedPermanently, "/Normal")

	} else {
		// 错误数据
		c.HTML(200, "Login.html", gin.H{
			"error": "Login Failed",
		}) // 错误信息
	}

}


// 正常数据显示处理
func ShowNormal(c *gin.Context) {
	var p int = 1
	page := c.Query("page")
	//fmt.Println(page)
	if page == "" || int(page[0]) > 57 || int(page[0]) < 48 {
		p = 1
	} else {
		p = int(page[0]) - 48
	}
	//fmt.Println(p)
	GetContentRet(c, p)
}

func ShowExcept(c *gin.Context) {
	var p int = 1
	page := c.Query("page")
	if page == "" || int(page[0]) > 57 || int(page[0]) < 48 {
		p = 1
	} else {
		p = int(page[0])
	}
	GetExceptRet(c, p)
}


//渲染模板
func Render(data []Config.Comment, count int) template.HTML {
	htmlContent := ""
	var err error
	t := template.Must(template.ParseFiles("Templates/view.html"))
	buffer := bytes.Buffer{}
	//就是将html文件里面的比那两替换为穿进去的数据
	for i:=0; i<count;i++ {
		if err = t.Execute(&buffer, data[i]); err != nil {
			log.Fatal(err, " --- template Execute Error")
		}
	}
	htmlContent += buffer.String()
	return template.HTML(htmlContent)
}


//渲染模板
func RenderExcept(data []Config.InvalidComment, count int) template.HTML {
	htmlContent := ""
	var err error
	fmt.Println("data", data)
	t := template.Must(template.ParseFiles("Templates/exception.html"))
	buffer := bytes.Buffer{}
	//就是将html文件里面的比那两替换为穿进去的数据
	for i:=0; i<count;i++ {
		if err = t.Execute(&buffer, data[i]); err != nil {
			fmt.Println(err, " --- template Execute Error")
		}
	}
	htmlContent += buffer.String()
	return template.HTML(htmlContent)
}