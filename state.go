package main

type State struct {
	clients     map[*Client]bool
	subscribe   chan *Client
	unsubscribe chan *Client
	broadcast   chan *Message
}

func createState() *State {
	return &State{
		make(map[*Client]bool),
		make(chan *Client),
		make(chan *Client),
		make(chan *Message),
	}
}

func (state *State) start() {
	for {
		select {
		case client := <-state.subscribe:
			state.clients[client] = true
		case client := <-state.unsubscribe:
			if _, ok := state.clients[client]; ok {
				delete(state.clients, client)
				close(client.send)
			}
		case message := <-state.broadcast:
			for client := range state.clients {
				var isRelatedId = client.RelatedObjectId == message.RelatedObjectId
				var isRelatedType = client.RelatedObjectType == message.RelatedObjectType
				var isNotSelf = client.UserId != message.UserId

				if isRelatedType && isRelatedId && isNotSelf {
					select {
					case client.send <- message.toSend():
					default:
						close(client.send)
						delete(state.clients, client)
					}
				}
			}
		}
	}
}
