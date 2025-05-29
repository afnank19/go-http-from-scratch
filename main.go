package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
)

// Step 1 Send Text over a TCP connection [DONE]
// Step 2 Probably parse it

type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	msg        chan []byte
	getPaths   map[string]string
}

type Headers struct {
	hostHeaderExists bool
	acceptType       string
	userAgent        string
}

type Request struct {
	method      string
	path        string
	httpVersion string
	headers     Headers
}

type Response struct {
	statusLine    string
	contentLength string
	contentType   string
	body          string
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msg:        make(chan []byte, 10),
		getPaths:   make(map[string]string),
	}
}

func (s *Server) Start() error {
	fmt.Println("Server up on: ", s.listenAddr)
	listener, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	s.ln = listener
	go s.acceptLoop()

	<-s.quitch
	close(s.msg)

	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("ERR: ", err)
			continue
		}

		fmt.Println("New Connection at: ", conn.RemoteAddr().String())

		go s.readLoop(conn)
	}
}

func (s *Server) readLoop(conn net.Conn) {
	// defer conn.Close()

	buffer := make([]byte, 2048)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			break
		}

		// fmt.Printf("Received: %s\n", buffer[:n])
		buff := buffer[:n]
		req := string(buff)
		// fmt.Println(string(buff))
		parsedReq := ParseHTTPRequest(req)
		response := handleRequest(parsedReq, s)
		res := buildResponse(response)

		conn.Write([]byte(res))

		// s.msg <- buff
	}
}

func (s *Server) Get(path string, filepath string) {
	s.getPaths[path] = filepath
}

func main() {
	fmt.Println("hello nano")
	fmt.Println("PID: ", os.Getpid())

	server := NewServer("localhost:8080")

	// go func() {
	// 	for msg := range server.msg {
	// 		fmt.Println(string(msg))
	// 	}
	// }()
	server.Get("/", "./html/index.html")
	server.Get("/about", "./html/about.html")
	server.Get("/style.css", "style.css")

	server.Start()
}

func handleRequest(parsedReq *Request, s *Server) Response {
	var response Response

	var contentType string = "text/html"
	fileExtension := filepath.Ext(parsedReq.path)
	if fileExtension != "" {
		contentType = getContentType(fileExtension)
	}

	filepath, exists := s.getPaths[parsedReq.path]
	if !exists {
		handleNotFound(&response)
		return response
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		handleNotFound(&response)
		return response
	}

	// fmt.Println(string(data))

	response.statusLine = "HTTP/1.1 200 OK\r\n"
	response.contentLength = fmt.Sprintf("Content-Length: %d\r\n", len(data))
	response.contentType = "Content-Type: " + contentType + "\r\n"
	response.body = string(data)

	return response
}

func buildResponse(response Response) string {
	res := response.statusLine + response.contentLength + response.contentType + "\r\n" + response.body
	return res
}

func handleNotFound(response *Response) {
	data, err := os.ReadFile("./notfound.html")
	if err != nil {
		fallBackData := "404 Not Found, this is automatically handled by the server, if you want a better UI, create a 'notfound.html' in the root dir"
		response.statusLine = "HTTP/1.1 404 NOT FOUND\r\n"
		response.contentLength = fmt.Sprintf("Content-Length: %d\r\n", len(fallBackData))
		response.contentType = "Content-Type: text/plain\r\n"
		response.body = string(fallBackData)

		return
	}

	response.statusLine = "HTTP/1.1 404 NOT FOUND\r\n"
	response.contentLength = fmt.Sprintf("Content-Length: %d\r\n", len(data))
	response.contentType = "Content-Type: text/html\r\n"
	response.body = string(data)
}

func getContentType(ext string) string {
	if ext == ".css" {
		return "text/css"
	}

	//default behaviour
	return "text/html"
}
