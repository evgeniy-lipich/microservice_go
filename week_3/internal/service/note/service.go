package note

import (
	"github.com/evgeniy-lipich/microservice_go/week_3/internal/repository"
	"github.com/evgeniy-lipich/microservice_go/week_3/internal/service"
)

type serv struct {
	noteRepository repository.NoteRepository
}

func NewService(noteRepository repository.NoteRepository) service.NoteService {
	return &serv{
		noteRepository: noteRepository,
	}
}
