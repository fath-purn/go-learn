package realtime

type Service interface {
	SaveMessage(input Message) (Message, error)
	GetMessage() ([]Message, error)
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

	return s.repository.Save(message)
}

func (s *service) GetMessage() ([]Message, error) {
	return s.repository.FindAll()
}
