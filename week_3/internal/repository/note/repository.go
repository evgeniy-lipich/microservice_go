package note

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/evgeniy-lipich/microservice_go/week_3/internal/client/db"
	"github.com/evgeniy-lipich/microservice_go/week_3/internal/repository"
	"github.com/evgeniy-lipich/microservice_go/week_3/internal/repository/note/converter"
	"github.com/evgeniy-lipich/microservice_go/week_3/internal/repository/note/model"
	modelRepo "github.com/evgeniy-lipich/microservice_go/week_3/internal/repository/note/model"
	desc "github.com/evgeniy-lipich/microservice_go/week_3/pkg/note_v1"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	tableName = "note"

	idColumn        = "id"
	titleColumn     = "title"
	contentColumn   = "content"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db *pgxpool.Pool
}

func (r *repo) Create(ctx context.Context, info *model.NoteInfo) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(titleColumn, contentColumn).
		Values(info.Title, info.Content).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "note_repository.Create",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*desc.Note, error) {
	builder := sq.Select(idColumn, titleColumn, contentColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "note_repository.Get",
		QueryRaw: query,
	}

	var note modelRepo.Note
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&note.ID, &note.Info.Title, &note.Info.Content, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return converter.ToNoteFromRepo(&note), nil
}

func NewRepository(db *pgxpool.Pool) repository.NoteRepository {
	return &repo{
		db: db,
	}
}