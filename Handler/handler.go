package Handler

import (
	"github.com/gin-gonic/gin"
	"Crewl/Config"
	"Crewl/Model"
	"github.com/satori/uuid"
)


// 开始界面
func Start(c *gin.Context) {
	if uid, err := c.Cookie("UID"); err != nil {
		// 查找数据库uid是否过期 if (没有过期)
		c.JSON(200, uid)
	}
	c.HTML(200, "Login.html", nil)

}

//登录处理
func Login(c *gin.Context) {
	var data Config.LoginForm
	c.Bind(data)
	var err error
	if Model.Authorized(&data) == true {
		u1 := uuid.Must(uuid.NewV4(), err)  //Must edal with only the error auto...
		// 加入redis

		// 设置cookie
		c.SetCookie("CREWL_UID", u1.String(), 3600, "/", "", false, true)
		// 查表显示 Model

		// 返回
		c.HTML(200, "show.html", nil)	// 直接返回

	} else {
		// 错误数据
		var e LoginErrorStruct
		c.HTML(200, "Login.html", e) // 错误信息
	}

}
