package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/tencentyun/cos-go-sdk-v5"
)

func main() {
	//环境变量
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env文件加载错误！")
	}
	// 先预定义好COS相关东西
	// 明文存id与key无异于自杀行为
	// 请做好跑路准备
	u, _ := url.Parse(os.Getenv("Sync_Url"))
	w := os.Getenv("Sync_WebhookUrl")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("Sync_SecretID"),
			SecretKey: os.Getenv("Sync_SecretKey"),
		},
	})
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: "text/html",
		},
		ACLHeaderOptions: &cos.ACLHeaderOptions{
			// 如果不是必要操作，建议上传文件时不要给单个文件设置权限，避免达到限制。若不设置默认继承桶的权限。
			XCosACL: "private",
		},
	}
	// 然后读本地文件
	file, err := os.Open("./data")
	// 任何时候都得做好出错准备
	// 出错就得擦屁股
	if err != nil {
		log.Fatalf("打开文件夹失败: %s", err)
	}
	// 养成好习惯
	// 打开的门记得给人关上
	defer file.Close()
	// 拿出来了就读吧
	// 0 就是读所有文件文件夹
	list, _ := file.Readdirnames(0)
	count := 0
	for _, name := range list {
		_, err := c.Object.Head(context.Background(), name, nil)
		if err != nil {
			log.Print("数据不存在，开始上传...")
			log.Print(name)
			_, err = c.Object.PutFromFile(context.Background(), name, "./data/"+name, opt)
			if err != nil {
				panic(err)
			}
			count++
			log.Print("已上传.")
		} else {
			log.Print("数据存在，跳过...")
		}
	}
	// webhook
	postBody, _ := json.Marshal(map[string]string{
		"msgtype": "text",
		"text":    `{"content": "同步完成, 本次共计同步` + strconv.Itoa(count) + `条数据."}`,
	})
	responseBody := bytes.NewBuffer(postBody)
	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post(w, "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("webhook错误 %v", err)
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Printf(sb)
}
