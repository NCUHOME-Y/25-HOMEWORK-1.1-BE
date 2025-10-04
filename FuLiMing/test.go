package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

type phoneVerification struct {
	phonenumber string
	choose int
	codenumber int
	x int
	code int
	contacts1 map[string]int
	contacts2 map[string]time.Time
	contacts3 map[string]int
}

func NewPhoneVerification() *phoneVerification {
	return &phoneVerification{
		contacts1: make(map[string]int),
		contacts2: make(map[string]time.Time),
		contacts3: make(map[string]int),
	}
}

func main() {
	verification := NewPhoneVerification()
	verification.Start()
}

func (pv *phoneVerification) Start() {
	fmt.Printf("请输入手机号码：")
	fmt.Scanln(&pv.phonenumber)
	
	if pv.contacts3[pv.phonenumber] == 0 {
		pv.contacts3[pv.phonenumber] = 5
	}
	
	pv.testPhoneNumber()
}

func (pv *phoneVerification) testPhoneNumber() {
	patern := `^1[3456789]\d{9}$`
	exam := regexp.MustCompile(patern)
	
	if exam.MatchString(pv.phonenumber) {
		pv.x = 0
		pv.test()
	} else {
		fmt.Println("请输入正确形式的电话号码。")
		pv.Start()
	}
}

func (pv *phoneVerification) test() {
	fmt.Println("1：请输入六位验证码进行登录  2：获取验证码")
	fmt.Scanln(&pv.choose)
	
	switch pv.choose {
	case 1:
		pv.choose1()

	case 2:
		pv.choose2()

	default:
		fmt.Println("内容错误")
		pv.test()
	}
}

func (pv *phoneVerification) choose1() {
	fmt.Println("请输入验证码")
	fmt.Scanln(&pv.code)
	
	if pv.code == pv.contacts1["验证码"] {
		fmt.Println("登录成功！")
		fmt.Println("验证码已失效")
		delete(pv.contacts1, "验证码")
	} else {
		fmt.Println("登录失败")
		pv.test()
	}
}

func (pv *phoneVerification) choose2() {
	now := time.Now()
	
	if pv.x == 0 {
		pv.x++
		pv.generateVerificationCode()
	} else if savedTime, exists := pv.contacts2["savetime"]; exists {
		t1 := savedTime.Add(1 * time.Minute)
		t2 := savedTime.Add(5 * time.Minute)
		
		if now.Before(t1) {
			fmt.Println("一分钟内无法再次获取验证码")
			pv.test()
			return
		} else if now.Before(t2) {
			pv.generateVerificationCode()
		} else {
			fmt.Println("验证码超过5分钟，已删除")
			delete(pv.contacts1, "验证码")
			pv.x--
			pv.test()
		}
	} else {
		pv.generateVerificationCode()
	}
}

func (pv *phoneVerification) generateVerificationCode() {
	if pv.contacts3[pv.phonenumber] > 0 {
		var codearr [6]int
		
		for i := 0; i < 6; i++ {
			codearr[i] = rand.Intn(10)
		}
		
		pv.codenumber = codearr[0]*100000 + codearr[1]*10000 + codearr[2]*1000 + 
			codearr[3]*100 + codearr[4]*10 + codearr[5]
		
		pv.contacts1["验证码"] = pv.codenumber
		fmt.Printf("验证码为：%06d\n", pv.contacts1["验证码"])
		
		pv.contacts3[pv.phonenumber]--
		
		nowTime := time.Now()
		pv.contacts2["savetime"] = nowTime
		
		pv.test()
	} else {
		fmt.Println("一天内最多只能获取5次验证码哦")
		pv.test()
	}
}