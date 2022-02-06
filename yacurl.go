package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port\n\n", os.Args)
		fmt.Println(os.Args)
		os.Exit(1)
	}
	response, connection := listener()

	links := getLinks(response)
	getResources(links, connection)

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
	_, err := os.Create(name)
	checkError(err)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", os.Args[1])
	checkError(err)
	connection, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	defer connection.Close()
	fmt.Println("GET " + link + " HTTP/1.0 \r\n\r\n")
	_, err = connection.Write([]byte("GET " + link + " HTTP/1.0 \r\n\r\n"))
	checkError(err)
	response, err := ioutil.ReadAll(connection)
	checkError(err)
	fmt.Println(string(response))
}
func getLinks(response string) []string {

	links := regexp.MustCompile("src=\".*?\"")

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

	url := os.Args[1]
	tcpAddr, err := net.ResolveTCPAddr("tcp4", url)
	checkError(err)
	connection, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	defer connection.Close()
	_, err = connection.Write([]byte("GET / HTTP/1.0\r\n\r\n"))
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
