package service

import (
	"https://github.com/riri95500/go-chat/broadcast"
)

type Manager interface {
	OpenListener(roomid string) chan interface{}
	CloseListener(roomid string, channel chan interface{})
	Submit(userid, roomid, text string)
	DeleteBroadcast(roomid string)
}

type Message struct {
	UserId string
	RoomId string
	Text   string
}

type Listener struct {
	RoomId string
	Chan   chan interface{}
}

type manager struct {
	roomChannels map[string]broadcast.Broadcaster
	open         chan *Listener
	close        chan *Listener
	delete       chan string
	messages     chan *Message
}

// Cette fonction déclenchera register
func (m *manager) OpenListener(roomid string) chan interface{} {
	listener := make(chan interface{})
	m.open <- &Listener{
		RoomId: roomid,
		Chan:   listener,
	}
	return listener
}

// Cette fonction déclenchera deregister
func (m *manager) CloseListener(roomid string, channel chan interface{}) {
	m.close <- &Listener{
		RoomId: roomid,
		Chan:   channel,
	}
}

// Cette fonction déclenchera deleteBroadcast
func (m *manager) DeleteBroadcast(roomid string) {
	m.delete <- roomid
}

// Cette fonction déclenchera broadcast.Submit
func (m *manager) Submit(userid, roomid, text string) {
	msg := &Message{
		UserId: userid,
		RoomId: roomid,
		Text:   text,
	}
	m.messages <- msg
}

func (m *manager) register(listener *Listener) {
	m.room(listener.RoomId).Register(listener.Chan)
}

func (m *manager) deregister(listener *Listener) {
	m.room(listener.RoomId).Unregister(listener.Chan)
	close(listener.Chan)
}

func (m *manager) deleteBroadcast(roomid string) {
	b, ok := m.roomChannels[roomid]
	if ok {
		b.Close()
		delete(m.roomChannels, roomid)
	}
}

/*
Get the room with the id roomid, or creates and registers it
*/
func (m *manager) room(roomid string) broadcast.Broadcaster {
	b, ok := m.roomChannels[roomid]
	if !ok {
		b = broadcast.NewBroadcaster(10)
		m.roomChannels[roomid] = b
	}
	return b
}

func (m *manager) run() {
	for {
		select {
		//Cette fonction sera déclenché à l'appel de OpenListener
		case listener := <-m.open:
			m.register(listener)
		//Cette fonction sera déclenché à l'appel de CloseListener
		case listener := <-m.close:
			m.deregister(listener)
		//Cette fonction sera déclenché à l'appel de DeleteBroadcast
		case roomid := <-m.delete:
			m.deleteBroadcast(roomid)
		//Cette fonction sera déclenché à l'appel de Submit
		case message := <-m.messages:
			m.room(message.RoomId).Submit(*message)
		}
	}
}

var managerSingleton *manager

func GetRoomManager() Manager {
	if managerSingleton == nil {
		managerSingleton = &manager{
			roomChannels: make(map[string]broadcast.Broadcaster),
			open:         make(chan *Listener, 100),
			close:        make(chan *Listener, 100),
			delete:       make(chan string, 100),
			messages:     make(chan *Message, 100),
		}

		go managerSingleton.run()
	}

	return managerSingleton
}
