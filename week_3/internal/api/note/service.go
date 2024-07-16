package note

import (
	"github.com/evgeniy-lipich/microservice_go/week_3/internal/service"
	desc "github.com/evgeniy-lipich/microservice_go/week_3/pkg/note_v1"
)

type Implementation struct {
	desc.UnimplementedNoteV1Server
	noteService service.NoteService
}

func NewImplementation(noteService service.NoteService) *Implementation {
	return &Implementation{
		noteService: noteService,
	}
}
