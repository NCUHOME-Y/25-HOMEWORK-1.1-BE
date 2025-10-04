package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"sync"
	"time"
)

type CodeSend struct {
	lastphonetime map[string]time.Time //时间
	mutex         sync.Mutex           //锁
	password      map[string]string    //存储验证码
}

// 初始化:用常量赋予结构体变量
func CodeSendInit() *CodeSend {
	return &CodeSend{
		lastphonetime: make(map[string]time.Time),
		password:      make(map[string]string),
	}
}

// 接口
type CodeService interface {
	CanSendCode(phone string) error
	JudgeCode(phone string, m string) bool
	Clear(phone string)
	Random() string
}

func (c *CodeSend) CanSendCode(phone string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock() //加锁

	//60秒内,同⼀⼿机号不能重复获取验证码
	lastTime, exits := c.lastphonetime[phone]
	if exits && time.Since(lastTime) <= 60*time.Second { //检验存在和是否再60秒之内
		return fmt.Errorf("60秒内,同⼀⼿机号不能重复获取验证码")
	}
	if counter >= 5 {
		return fmt.Errorf("当日发送次数已达上限(5次)")
	}
	//记录时间
	c.lastphonetime[phone] = time.Now()
	counter++
	return nil
}

// 当日最多发送五次
var counter int = 0 //全局变量，记录发送次数
func (c *CodeSend) DaliyLimit(phone string) int {
	// 计算到第二天0点的时间
	now := c.lastphonetime[phone] //获取上次发送验证码的时间
	nextDay := now.Add(24 * time.Hour)
	zeroTime := time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 0, nextDay.Location()) //nextDay.Location() 的作用是获取 nextDay 这个时间对象所属的时区信息
	duration := zeroTime.Sub(now)

	// 等待到第二天0点
	time.Sleep(duration)

	// 清零计数器
	counter = 0
	return counter
}

// 验证验证码是否有效
func (c *CodeSend) JudgeCode(phone string, m string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	info, exit := c.password[phone]
	if !exit {
		return false
	}
	return info == m && time.Since(c.lastphonetime[phone]) <= 5*time.Minute
}

// 登录后删除验证码
func (c *CodeSend) Clear(phone string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.password, phone)
}

// 生成随机数
func (c *CodeSend) Random() string {
	const char = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	charlen := big.NewInt(int64(len(char)))
	result1 := make([]byte, 6)
	for i := range result1 {
		num, _ := rand.Int(rand.Reader, charlen)
		result1[i] = char[num.Int64()]
	}
	return string(result1)
}

func main() {
	var service CodeService = CodeSendInit()

	//判定手机号
	fmt.Println("请输入手机号")
	var phone string
	fmt.Scanln(&phone)
	//正则表达式
	reg := regexp.MustCompile(`^1[3-9]\d{9}$`)
	result1 := reg.MatchString(phone) //返回切片. 而Matchstring  返回一个bool类型

	if !result1 {
		fmt.Println("登录失败,请输入正确手机号码")

	} //判断手机号 https://studygolang.com/pkgdoc 语法链接

	if result1 {
		//判断是否生成验证码
		//循环选择
		for {
			var order int
			fmt.Println("输⼊1:表⽰使⽤验证码登录	输⼊2:表⽰获取验证码")
			fmt.Scanln(&order)
			switch order {
			case 1:
				var m string
				fmt.Println("输入验证码")
				fmt.Scanln(&m) //输入验证码
				code := service.(*CodeSend).password[phone]
				if m == code {
					if service.JudgeCode(phone, m) {
						fmt.Println("登录成功")
						service.Clear(phone)
						return
					}
					if !service.JudgeCode(phone, m) {
						fmt.Println("超时")
					}
				} else {
					fmt.Println("请输入正确验证码")
				}
			case 2:

				if err := service.CanSendCode(phone); err != nil { //验证是否可发送验证码，并输出错误结果
					fmt.Println("发送失败，", err)
				} else {
					//生成随机数
					result := service.Random()
					fmt.Printf("验证码：%v\n", result)
					service.(*CodeSend).password[phone] = result
				}
			default:
				fmt.Println("无效操作")
			}
		}
	}
}
