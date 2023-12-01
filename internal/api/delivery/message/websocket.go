package message

// "github.com/gorilla/websocket"

// conn, err := websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)
// if err != nil {
// 	httpResponse.Fail(err, h.logger).ToJSON(c)
// 	return
// }
// if err = h.messageSrv.SendMessage(ctx, conn, "123"); err != nil {
// 	httpResponse.Fail(err, h.logger).ToJSON(c)
// 	return
// }
// conn, err := websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)
// if err != nil {
// 	panic(err)
// }
// client := &Client{
// 	Connection: conn,
// 	Send:       make(chan []byte),
// }
// clients[client] = true
// fmt.Println(clients)
// go handleWebSocketConnection(client, clients, broadcast)
// conn.Close()

// func handleWebSocketConnection(client *Client, clients map[*Client]bool, broadcast chan []byte) {
// 	defer func() {
// 		client.Connection.Close()
// 		delete(clients, client)
// 	}()
// 	for {
// 		_, message, err := client.Connection.ReadMessage()
// 		if err != nil {
// 			break
// 		}
// 		broadcast <- message
// 	}
// }

// func handleBroadcast(clients map[*Client]bool, broadcast chan []byte) {
// 	for {
// 		message := <-broadcast
// 		for client := range clients {
// 			select {
// 			case client.Send <- message:
// 			default:
// 				close(client.Send)
// 				delete(clients, client)
// 			}
// 		}
// 	}
// }
