package note

import (
	"context"

	"github.com/evgeniy-lipich/microservice_go/week_3/internal/model"
)

func (s *serv) Get(ctx context.Context, id int64) (*model.Note, error) {
	note, err := s.noteRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return note, nil
}
