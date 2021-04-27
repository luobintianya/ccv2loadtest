package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	totalstart := time.Now()
	done := make(chan bool, 1)
	for i := 0; i < 200; i++ {
		go postData(i, client, done)
		time.Sleep(time.Second / time.Duration(rand.Intn(9)+1))
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			_, cancel := context.WithTimeout(context.Background(), 35*time.Second)
			defer cancel()

			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}

	totalend := time.Now()

	spec := float32((totalend.Unix() - totalstart.UnixNano()))
	<-done
	time.Sleep(2 * time.Second)
	fmt.Printf("total  %f 秒", spec)
}

func postData(thread int, client *http.Client, done chan bool) float32 {
	var totalTime *float32=new( float32)



	for i := 0; i < 10000; i++ {
		doRequst(i, thread, client, done, totalTime)
	}

	return *totalTime
}

func doRequst(i int, thread int, client *http.Client, done chan bool, totalTime *float32) {

	r := rand.Int63n(9) + 1
	start := time.Now()
	order := "{\"code\":\"111234\",\"name\":\"111234\",\"description\":\"abc\",\"net\":false,\"entries\":[{\"quantity\":10,\"product\":{\"code\":\"4567130\"},\"updateable\":false,\"addable\":false},{\"quantity\":10,\"product\":{\"code\":\"4567133\"},\"updateable\":false,\"addable\":false}]}\""
	req := &http.Request{}
	//req.PostForm = url.Values{"order": {order}}
	req.Header = http.Header{}
	req.Method = "POST"
	//req.SetBasicAuth("eic", "1234")
	//req.Header.Add("authorization", "Bearer A4gNnHJlTeMGl_a6DrcnfGzvk2M")
	req.Header.Add("authorization", "Bearer 6trBPMPuXMjAgHSmR0YMHAVX3Y0")
	//req.Header.Add("Content-Type","application/x-www-form-urlencoded")
	//req.URL, _ = url.Parse("https://api.cg8cr5xpc9-xxxhold1-d1-public.model-t.ccv2prod.sapcloud.cn/occ/v2/powertools/order/poc")
	//req.URL, _ = url.Parse("https://localhost:9002/occ/v2/powertools/order/poc?order="+url.QueryEscape(order))
	//req.URL, _ = url.Parse("https://api.cg8cr5xpc9-xxxhold1-d1-public.model-t.ccv2prod.sapcloud.cn/occ/v2/powertools/order/poc?order=" + url.QueryEscape(order))
	req.URL, _ = url.Parse("https://api.cg8cr5xpc9-xxxhold1-s1-public.model-t.ccv2prod.sapcloud.cn/occ/v2/powertools/order/poc?order=" + url.QueryEscape(order))
	req.ParseForm()
	res, err := client.Do(req)
	if (err != nil || res.StatusCode != 200) && res != nil {
		//fmt.Println(err)
		bodyx, _ := ioutil.ReadAll(res.Body)
		str := string(bodyx)
		fmt.Println(str)
		done <- true

	}
	end := time.Now()
	spec := float32((end.UnixNano() - start.UnixNano())) / 1000000000.0
	*totalTime = *totalTime + spec

	time.Sleep(2 * (time.Second / time.Duration(r)))
	if res != nil {
		fmt.Printf("协程号 %2d , 新建订单编号 %4d, spec:%f 秒 , 服务器响应Code: %d \n", thread, i, spec, res.StatusCode)
	} else {
		fmt.Printf("协程号 %2d , 新建订单编号 %4d, spec:%f 秒 , 服务器响应错误Code: %s \n", thread, i, spec, err)
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("错误%s\n　继续！", err)
		}
	}()
}
