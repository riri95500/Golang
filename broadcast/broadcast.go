package broadcast

type Broadcaster interface {
	// Register a new channel to receive broadcasts
	Register(chan<- interface{})
	// Unregister a channel so that it no longer receives broadcasts.
	Unregister(chan<- interface{})
	// Shut this broadcaster down.
	Close() error
	// Submit a new object to all subscribers
	Submit(interface{}) bool
}

type broadcaster struct {
	input chan interface{}
	reg   chan chan<- interface{}
	unreg chan chan<- interface{}

	outputs map[chan<- interface{}]bool
}

func (b *broadcaster) broadcast(m interface{}) {
	//On diffuse le msg a tout les listeners(tout les viewers du chat)
	for ch := range b.outputs {
		ch <- m
	}
}

func (b *broadcaster) run() {
	for {
		//Le select attends qu'un de ses case s'éxécute
		select {
		case m := <-b.input:
			b.broadcast(m)
		//ok a true si channel ouverte false sinon
		//c'est Close() qui fermera la channel
		//La channel sort ici
		case ch, ok := <-b.reg:
			if ok {
				b.outputs[ch] = true
			} else {
				return
			}
		case ch := <-b.unreg:
			delete(b.outputs, ch)
		}
	}
}

// Chanel qui recoit uniquement des msg et n'en envoi pas
func (b *broadcaster) Register(newch chan<- interface{}) {
	//On enregistre newch dans la chanel reg
	//Une channel entre et doit obligatoirement sortir
	//Il faut trouver ou la channel sort
	b.reg <- newch
}

func (b *broadcaster) Unregister(ch chan<- interface{}) {
	b.unreg <- ch
}

func (b *broadcaster) Close() error {
	close(b.reg)
	close(b.unreg)
	return nil
}

func (b *broadcaster) Submit(m interface{}) bool {
	if b == nil {
		return false
	}
	select {
	case b.input <- m:
		return true
	default:
		return false
	}
}

func NewBroadcaster(buflen int) Broadcaster {
	//Initialisation des channels avec make
	b := &broadcaster{
		input:   make(chan interface{}, buflen),
		reg:     make(chan chan<- interface{}),
		unreg:   make(chan chan<- interface{}),
		outputs: make(map[chan<- interface{}]bool),
	}

	go b.run()

	return b
}
