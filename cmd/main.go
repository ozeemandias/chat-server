package main

import (
	"context"
	"flag"
	"log"
	"net"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ozeemandias/chat-server/internal/config"
	"github.com/ozeemandias/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var envPath string

func init() {
	flag.StringVar(&envPath, "env", ".env", "path to env file")
}

type server struct {
	chat_v1.UnimplementedChatV1Server
	dbpool *pgxpool.Pool
}

func (s *server) Create(ctx context.Context, req *chat_v1.CreateRequest) (*chat_v1.CreateResponse, error) {
	var chatID int64
	tx, err := s.dbpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "cannot start transaction: %v", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	err = tx.QueryRow(ctx, "INSERT INTO chats DEFAULT VALUES RETURNING id;").Scan(&chatID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert chat: %v", err)
	}

	for _, userID := range req.UserIds {
		createUserChatsQuery, args, err := sq.Insert("user_chats").
			Columns("user_id", "chat_id").
			Values(userID, chatID).
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to build query: %v", err)
		}

		_, err = tx.Exec(ctx, createUserChatsQuery, args...)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to insert user_chats: %v", err)
		}
	}

	return &chat_v1.CreateResponse{
		Id: chatID,
	}, nil
}

func (s *server) Delete(ctx context.Context, req *chat_v1.DeleteRequest) (*emptypb.Empty, error) {
	query, args, err := sq.Delete("chats").
		Where(sq.Eq{"id": req.GetId()}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build query: %v", err)
	}

	ct, err := s.dbpool.Exec(ctx, query, args...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete chat: %v", err)
	}

	log.Printf("deleted chats count: %d", ct.RowsAffected())

	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *chat_v1.SendMessageRequest) (*emptypb.Empty, error) {
	query, args, err := sq.Insert("messages").
		Columns("chat_id", "user_id", "message").
		Values(req.ChatId, req.From, req.Text).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build query: %v", err)
	}

	_, err = s.dbpool.Exec(ctx, query, args...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert message: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func main() {
	ctx := context.Background()

	// Считываем переменные окружения
	err := config.Load(envPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	ln, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	dbpool, err := pgxpool.New(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("unable to create connection pool: %v\n", err)
	}
	defer dbpool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	chat_v1.RegisterChatV1Server(s, &server{dbpool: dbpool})

	log.Printf("server listening at %v", ln.Addr())

	if err = s.Serve(ln); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
