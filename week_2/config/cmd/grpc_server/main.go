package main

import (
	"context"
	"database/sql"
	"flag"
	"github.com/evgeniy-lipich/microservice_go/week_2/config/internal/config/env"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	desc "github.com/evgeniy-lipich/microservice_go/week_2/config/pkg/note_v1"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/evgeniy-lipich/microservice_go/week_2/config/internal/config"
)

var configPath string

func init() {
	// при запуске будет передаваться наименования конфига в параметрах
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedNoteV1Server
	pool *pgxpool.Pool
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	// запрос на создание записи в БД
	builderInsert := sq.Insert("note").
		PlaceholderFormat(sq.Dollar).
		Columns("title", "body").
		Values(gofakeit.City(), gofakeit.Address().Street).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var noteID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&noteID)
	if err != nil {
		log.Fatalf("failed to insert note: %v", err)
	}

	log.Printf("inserted note with id: %d", noteID)

	return &desc.CreateResponse{
		Id: noteID,
	}, nil
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	// запрос на получение записи из БД
	builderSelectOne := sq.Select("id", "title", "body", "created_at", "updated_at").
		From("note").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()}).
		Limit(1)

	query, args, err := builderSelectOne.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var id int64
	var title, body string
	var createdAt time.Time
	var updatedAt sql.NullTime

	err = s.pool.QueryRow(ctx, query, args...).Scan(&id, &title, &body, &createdAt, &updatedAt)
	if err != nil {
		log.Fatalf("failed to fetch note: %v", err)
	}

	log.Printf("id: %d, title: %s, body: %s, created_at: %v, updated_at: %v\n", id, title, body, createdAt, updatedAt)

	var updatedAtTime *timestamppb.Timestamp
	if updatedAt.Valid {
		updatedAtTime = timestamppb.New(updatedAt.Time)
	}

	return &desc.GetResponse{
		Note: &desc.Note{
			Id:        id,
			Info:      &desc.NoteInfo{Title: title, Content: body},
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: updatedAtTime,
		},
	}, nil
}

func main() {
	// считываем параметр с консоли
	flag.Parse()
	ctx := context.Background()

	// считываем переменные окружения
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterNoteV1Server(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
