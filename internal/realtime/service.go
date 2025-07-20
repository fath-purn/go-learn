package realtime

type Service interface {
	SaveMessage(input Message) (Message, error)
	GetMessageByRoom(roomID string) ([]Message, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) SaveMessage(input Message) (Message, error) {
	message := Message{}
	message.Content = input.Content
	message.SenderID = input.SenderID
	message.RoomID = input.RoomID

	return s.repository.Save(message)
}

func (s *service) GetMessageByRoom(roomID string) ([]Message, error) {
	return s.repository.FindByRoomID(roomID)
}
