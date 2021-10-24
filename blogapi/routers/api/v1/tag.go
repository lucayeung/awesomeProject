package v1

import (
	"blogapi/models"
	"blogapi/pkg/e"
	"blogapi/pkg/setting"
	"blogapi/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"

	"github.com/astaxie/beego/validation"
)

// curl http://localhost:8000/api/v1/tags
func GetTags(c *gin.Context) {
	// 组装数据访问参数
	maps := make(map[string]interface{})
	name := c.Query("name") // TODO 处理?name=test&state=1这种请求参数
	//name = c.DefaultQuery("name", "") // 支持默认值
	if name != "" {
		maps["name"] = name
	}
	var state = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	// 调用数据访问
	code := e.SUCCESS
	data := make(map[string]interface{})
	data["list"] = models.GetTags(util.GetPage(c), setting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)

	// 响应请求
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

// curl -X POST "http://localhost:8000/api/v1/tags?name=Kubernetes&state=1&created_by=Luca"
func AddTag(c *gin.Context) {
	name := c.Query("name") // req.getParameter("name");
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	createdBy := c.Query("created_by")

	vaild := validation.Validation{}
	vaild.Required(name, "name").Message("名称不能为空")
	vaild.MaxSize(name, 12, "name").Message("名称最长为12个字符")
	vaild.Required(createdBy, "created_by").Message("创建人不能为空")
	vaild.MaxSize(createdBy, 8, "created_by").Message("创建人最长为8个字符")
	vaild.Range(state, 0, 1, "state").Message("状态只允许0或1")

	var code int
	if vaild.HasErrors() {
		code = e.INVALID_PARAMS
	} else if models.ExistTagByName(name) {
		code = e.ERROR_EXIST_TAG
	} else {
		models.AddTag(name, state, createdBy)
		code = e.SUCCESS
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

func EditTag(c *gin.Context) {

}

func DeleteTag(c *gin.Context) {

}
