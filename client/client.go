package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"golang.org/x/net/websocket"
)

type Client struct {
	conn   *websocket.Conn
	reader *bufio.Reader
	writer *bufio.Writer
	url    string
	origin string
  quitch chan struct{}
}

func NewClient(url string, origin string) *Client {
	client := &Client{
		url:    url,
		origin: origin,
    quitch: make(chan struct{}),
	}
	return client
}

func (client *Client) DialClient() error {
	ws, err := websocket.Dial(client.url, "", client.origin)
	if err != nil {
		return err
	}
  fmt.Println("[Client] Connected to: ", ws.RemoteAddr().String())
	client.conn = ws
	client.reader = bufio.NewReader(ws)
	client.writer = bufio.NewWriter(ws)
  go client.ReadStream()
  go client.WriteStream()
  <-client.quitch
	return nil
}

func (client *Client) ReadStream(){
  defer client.conn.Close()
  buf := make([]byte, 1024)
  for{
    n, err:=client.reader.Read(buf)
    if err != nil {
      if err == io.EOF{
        fmt.Println("[Client] reached end of stream")
        return
      }
      fmt.Println("[Client] Error reading stream: ", err)
      continue
    }
    fmt.Println(string(buf[:n]))
  }
}

func (client *Client) WriteStream(){
  scanner:= bufio.NewScanner(os.Stdin)
  defer client.conn.Close()
  for scanner.Scan(){
    _, err:=client.writer.Write(scanner.Bytes())
    if err != nil {
      fmt.Println("[Client] Error writing stream: ", err)
      return
    }
    err=client.writer.Flush()
    if err != nil {
      fmt.Println("[Client] Error flushing stream: ", err)
      return
    }
  }
}

func main(){
  var client *Client = NewClient("ws://localhost:3000/log", "http://localhost/")
  client.DialClient()
}
