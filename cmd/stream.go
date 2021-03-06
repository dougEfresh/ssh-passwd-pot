// Copyright © 2017 Douglas Chimento <dchimento@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"
)

//var eventClient api.Transporter

// streamCmd represents the stream command
/*
const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512

	eventTable    = "event"
	eventGeoTable = "event_geo"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var db potdb.DB

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
*/
var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Stream events to websocket clients",
	Long:  "",
	Run:   runCobra,
}

func runCobra(cmd *cobra.Command, args []string) {
	/*
		var err error
		setupLogger(cmd.Name())
		setup(cmd, args)
		r := mux.NewRouter()
		r.HandleFunc(api.StreamURL, streamEvents).Methods("GET")
		r.HandleFunc(api.StreamURL+"/random", streamEvents).Methods("GET")

		srv := &http.Server{
			Handler:      r,
			Addr:         config.BindAddr,
			WriteTimeout: 10 * time.Second,
			ReadTimeout:  10 * time.Second,
		}

		logger.Infof("Listing on %s", config.BindAddr)
		db, err = potdb.Open(config.Dsn)
		if err != nil {
			logger.Errorf("Error loading DB %s", err)
			os.Exit(1)
		}
		eventClient, _ = event.New(event.SetDB(db))
		//websocket requests
		go hub.run()
		go randomDataHub.run()
		go startRandomHub(randomDataHub)
		err = srv.ListenAndServe()

		if err != nil {
			logger.Errorf("Caught error %s", err)
			os.Exit(-1)
		}
	*/
}

/*
// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				logger.Infof("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func streamEvents(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "random") {
		serveRandomWs(randomDataHub, w, r)
	} else {
		serveWs(hub, w, r)
	}
}

// serveWs handles websocket requests from the peer.
func serveRandomWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 512)}
	client.hub.register <- client
	go client.writePump()
	client.readPump()
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 512)}
	client.hub.register <- client
	go client.writePump()
	client.readPump()
}

var lastRandomEvent *api.EventGeo

func startRandomHub(hub *Hub) {
	logger.Info("Starting random hub")

	var id int64
	query := fmt.Sprintf("SELECT id FROM %s WHERE id != $1 and origin_country != '' and remote_latitude != $2 and remote_longitude != $3 and id >= (SELECT max(id) * RANDOM() FROM %s) ORDER BY id LIMIT 1 ", eventGeoTable, eventTable)
	for {
		time.Sleep(1250 * time.Millisecond)
		if len(hub.clients) == 0 {
			continue
		}

		if lastRandomEvent == nil {
			r := db.QueryRow("SELECT max(id) FROM event_geo")
			err := r.Scan(&id)
			if err != nil {
				logger.Error("Error getting max id")
				continue
			}
		} else {
			r := db.QueryRow(query, lastRandomEvent.ID, lastRandomEvent.RemoteLatitude, lastRandomEvent.RemoteLongitude)
			err := r.Scan(&id)
			if err != nil {
				logger.Errorf("Error getting next id %s", err)
			}
		}
		if id == 0 {
			logger.Error("Could not find an random event!")
		}
		if id != 0 {
			lastRandomEvent = broadcastEvent(id, hub)
		}
	}
}

func broadcastEvent(id int64, hub *Hub) *api.EventGeo {
	if len(hub.clients) == 0 {
		return nil
	}
	gEvent, _ := eventClient.GetEvent(id)
	if gEvent == nil {
		return nil
	}
	if b, err := json.Marshal(gEvent); err != nil {
		logger.Errorf("Error decoding geo event %d %s", id, err)
	} else {
		hub.broadcast <- b
	}
	return gEvent
}
*/
func init() {
	RootCmd.AddCommand(streamCmd)
	streamCmd.PersistentFlags().StringVar(&config.Dsn, "dsn", "postgres://postgres:@172.17.0.1/?sslmode=disable", "DSN database url")
	streamCmd.PersistentFlags().StringVar(&config.BindAddr, "bind", "localhost:8080", "bind to this address:port")
	streamCmd.PersistentFlags().StringVar(&config.NewRelic, "new-relic", "", "new relic api key")
}
