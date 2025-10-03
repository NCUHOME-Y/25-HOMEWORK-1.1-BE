package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

// 存入data.json
func (s *Map) SaveToFile(filename string) error {
	file, err1 := os.Create(filename)
	if err1 != nil {
		return err1
	}
	err2 := json.NewEncoder(file).Encode(s)
	if err2 != nil {
		return err2
	}
	defer file.Close()
	return nil
}

// 启动读取数据
func LoadFromFile(filename string) (Map, error) {
	file, err1 := os.Open(filename)
	s := make(Map)
	if err1 != nil {
		return s, err1
	}
	err := json.NewDecoder(file).Decode(&s)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return s, err
}

func IsPhone(number string) error {
	phone := `^1[3-9]\d{9}$`
	reg := regexp.MustCompile(phone)
	if !reg.MatchString(number) {
		return errors.New("请输入正确的手机号")
	} else {
		return nil
	}
}
func GetCode() string {
	q := `1234567890QWERTYUIOPASDFGHJKLZXCVBNMqwertyuioplkjhgfdsazxcvbnm`
	var y = make([]string, 0)
	for r := 0; r < 6; r++ {
		idx := rand.Intn(len(q))
		y = append(y, string(q[idx]))
	}
	return strings.Join(y, "")
}

type CodeRecord struct {
	Code     string
	FormTime time.Time
	Used     bool
	Times    int
	Date     string
}
type Map map[string]*CodeRecord

func (s *Map) AddCode(number string) (string, error) {
	now := time.Now()
	today := now.Format("2006-01-02")
	d, ok := (*s)[number]
	if !ok {
		code := GetCode()
		(*s)[number] = &CodeRecord{
			Code:     code,
			FormTime: now,
			Used:     false,
			Times:    1,
			Date:     today,
		}
		return code, nil
	}
	if d.Date != today {
		d.Times = 0
		d.Date = today
	}
	if d.Times == 5 {
		return "", errors.New("今日发送次数已达上限")
	}
	if now.Sub(d.FormTime) < 60*time.Second {
		return "", errors.New("获取太频繁，请稍后再试")
	}
	code := GetCode()
	d.Code = code
	d.FormTime = now
	d.Used = false
	d.Times++
	return code, nil

}
func (s *Map) TestCode(number, code string) error {
	now := time.Now()
	d, ok := (*s)[number]
	if !ok {
		return errors.New("请先获取验证码")
	}
	if d.Code != code || d.Used || now.Sub(d.FormTime) > 5*time.Minute {
		return errors.New("验证码错误")
	}
	d.Used = true
	return nil
}

func main() {
	m := make(Map)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n请输入手机号：")

		phone, _ := reader.ReadString('\n')
		phone = strings.TrimSpace(phone)
		err := IsPhone(phone)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for {
			fmt.Println("请输入操作：1=使用验证码登录  2=获取验证码")
			t, _ := reader.ReadString('\n')
			t = strings.TrimSpace(t)
			switch t {
			case "2":
				code, err := m.AddCode(phone)
				if err != nil {
					fmt.Println("获取失败：", err.Error())
					continue
				} else {
					fmt.Println(code)
					continue
				}

			case "1":
				fmt.Println("请输入验证码：")
				g, _ := reader.ReadString('\n')
				g = strings.TrimSpace(g)
				if err := m.TestCode(phone, g); err != nil {
					fmt.Println("登陆失败：", err.Error())
					continue
				} else {
					fmt.Println("登录成功！")
				}
			}
			break
		}
		break
	}

}
