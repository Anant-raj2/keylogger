package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

type Server struct {
  peers map[*websocket.Conn]bool
  reader *bufio.Reader
  writer *bufio.Writer
}

func NewServer() *Server{
  return &Server{
    peers: map[*websocket.Conn]bool{},
  }
}

func (server *Server) HandleSocket(conn *websocket.Conn){
  server.peers[conn] = true
  server.reader = bufio.NewReader(conn)
  server.writer = bufio.NewWriter(conn)
  fmt.Println("[Server] Connection from: ", conn.RemoteAddr().String())
  server.ReadStream(conn)
}

func (server *Server) ReadStream(conn *websocket.Conn){
  buf:= make([]byte, 1024)
  for{
    n, err := server.reader.Read(buf)
    if err != nil {
      if err==io.EOF{
        fmt.Println("Reached end of file")
        break
      }
      fmt.Println("[Server] Error Reading from stream: ", err)
      continue
    }
    data:=buf[:n]
    fmt.Println("Recieved message from: ", conn.RemoteAddr().String())
    _, err=server.writer.Write(data)
    if err != nil {
      fmt.Println("[Server] Error writing: ", err)
      continue
    }
    server.writer.Flush()
    fmt.Println("Wrote message to: ", conn.RemoteAddr().String())
  }
}

func main(){
  var server *Server = NewServer()
  http.Handle("/log", websocket.Handler(server.HandleSocket))
  http.ListenAndServe(":3000", nil)
}
