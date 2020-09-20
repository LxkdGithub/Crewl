package Handler

import (
	"github.com/gin-gonic/gin"
	"Crewl/Config"
	"Crewl/Model"
)


// 开始界面
func Start(c *gin.Context) {
	c.HTML(200, "Login.html", nil)
}

//登录处理
func Login(c *gin.Context) {
	var data Config.LoginForm
	c.Bind(data)
	if Model.Authorized(&data) == true {
		c.SetCookie(data.Username, "sd")
		c.HTML()	// 直接返回
	} else {
		c.HTML() // 错误信息
	}

}
