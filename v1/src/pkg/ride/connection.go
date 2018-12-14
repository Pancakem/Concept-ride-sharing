package ride

// Hub stores active drivers by their connections
type Hub struct {
	// connected clients (mostly drivers)
	clients map[string]*Client
	// add a new driver to the list
	register chan *aggreg
	// remove an inactive driver
	unregister chan *aggreg
}

type aggreg struct {
	id string
	*Client
}

// NewHub returns an instance of Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *aggreg),
		unregister: make(chan *aggreg),
	}
}

// Check finds a writer in the list of riders made online
func (h *Hub) Check(driverid string) *Client {
	if val, ok := h.clients[driverid]; ok {
		return val
	}
	return nil
}

// run adds connections to hub and also removes them
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.id] = client.Client
		case client := <-h.unregister:
			if _, ok := h.clients[client.id]; ok {
				delete(h.clients, client.id)
				close(client.send)
			}
		}
	}
}
