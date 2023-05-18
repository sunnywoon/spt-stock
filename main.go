package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

type Product struct {
	Name  string `json:"name"`
	Num   int    `json:"num"`
	Price string `json:"price"`
	Link  string `json:"link"`
}

type Msg struct {
	E3List  []Product `json:"e3_list"`
	E5List  []Product `json:"e5_list"`
	AmdList []Product `json:"amd_list"`
}

var proxy string

func main() {

	godotenv.Load()

	var sleep string
	var proxyStr string

	flag.StringVar(&proxyStr, "proxy", "", "HTTP proxy address (e.g. http://proxy.example.com:8080)")

	flag.StringVar(&sleep, "sleep", "60", "sleep second")

	flag.Parse()

	sleepNum, _ := strconv.Atoi(sleep)

	if len(proxyStr) > 0 {
		proxy = proxyStr
	}

	for {
		check()
		time.Sleep(time.Second * time.Duration(sleepNum))
	}
}

func check() {
	listAmd := checkAmd()
	listE5 := checkE5()
	listE3 := checkE3()

	var msg Msg

	if len(listAmd) > 0 {
		msg.AmdList = listAmd
	}

	if len(listE5) > 0 {
		msg.E5List = listE5
	}

	if len(listE3) > 0 {
		msg.E3List = listE3
	}

	msgJson, _ := json.Marshal(msg)

	fmt.Println(string(msgJson))

	// 通知webhook
	url := os.Getenv("WEBHOOK_URL")

	sendJson(url, msgJson)
}

func sendJson(url string, json []byte) {
	// 创建一个 POST 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
	if err != nil {
		fmt.Println("Request creation error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求并获取响应
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Request failed. Status:", resp.StatusCode)
		return
	}
}

func setProxy(c *colly.Collector) {
	if len(proxy) > 0 {
		c.SetProxy(proxy)
	}
}

func checkE5() []Product {
	collector := colly.NewCollector(
		colly.AllowedDomains("billing.spartanhost.net"),
	)

	setProxy(collector)

	list := make([]Product, 0)
	collectProduct(collector, &list)

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	collector.Visit("https://billing.spartanhost.net/store/ddos-protected-ssd-e5-kvm-vps-seattle")

	return list
}

func checkE3() []Product {
	collector := colly.NewCollector(
		colly.AllowedDomains("billing.spartanhost.net"),
	)

	setProxy(collector)

	list := make([]Product, 0)
	collectProduct(collector, &list)

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	collector.Visit("https://billing.spartanhost.net/store/ddos-protected-hdd-e3-kvm-vps-seattle")

	return list
}

func checkAmd() []Product {
	collector := colly.NewCollector(
		colly.AllowedDomains("billing.spartanhost.net"),
	)

	setProxy(collector)

	list := make([]Product, 0)
	collectProduct(collector, &list)

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	collector.Visit("https://billing.spartanhost.net/store/ddos-protected-ssd-premium-kvm-vps-seattle")

	return list
}

func collectProduct(collector *colly.Collector, list *[]Product) {
	collector.OnHTML("div.products", func(h *colly.HTMLElement) {
		h.ForEach("div.product", func(_ int, h *colly.HTMLElement) {
			product := Product{}
			h.ForEach("span", func(i int, h *colly.HTMLElement) {
				if i == 0 {
					product.Name = h.Text
				}

				if i == 1 {
					qty := h.Text
					qty = strings.Replace(qty, "Available", "", -1)
					qty = strings.Trim(qty, "")
					num, err := strconv.Atoi(qty)
					if err != nil {
						num = 0
					}

					product.Num = num
				}
			})

			price := h.ChildText("span.price")
			link := h.ChildAttr("a", "href")

			product.Price = price
			product.Link = link

			if product.Num > 0 {
				*list = append(*list, product)
			}
		})

	})
}
