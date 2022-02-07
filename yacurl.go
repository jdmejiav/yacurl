package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strings"
)

var host string
var path string
var port string

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s host port\n\n", os.Args)
		os.Exit(1)
	}

	idx := strings.Index(os.Args[1], "/")
	port = os.Args[2]
	if idx != -1 {
		host = os.Args[1][:idx]
		path = os.Args[1][idx:]
	} else {
		host = os.Args[1]
		path = "/"
	}
	fmt.Println("port " + port)
	fmt.Println("host " + host)
	fmt.Println("path " + path)
	response, connection := listener()
	fmt.Println(response)
	createHtml(response)
	links := getLinks(response)
	getResources(links, connection)

}
func createHtml(doc string) {
	idx := strings.Index(doc, "<")
	if idx != -1 {
		f, err := os.Create("index.html")
		checkError(err)
		defer f.Close()
		ioutil.WriteFile("index.html", []byte(doc[idx+1:]), 0644)
	}
}
func getResources(links []string, connection net.Conn) {
	done := make(chan bool)
	for _, l := range links {
		go downloadResource(l, done)
	}
	for range links {
		<-done
	}

}
func downloadResource(link string, done chan bool) {
	var name string = ""
	idx := strings.LastIndex(link, "/")
	if idx == -1 {
		name = link
	} else {
		name = link[idx+1:]
	}
	file, err := os.Create(name)
	checkError(err)
	defer file.Close()
	tcpAddr, err := net.ResolveTCPAddr("tcp4", host+":"+port)
	checkError(err)
	connection, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	defer connection.Close()

	_, err = connection.Write([]byte("GET " + link + " \r\n\r\n"))
	checkError(err)
	response, err := ioutil.ReadAll(connection)
	checkError(err)
	fmt.Println(string(response))
	idxDoc := strings.Index(string(response), "\n")
	if idxDoc != -1 {
		doc := string(response)[idxDoc+2:]
		ioutil.WriteFile(name, []byte(doc), 0644)
	}
	done <- true
}
func getLinks(response string) []string {

	links := regexp.MustCompile("src *=\".*?\"")

	ls := links.FindAllString(response, -1)
	out := []string{}
	for _, l := range ls {
		l = strings.Replace(l, "src=", "", 1)
		l = strings.Replace(l, "\"", "", 2)
		out = append(out, l)
	}

	return out
}
func listener() (string, net.Conn) {

	tcpAddr, err := net.ResolveTCPAddr("tcp4", host+":"+port)
	checkError(err)
	connection, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	defer connection.Close()
	_, err = connection.Write([]byte("GET " + path + " HTTP/1.0\r\n\r\n"))
	checkError(err)
	response, err := ioutil.ReadAll(connection)
	checkError(err)
	return string(response), connection
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
