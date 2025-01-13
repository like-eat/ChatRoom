package v1

import (
	"github.com/gin-gonic/gin"
	"kama_chat_server/internal/dto/request"
	"kama_chat_server/internal/service/gorm"
	"kama_chat_server/pkg/constants"
	"net/http"
)

// GetMessageList 获取聊天记录
func GetMessageList(c *gin.Context) {
	var req request.GetMessageListRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, rsp, ret := gorm.MessageService.GetMessageList(req.UserOneId, req.UserTwoId)
	JsonBack(c, message, ret, rsp)
}

// GetGroupMessageList 获取群聊消息记录
func GetGroupMessageList(c *gin.Context) {
	var req request.GetGroupMessageListRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, rsp, ret := gorm.MessageService.GetGroupMessageList(req.GroupId)
	JsonBack(c, message, ret, rsp)
}

// UploadAvatar 上传头像
func UploadAvatar(c *gin.Context) {
	message, ret := gorm.MessageService.UploadAvatar(c)
	JsonBack(c, message, ret, nil)
}

// UploadFile 上传头像
func UploadFile(c *gin.Context) {
	message, ret := gorm.MessageService.UploadFile(c)
	JsonBack(c, message, ret, nil)
}

/*downloadFile 通用API，比如，下载文件，大家都可以调用这个api*/
func downloadFile(c *gin.Context) {
	//fmt.Println(c.Request.URL)
	////filePath := c.Query("path")
	//requestURL := fmt.Sprintf("%v", c.Request.URL)
	//requestURLarray := strings.Split(requestURL, "url=")
	//if len(requestURLarray) < 2 {
	//	sbjlog.Debug("downloadFile 失败 文件地址：%s ", requestURL)
	//	return
	//}
	//filePath := requestURLarray[1]
	//filePath = "./" + filePath
	////打开文件
	//fileTmp, errByOpenFile := os.Open(filePath)
	//if errByOpenFile != nil {
	//	sbjlog.Debug("downloadFile 失败 文件地址：%s 异常:%v", filePath, errByOpenFile)
	//	c.Redirect(http.StatusFound, "/404")
	//	return
	//}
	//defer fileTmp.Close()
	//
	////获取文件的名称
	//fileName := path.Base(filePath)
	//isExit, err := PathExists(filePath)
	//if err != nil {
	//	sbjlog.Debug("downloadFile 失败 文件地址：%s 异常:%v", filePath, err)
	//	c.Redirect(http.StatusFound, "/404")
	//	return
	//}
	//if !isExit {
	//	sbjlog.Debug("downloadFile 失败 文件地址：%s ", filePath)
	//	c.Redirect(http.StatusFound, "/404")
	//	return
	//}
	//c.Header("Content-Type", "application/octet-stream")
	////强制浏览器下载
	//c.Header("Content-Disposition", "attachment; filename="+fileName)
	////浏览器下载或预览
	//c.Header("Content-Disposition", "inline;filename="+fileName)
	//c.Header("Content-Transfer-Encoding", "binary")
	//c.Header("Cache-Control", "no-cache")
	//
	//c.File(filePath)
	//return
}
