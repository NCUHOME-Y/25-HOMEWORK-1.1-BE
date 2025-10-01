package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

func main() {
	var c phonenumbers
	c.Cin()
	c.Cout()
}

type phonenumbers struct {
	Number    string
	Time1     time.Time
	Timejugde bool
	conuts    int
}

type Cincheckphone interface {
	Cin()
	Check() bool
	Cout()
}
type Randome interface {
	Random() string
}
type Time101 interface {
	Time() bool
	Date() bool
}

// 输入手机号
func (c *phonenumbers) Cin() {
	c.Number = "-1"
	c.conuts = 0
	fmt.Println("请输入手机号:")
	pattern := `^1[3-9]\d{9}$`
	var b string
	fmt.Scan(&b)
	re, _ := regexp.Compile(pattern)
	Matched := re.MatchString(b)
	if Matched {
		return
	} else {
		fmt.Println("您输入号码有误,请重试....")
		c.Cin()
	}
}

// 生产验证码
func (c *phonenumbers) Random() string {
	c.Number = fmt.Sprintf("%06d", rand.Intn(100000))
	c.Time1 = time.Now()
	c.Timejugde = true
	c.conuts++
	return c.Number
}

// 判断时间
func (c *phonenumbers) Time() bool {
	duration := time.Since(c.Time1)
	if duration < 1*time.Minute {
		c.Timejugde = false
	}
	return c.Timejugde
}

// 判断次数
func (c *phonenumbers) Date() bool {
	duration := time.Since(c.Time1)
	if duration < 24*time.Hour && c.conuts >= 5 {
		c.Timejugde = false
	}
	return c.Timejugde
}

// 验证输入的验证码
func (c *phonenumbers) Check() bool {
	var b string
	fmt.Scan(&b)
	patten1 := `\d{6}$`
	re, _ := regexp.Compile(patten1)
	Matched := re.MatchString(b)
	if Matched {
		return c.Number == b
	} else {
		return false
	}
}

// 判断输出验证码程序
func (c *phonenumbers) Cout() {
	var choice int
	fmt.Println("验证码登录请按1  请求验证码请按2")
	fmt.Scan(&choice)
	switch choice {
	case 1:
		{
			fmt.Println("请输入验证码:")
			if c.Check() {
				fmt.Println("登录成功")
			} else {
				fmt.Println("无效验证码")
				c.Cout()
			}
		}
	case 2:
		{
			if c.conuts == 0 {
				fmt.Println(c.Random())
			} else if c.Time() && c.Date() {
				fmt.Println(c.Random())
			} else {
				fmt.Println("请稍后再试.....")
			}
			c.Cout()
		}
	default:
		{
			fmt.Println("请输入1，2")
			c.Cout()
		}
	}
}
