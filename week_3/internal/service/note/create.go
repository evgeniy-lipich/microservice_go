package note

import (
	"context"

	"github.com/evgeniy-lipich/microservice_go/week_3/internal/model"
)

func (s *serv) Create(ctx context.Context, info *model.NoteInfo) (int64, error) {
	return 0, nil
}
