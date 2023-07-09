package rooms

type Room struct {
	ID   string
	Name string
}

type RoomRepository interface {
	CreateRoom(room Room) error
	DeleteRoom(roomID string) error
	GetRoomByID(roomID string) (Room, error)
}

type roomRepository struct {
	rooms map[string]Room
}

func (repo *roomRepository) CreateRoom(room Room) error {
	// Votre implémentation pour créer un salon et l'ajouter à la liste des salons
	return nil
}

func (repo *roomRepository) DeleteRoom(roomID string) error {
	// Votre implémentation pour supprimer un salon de la liste des salons
	return nil
}

func (repo *roomRepository) GetRoomByID(roomID string) (Room, error) {
	// Votre implémentation pour obtenir un salon par son ID
	return Room{}, nil
}

func NewRoomRepository() RoomRepository {
	return &roomRepository{
		rooms: make(map[string]Room),
	}
}
