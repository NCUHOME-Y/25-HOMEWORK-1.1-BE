package main

// 本次作业较面试考核修改了随机数生成的方式，不再使用已弃用的随机数生成方式；改变了部分变量的命名方式，尽可能采取了驼峰命名法；修复了上次考核中
// 未能实现的每个手机号每日只能申请五次验证码的功能；增加了足够的注释，提高了代码的可读性；修改了部分文字描述，更加人性化。


import (
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

// 验证码信息
type CodeInformation struct {
	code        string
	IsUsed      bool
	CreatedTime time.Time
}

// 用户登录信息
type UserInformation struct {
	phonenumber  string
	LastTime     time.Time
	LastDate     time.Time // 记录日期，精确到天
	CreatedCount int       // 当日申请验证码次数
	Information  *CodeInformation
}

// 用户认证服务
type UserService struct {
	users map[string]*UserInformation
}

// 初始化
func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]*UserInformation),
	}
}

// 生成字母和数字混合验证码（使用新版随机数生成方式）
func (s *UserService) GenerateAlphanumericCode() string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6)
	
	// 使用新版随机数生成方式
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range code {
		code[i] = charset[r.Intn(len(charset))]
	}
	return string(code)
}

// 检验手机号格式
func (s *UserService) Phonevalidate(phonenumber string) bool {
	var phoneRegexp = regexp.MustCompile(`^1[3-9]\d{9}$`)
	return phoneRegexp.MatchString(phonenumber)
}

// 检验能否请求验证码
func (s *UserService) CanRequestCode(phonenumber string) (bool, string) {
	user, ok := s.users[phonenumber]
	currentTime := time.Now()
	// 只保留日期部分（年/月/日）用于比较
	currentDate := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())

	if !ok {
		return true, ""
	}
	
	// 检查是否跨天，如果跨天则重置计数
	if !user.LastDate.Equal(currentDate) {
		user.CreatedCount = 0
		user.LastDate = currentDate
		return true, ""
	}
	
	// 检查每日次数上限
	if user.CreatedCount >= 5 {
		return false, "今日验证码请求次数已达上限"
	}
	
	// 检查一分钟冷却
	if time.Since(user.LastTime) < 60*time.Second {
		return false, "一分钟内已获取验证码，请稍后再试"
	}
	
	return true, ""
}

// 产生验证码
func (s *UserService) GenerateCode(phonenumber string) (string, error) {
	// 先验证手机号格式
	if !s.Phonevalidate(phonenumber) {
		return "", fmt.Errorf("手机号码格式不正确")
	}
	
	// 检查是否可以请求验证码
	if canRequest, reason := s.CanRequestCode(phonenumber); !canRequest {
		return "", fmt.Errorf(reason)
	}
	
	code := s.GenerateAlphanumericCode()
	currentTime := time.Now()
	currentDate := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())

	// 如果用户不存在则创建
	if _, ok := s.users[phonenumber]; !ok {
		s.users[phonenumber] = &UserInformation{
			phonenumber: phonenumber,
			LastDate:    currentDate,
		}
	}
	
	// 更新用户信息
	user := s.users[phonenumber]
	user.Information = &CodeInformation{
		code:        code,
		IsUsed:      false,
		CreatedTime: currentTime,
	}
	user.LastTime = currentTime
	user.CreatedCount++
	
	return code, nil
}

// 检验验证码是否符合要求
func (s *UserService) ValidateCode(phonenumber, code string) bool {
	user, ok := s.users[phonenumber]
	if !ok || user.Information == nil {
		return false
	}
	
	// 检查是否已使用
	if user.Information.IsUsed {
		return false
	}
	
	// 检查是否过期（5分钟）
	if time.Since(user.Information.CreatedTime) > 5*time.Minute {
		return false
	}
	
	// 检查验证码是否正确
	return code == user.Information.code
}

// 标记验证码为已使用
func (s *UserService) MarkCodeAsUsed(phonenumber string) {
	if user, ok := s.users[phonenumber]; ok && user.Information != nil {
		user.Information.IsUsed = true
	}
}

func main() {
	userService := NewUserService()
	var phonenumber string
	
	fmt.Println("请输入电话号码")
	fmt.Scanln(&phonenumber)
	
	if !userService.Phonevalidate(phonenumber) {
		fmt.Println("号码格式不正确")
		return
	}

	for {
		fmt.Println("1:输入验证码进行登录  2:获取验证码")
		var number string
		fmt.Scanln(&number)

		switch number {
		case "1":
			fmt.Println("请输入验证码登录")
			var code string
			fmt.Scanln(&code)

			if userService.ValidateCode(phonenumber, code) {
				fmt.Println("登陆成功")
				userService.MarkCodeAsUsed(phonenumber)
				return
			} else {
				fmt.Println("无效验证码或已过期")
			}
		case "2":
			code, err := userService.GenerateCode(phonenumber)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Printf("验证码已发送: %s\n", code)
			}
		default:
			fmt.Println("未知操作，请输入1或2")
		}
	}
}
