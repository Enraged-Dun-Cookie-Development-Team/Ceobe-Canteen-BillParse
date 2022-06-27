package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"

	iconv "github.com/djimenez/iconv-go"
)

func main() {
	ReadWechat()
	ReadAlipay()
	WriteWechat("wechat")
	WriteWechat("alipay")

}

var wechatList = make([]Record, 0)
var alipayList = make([]Record, 0)

type Record struct {
	Time        string
	Amount      string
	Person      string
	Description string
	BillNumber  string
}

func ReadAlipay() {
	opencast, err := os.Open("./alipay.csv")
	if err != nil {
		log.Println("支付宝表格不存在")
		return
	}
	defer opencast.Close()

	log.Println("开始读取支付宝表格数据")
	reader := csv.NewReader(opencast)
	i := 0
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if i < 3 || len(line) < 22 {
			i++
			continue
		}
		if !strings.HasPrefix(line[10], "刻") || line[2] != "收钱码收款" {
			continue
		}

		time, _ := iconv.ConvertString(line[1], "gb2312", "utf-8")
		amount, _ := iconv.ConvertString(line[7], "gb2312", "utf-8")
		person, _ := iconv.ConvertString(line[6], "gb2312", "utf-8")
		descrition, _ := iconv.ConvertString(line[22], "gb2312", "utf-8")
		billNumber, _ := iconv.ConvertString(line[4], "gb2312", "utf-8")

		if strings.HasPrefix(person, "*") {
			person = "匿名"
		}

		alipayInfo := Record{
			Time:        time,
			Amount:      amount,
			Person:      person,
			Description: descrition,
			BillNumber:  billNumber,
		}
		alipayList = append(alipayList, alipayInfo)

	}
	log.Println("支付宝表格数据读取完成")
}

func ReadWechat() {
	opencast, err := os.Open("./wechat.csv")
	if err != nil {
		log.Println("微信表格不存在")
		return
	}
	defer opencast.Close()
	log.Println("开始读取微信表格数据")
	reader := csv.NewReader(opencast)
	i := 0
	for {

		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if i < 17 {
			i++
			continue
		}
		if !strings.HasPrefix(line[10], "刻") || line[1] != "二维码收款" {
			continue
		}

		wechatInfo := Record{
			Time:        line[0],
			Amount:      strings.Trim(line[5], "¥"),
			Person:      line[2],
			Description: strings.Trim(line[10], "/"),
			BillNumber:  strings.Trim(line[8], "\t"),
		}
		wechatList = append(wechatList, wechatInfo)
	}
	log.Println("微信表格数据读取完成")
}

func WriteWechat(source string) {
	file := "./handle_alipay.csv"
	if source == "wechat" {
		file = "./handle_wechat.csv"
	}
	//读取文件，自动创建，覆盖写入
	File, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Println("文件打开失败！")
	}
	defer File.Close()

	WriterCsv := csv.NewWriter(File)

	if source == "wechat" {
		for _, v := range wechatList {
			line := []string{
				v.Time,
				v.Amount,
				v.Person,
				v.Description,
				v.BillNumber,
			}
			err := WriterCsv.Write(line)
			if err != nil {
				continue
			}
		}
	} else {
		for _, v := range alipayList {
			line := []string{
				v.Time,
				v.Amount,
				v.Person,
				v.Description,
				v.BillNumber,
			}
			err := WriterCsv.Write(line)
			if err != nil {
				continue
			}
		}
	}
	WriterCsv.Flush() //刷新，不刷新是无法写入的
	log.Println("数据写入成功...")
}
