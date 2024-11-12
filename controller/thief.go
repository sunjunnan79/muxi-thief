package controller

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"muxi-thief/api/response"
	"muxi-thief/pkg"
	"net/http"
	"os"
	"sync"
	"time"
)

type GenerateJWTer interface {
	GenerateToken(code string) (string, error)
}

type Controller struct {
	jwt GenerateJWTer
}

// 存储 IP 的攻击请求信息
type IPInfo struct {
	Count      int
	FirstVisit time.Time
}

var ipMap sync.Map
var requestLimit = 5
var limitDuration = 3 * time.Second

func NewAuthController(jwt GenerateJWTer) *Controller {
	return &Controller{
		jwt: jwt,
	}
}

// Login 用户登录
// @Summary github用户登录授权接口
// @Description github用户登录授权接口,会自动重定向到github的授权接口上
// @Tags Auth
// @Produce json
// @Success 200 {object} response.Success "登录成功"
// @Failure 400 {object} response.Err "请求参数错误"
// @Failure 500 {object} response.Err "内部错误"
// @Router /api/v1/login [post]
func (c *Controller) Login(ctx *gin.Context) {
	code := ctx.GetHeader("code")
	token, err := c.jwt.GenerateToken(code)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Success{
			Data: "",
			Msg:  "警告!缺少code",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Success{
		Data: token,
		Msg:  fmt.Sprintf("恭喜你:%s,你已经成功登陆了XXHBGS的内部系统,但是请小心不要被抓获,现在你需要把返回的token作为值加到Authorization请求头中去(以后的请求都要加上)，并尝试去破坏这个系统。这个系统我们曾经攻破过，在这个网站里似乎能找到什么线索：%s ", code, BASEURL+"getStrategy"),
	})
	return
}

// GetStrategy 获取攻略
// @Summary 获取攻略
// @Description 获取攻略
// @Tags Auth
// @Produce json
// @Success 200 {object} string "获取攻略成功"
// @Router /api/v1/getStrategy [get]
func (c *Controller) GetStrategy(ctx *gin.Context) {

	message := fmt.Sprintf(`似乎是某人的十年前的日记:
	2014年11月1日:木犀成立了有几个月了,我们克服了很多困难,虽然日子还是很艰难,但是会好的。
	2014年11月5日:今天XXHBGS突然想要把我们干掉,还抢走了我们宝贵的文献,我们几个月的成果全毁了。
	2014年11月10日：经过全体后端组连续5天不眠不修的努力我们终于找到了夺回我们文献的办法，我们知道了怎么攻击XXHBGS的内部网站
	2014年11月11日：我们对：%s 网站发起了全面进攻，在进攻的过程中我们发现原来我们的文献藏在：%s 。
	2014年11月16日：经过不断的攻击和探索最终我们发现只要我们在向：%s 进行疯狂的攻击的同时，尝试去访问：%s 就能够成功找回我们被抢走的文献。
				   具体攻击方法如下:访问：%s 获取工具用的图片,通过短时间内并发请求发送攻击图片的方式来扰乱系统。PS:(3s内达到5次以上访问,别攻击太狠把服务器打挂了)
	...
`, BASEURL+"attack", BASEURL+"paper", BASEURL+"attack", BASEURL+"paper", BASEURL+"eyes")

	ctx.JSON(http.StatusOK, response.Success{
		Data: "",
		Msg:  message,
	})
	return
}

// Eyes 获取瞳孔
// @Summary 获取瞳孔
// @Description 获取瞳孔
// @Tags Auth
// @Produce json
// @Success 200 {object} string "获取瞳孔成功"
// @Router /api/v1/eyes [get]
func (c *Controller) Eyes(ctx *gin.Context) {

	code, err := getCode(ctx)
	if err != nil {
		return
	}

	number := pkg.SashToRange(code, 8)
	// 指定文件路径
	filePath := fmt.Sprintf("./file/%d.jpg", number) // 文件路径可以根据实际情况调整

	// 检查文件是否存在
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			ctx.JSON(http.StatusInternalServerError, response.Err{Err: errors.New("系统出错!:未找到系统文件").Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.Err{Err: errors.New("系统出错!:无法读取文件").Error()})
		return
	}

	// 读取文件内容
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Err{Err: errors.New("系统出错!:无法读取文件内容").Error()})
		return
	}

	// 将文件内容进行Base64编码
	encodedData := base64.StdEncoding.EncodeToString(fileData)

	// 构造JSON响应
	ctx.JSON(http.StatusOK, response.Success{
		Data: encodedData,
		Msg: fmt.Sprintf(`恭喜你,你现在已经拿到了虹膜,请使用虹膜扰乱XXHBGS的系统,找回我们的文献。
攻击原理如下:
	1.将虹膜图片使用form进行发送,file作为key,文件名称请设置为%s%d.jpg。
	2.http方法设置为PUT,通过替换其资源来攻击网站
	3.这个虹膜是使用base64加密传输的,你可以先将其保存再读取或者不保存直接发送
	PS:不要惊讶为什么是二次元头像,这才能扰乱敌人的系统!`, code, number),
	})
}

// Attack 系统攻击接口
// @Summary 系统攻击接口
// @Description 系统攻击接口，用于模拟对目标系统的攻击请求。通过频繁的文件上传请求来扰乱系统。
// @Tags Auth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "上传的图片文件, 文件名称格式为：{code}{随机数}.jpg"
// @Success 200 {object} response.Success "攻击请求成功，系统未检测到异常"
// @Failure 400 {object} response.Err "请求参数错误"
// @Failure 403 {object} response.Err "攻击次数超限，系统发出警告"
// @Failure 500 {object} response.Err "系统错误或文件名不正确"
// @Router /api/v1/attack [put]
func (c *Controller) Attack(ctx *gin.Context) {

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Err{Err: errors.New("请将上传的表单参数改为file!").Error()})
		return
	}

	//获取个人身份
	code, err := getCode(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Err{Err: errors.New("系统错误!").Error()})
		return
	}

	// 获取上传文件的名称
	fileName := file.Filename
	if fileName != fmt.Sprintf("%s%d.jpg", code, pkg.SashToRange(code, 8)) {
		ctx.JSON(http.StatusInternalServerError, response.Err{Err: errors.New("上传的文件名称错误!请仔细检查!").Error()})
		return
	}

	//ip
	ip := ctx.ClientIP()
	value, _ := ipMap.LoadOrStore(ip, &IPInfo{Count: 0, FirstVisit: time.Now()})
	ipInfo := value.(*IPInfo)
	now := time.Now()

	// 重置计数条件,如果达到限制时间的话
	if now.Sub(ipInfo.FirstVisit) > limitDuration {
		ipInfo.Count = 0
		ipInfo.FirstVisit = now
	} else {
		ipInfo.Count++
		ipMap.Store(ip, ipInfo)
	}

	if ipInfo.Count > requestLimit {
		ctx.JSON(http.StatusInternalServerError, response.Err{Err: errors.New("警告!系统遭到攻击").Error()})
		return
	} else {
		ctx.JSON(http.StatusOK, response.Success{Data: "", Msg: "系统运行正常"})
		return
	}
}

// Paper 文献获取接口
// @Summary 文献获取接口
// @Description 获取特定文献接口。用户通过频繁访问达到请求上限后，可在一定时间内获取该文献。
// @Tags Auth
// @Produce json
// @Success 200 {object} response.Success "成功获取文献"
// @Failure 403 {object} response.Err "用户未满足请求条件，无法获取文献"
// @Failure 500 {object} response.Err "系统错误"
// @Router /api/v1/paper [get]
func (c *Controller) Paper(ctx *gin.Context) {
	ip := ctx.ClientIP()
	value, ok := ipMap.Load(ip)

	if !ok {
		ctx.JSON(http.StatusInternalServerError, response.Err{Err: errors.New("系统错误!").Error()})
		return
	}

	code, err := getCode(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Err{Err: errors.New("系统错误!").Error()})

		return
	}

	ipInfo := value.(*IPInfo)
	if ipInfo.Count >= requestLimit && time.Since(ipInfo.FirstVisit) <= limitDuration {
		ctx.JSON(http.StatusOK, response.Success{Data: "Golang is the best language! Muxi is the best group! And You are the best freshman!", Msg: fmt.Sprintf("The End,恭喜你%s特工成功追回了孙院士的论文,挽救了木犀", code)})
	} else {
		ctx.JSON(http.StatusForbidden, response.Err{Err: errors.New("您没有访问权限").Error()})
	}

}

// Start 开始
// @Summary 开始游戏
// @Description 获取开始的用户认证
// @Tags Auth
// @Produce json
// @Success 200 {object} string "开始游戏成功!"
// @Router /api/v1/start [get]
func (c *Controller) Start(ctx *gin.Context) {
	// 格式化字符串，包含故事内容和登录网址
	message := fmt.Sprintf(`你醒了，好久不见。
我是jacksie的学弟jackson。
我们刚刚接到通知，之前你帮助孙院士找回的论文，在你破解的时候被敌对组织XXHBGS截获了。
现在我们需要你破解XXHBGS的安保设施，夺回孙院士的论文，并在他们察觉之前抢先发布。
现在请你使用POST方法访问这个网址尝试登录他们的内部系统：%s 请在请求头附上你的行动代号: code，以便我们植入破解木马。(code的内容随便你取)
`, BASEURL+"login")

	// 返回带格式的字符串
	ctx.String(http.StatusOK, message)
	return
}

func getCode(ctx *gin.Context) (string, error) {
	code, exist := ctx.Get("code")
	if !exist {
		return "", errors.New("get code from ctx err")
	}
	Code, ok := code.(string)
	if !ok {
		return "", errors.New("transform interface{} to string err")
	}
	return Code, nil
}
