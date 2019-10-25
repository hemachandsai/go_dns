package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wait sync.WaitGroup
var start = time.Now()
var cnameRecord string
var noRecord string
var nxDomain string
var space = "                                                "
var count int
var length int

func main() {
	// array := []string{"greenhouse.io"}
	array := []string{"berush.com", "bitmoji.com", "bitstrips.com", "events.semrush.com", "gnip.com", "greenhouse.io", "hacker101.com", "hackerone-ext-content.com", "hackerone-user-content.com", "hackerone.com", "hackerone.net", "istarbucks.co.kr", "labs-semrush.com", "legalrobot.com", "mobpub.com", "onelogin.com", "paypal.com", "periscope.tv", "pscp.tv", "semrush.com", "shipt.com", "slack-files.com", "slack-imgs.com", "slack-redir.net", "slack.com", "slackatwork.com", "slackb.com", "spaces.pm", "starbucks.ca", "starbucks.co.jp", "starbucks.co.uk", "starbucks.com", "starbucks.com.br", "starbucks.com.cn", "starbucks.com.sg", "starbucks.de", "starbucks.fr", "starbucksreserve.com", "twimg.com", "twitter.com", "uber.com", "uber.com.cn", "ubunt.com", "ui.com", "vine.co"}
	for index := range array {
		fmt.Println("Doing task for domain: " + array[index] + " index: " + strconv.Itoa(index))
		cnameRecord = ""
		noRecord = ""
		nxDomain = ""
		count = 0
		readFile(array[index])
	}
}
func readFile(domain string) {
	content, err := ioutil.ReadFile("C:\\Users\\QDS\\Desktop\\misc\\dns\\output\\" + domain + "_result.txt")
	if err == nil {
		contentString := string(content)
		data := strings.Split(contentString, "\r\n")
		length = len(data)
		fmt.Println(length)
		for ele := range data {
			wait.Add(1)
			go resolveCNAME(data[ele])
		}
	} else {
		fmt.Println(err)
		return
	}
	wait.Wait()
	writeToFile(domain)
}
func returnSpace(length int) string {
	space := ""
	for i := 0; i < 50-length; i++ {
		space += " "
	}
	return space
}
func resolveCNAME(dname string) {
	cname, err := net.LookupCNAME(dname)
	if err == nil {
		go resolveIP(cname, dname)
	} else {
		noRecord += dname + returnSpace(len(dname)) + "\n"
	}
	wait.Done()
	count++
	if count%100 == 0 {
		fmt.Println((float64(count) / float64(length)) * float64(100))
	}
}

func resolveIP(cname, dname string) {
	wait.Add(1)
	ip, err := net.LookupIP(cname)
	if err == nil {
		cnameRecord += dname + returnSpace(len(dname)) + cname + returnSpace(len(cname)) + ip[0].String() + "\n"
	} else {
		nxDomain += dname + returnSpace(len(dname)) + cname + returnSpace(len(cname)) + err.Error() + "\n"
	}
	wait.Done()
}

func writeToFile(domain string) {
	var finalContent string
	finalContent += "NXDOMAIN DATA\r\n\r\n"
	finalContent += nxDomain
	finalContent += "\r\n\r\nCNAME RECORD DATA\r\n\r\n"
	finalContent += cnameRecord
	finalContent += "\r\n\r\nNORECORD RECORD DATA\r\n\r\n"
	finalContent += noRecord
	ioutil.WriteFile("C:\\Users\\QDS\\Desktop\\misc\\dns\\output\\finalresult_"+domain+".txt", []byte(finalContent), 222)
	end := time.Since(start)
	fmt.Println("Time taken for domain: "+domain+" seconds: ", end)
}
