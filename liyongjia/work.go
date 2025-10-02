package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type CodeSend struct {
	lastphonetime map[string]time.Time //时间
	mutex         sync.Mutex           //锁
	password      map[string]int       //存储验证码
	sendCont      map[string]int       //发送次数

}

// 初始化:用常量赋予结构体变量
func CodeSendInit() *CodeSend {
	return &CodeSend{
		lastphonetime: make(map[string]time.Time),
		password:      make(map[string]int),
		sendCont:      make(map[string]int),
	}
}

// 接口
type CodeServier interface {
	CanSendCode(phone string) error
	JudgeCode(phone string, m int) bool
	Clear(phone string)
	Random(s []int)
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
	c.mutex.Lock()
	defer c.mutex.Unlock()
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
func (c *CodeSend) JudgeCode(phone string, m int) bool {
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
func (c *CodeSend) Random(s []int) {
	for i := 0; i < len(s); i++ {
		//随机整数
		s[i] = rand.Intn(10)
	}
}

func main() {
	var service CodeServier = CodeSendInit()

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
			var result2 int
			var order int
			fmt.Println("输⼊1:表⽰使⽤验证码登录	输⼊2:表⽰获取验证码")
			fmt.Scanln(&order)
			switch order {
			case 1:
				var m int
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
					n := 6
					s := make([]int, n)
					service.Random(s)

					//类型转化  []int>>str>>int
					var str string

					for _, date := range s {
						str += strconv.Itoa(date) // 遍历切片并拼接字符串strconv.Itoa(date)：将整数 date 转换为字符串str
					}
					result2, _ = strconv.Atoi(str) //strconv.Atoi(str)：将字符串str转换回整数result2 Atoi 是 "ASCII to Integer" 的缩写。因为一定为整数类型所以省去了err
					fmt.Printf("验证码：")
					fmt.Println(result2)
					service.(*CodeSend).password[phone] = result2
				}
			default:
				fmt.Println("无效操作")
			}
		}
	}
}
