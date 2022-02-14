package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var host string
var path string
var port string
var archivos int

func main() {
	archivos = 0
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
	response, connection := handleConnection()
	createHtml(response)
	links := getLinks(response)
	getResources(links, connection)
}

func createHtml(doc string) {
	exp := regexp.MustCompile("text/html")
	if string(exp.Find([]byte(doc))) != "" {
		doc = removeHeader(doc)
		file, err := os.Create("index.html")
		checkError(err)
		defer file.Close()
		ioutil.WriteFile("index.html", []byte(doc), 0644)
	} else {
		doc = removeHeader(doc)
		var name string = ""
		idx := strings.LastIndex(path, ".")
		if idx == -1 {
			name = strconv.Itoa(archivos)
		} else {
			name = strconv.Itoa(archivos) + path[idx:]
		}
		archivos += 1
		file, err := os.Create(name)
		checkError(err)
		defer file.Close()
		ioutil.WriteFile(name, []byte(doc), 0644)
		archivos += 1
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

func removeHeader(doc string) string {
	idx := strings.Index(doc, "\r\n\r")
	if idx != -1 {
		fmt.Println("")
		fmt.Println(doc[:idx])
		fmt.Println("")
		return doc[idx+4:]
	} else {
		return doc
	}
}

func downloadResource(link string, done chan bool) {
	var name string = ""
	name = strings.ReplaceAll(link, "/", "-")

	archivos += 1
	tcpAddr, err := net.ResolveTCPAddr("tcp4", host+":"+port)
	checkError(err)
	connection, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	defer connection.Close()
	_, err = connection.Write([]byte("GET " + link + " \r\n\r\n"))
	if err != nil {
		return
	}
	response, err := ioutil.ReadAll(connection)
	checkError(err)
	//fmt.Println(string(response))

	file, err := os.Create(name)
	checkError(err)
	defer file.Close()

	idxDoc := strings.Index(string(response), "\n")
	if idxDoc != -1 {
		doc := string(response)[idxDoc+2:]
		doc = removeHeader(doc)
		ioutil.WriteFile(name, []byte(doc), 0644)
	}
	done <- true
}
func getLinks(response string) []string {
	links := regexp.MustCompile("(src=\".*?\")|(src='.*?')")

	ls := links.FindAllString(response, -1)

	out := []string{}
	for _, l := range ls {
		l := strings.Trim(l, " ")
		l = strings.Replace(l, "src=", "", 1)
		l = strings.Replace(l, "\"", "", 2)
		l = strings.ReplaceAll(l, "'", "")
		l = strings.ReplaceAll(l, ":", "")
		out = append(out, l)
	}
	for _, l := range out {
		fmt.Println(string(l))
	}

	return out
}
func handleConnection() (string, net.Conn) {

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
