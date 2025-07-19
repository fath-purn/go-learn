package book

type Service interface {
	Create(book BookRequest) (Book, error)
	FindAll() ([]Book, error)
	FIndByID(ID int) (Book, error)
	Update(ID int, book BookRequest) (Book, error)
	Delete(ID int) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) Create(bookRequest BookRequest) (Book, error) {
	book := Book{
		Title:       bookRequest.Title,
		Price:       bookRequest.Price,
		Synopsis:    bookRequest.Synopsis,
		Description: bookRequest.Description,
		Rating:      bookRequest.Rating,
	}

	createdBook, err := s.repository.Create(book)
	if err != nil {
		return Book{}, err
	}
	return createdBook, nil
}

func (s *service) FindAll() ([]Book, error) {
	books, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (s *service) FIndByID(ID int) (Book, error) {
	book, err := s.repository.FIndByID(ID)
	if err != nil {
		return Book{}, err
	}
	return book, nil
}

func (s *service) Update(ID int, bookRequest BookRequest) (Book, error) {
	book, err := s.repository.FIndByID(ID)
	if err != nil {
		return Book{}, err
	}

	book.Title = bookRequest.Title
	book.Price = bookRequest.Price
	book.Synopsis = bookRequest.Synopsis
	book.Description = bookRequest.Description
	book.Rating = bookRequest.Rating

	updatedBook, err := s.repository.Update(book)
	if err != nil {
		return Book{}, err
	}
	return updatedBook, nil
}

func (s *service) Delete(ID int) error {

	if err := s.repository.Delete(ID); err != nil {
		return err
	}
	return nil
}
