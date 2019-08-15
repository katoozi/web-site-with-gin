package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

// RunWebSocketCommand is a cobra command that start web server
var RunWebSocketCommand = &cobra.Command{
	Use:   "runwebsocket",
	Short: "start websocket server",
	Run:   runWebSocket,
}

// Packet represents application level data.
type Packet struct {
	ActionType  string `json:"action_type"`
	Text        string `json:"text"`
	MessageType int
}

// Channel wraps user connection.
type Channel struct {
	conn *websocket.Conn // WebSocket connection.
	send chan Packet     // Outgoing packets queue.
}

func (ch *Channel) reader() {
	for {
		messageType, p, err := ch.conn.ReadMessage()
		if err != nil {
			return
		}
		pkt := readPacket(p, messageType)
		ch.send <- *pkt
	}
}

func (ch *Channel) writer() {
	for pkt := range ch.send {
		data, err := json.Marshal(pkt)
		if err != nil {
			log.Printf("Error while marshaling json: %v", err)
		}
		if err := ch.conn.WriteMessage(pkt.MessageType, data); err != nil {
			return
		}
	}
}

// NewChannel is the Channel factorty function
func NewChannel(conn *websocket.Conn) *Channel {
	c := &Channel{
		conn: conn,
		send: make(chan Packet, 10),
	}

	go c.reader()
	go c.writer()

	return c
}

func readPacket(data []byte, messageType int) *Packet {
	packetData := &Packet{}
	err := json.Unmarshal(data, packetData)
	if err != nil {
		log.Printf("Error while unmarshal json: %v", err)
	}
	packetData.MessageType = messageType
	return packetData
}

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// We'll need to check the origin of our connection
	// this will allow us to make requests from our React
	// development server to here.
	// For now, we'll do no checking and just allow any connection
	// in production must be nil
	// in develop must be `func(r *http.Request) bool { return true }`
	CheckOrigin: func(r *http.Request) bool { return true },
}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
// func reader(conn *websocket.Conn) {
// 	for {
// 		// read in a message
// 		messageType, p, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}
// 		// print out that message for clarity
// 		fmt.Println(string(p))

// 		if err := conn.WriteMessage(messageType, p); err != nil {
// 			log.Println(err)
// 			return
// 		}

// 	}
// }

// define our WebSocket endpoint
func serveWs(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error while upgrade connection: %v", err)
		return
	}
	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	// ch := NewChannel(ws)
	NewChannel(ws)
	// reader(ws)
}

func httpFileHandler(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "./web/build/index.html")
}

func setupRoutes() {
	// mape our `/ws` endpoint to the `serveWs` function
	http.HandleFunc("/ws", serveWs)
	fs := http.FileServer(http.Dir("./web/build/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", httpFileHandler)
}

func runWebSocket(cmd *cobra.Command, args []string) {
	fmt.Println("Start to Listen...")
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}