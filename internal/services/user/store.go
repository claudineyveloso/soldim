package user

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/claudineyveloso/soldim.git/internal/db"
	"github.com/claudineyveloso/soldim.git/internal/services/auth"
	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/google/uuid"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(user types.UserPayload) error {
	queries := db.New(s.db)
	ctx := context.Background()

	user.ID = uuid.New()
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		slog.Error("Erro ao gerar hash da senha", slog.String("error", err.Error()))
		return err
	}

	createUserParams := db.CreateUserParams{
		ID:        user.ID,
		Email:     user.Email,
		Password:  string(hashedPassword),
		IsActive:  user.IsActive,
		UserType:  user.UserType,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if err := queries.CreateUser(ctx, createUserParams); err != nil {
		slog.Error("Erro ao criar um usuário", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (s *Store) GetUsers() ([]*types.User, error) {
	queries := db.New(s.db)
	ctx := context.Background()

	dbUsers, err := queries.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	var users []*types.User
	for _, dbUser := range dbUsers {
		user := convertDBUserToUser(dbUser)
		users = append(users, user)
	}
	return users, nil
}

func (s *Store) GetUserByID(userID uuid.UUID) (*types.User, error) {
	queries := db.New(s.db)
	ctx := context.Background()
	dbUser, err := queries.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	user := convertDBUserToUser(dbUser)

	return user, nil
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	queries := db.New(s.db)
	ctx := context.Background()

	dbUser, err := queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	user := &types.User{
		ID:       dbUser.ID,
		Email:    dbUser.Email,
		Password: dbUser.Password,
		IsActive: dbUser.IsActive,
		UserType: dbUser.UserType,
	}

	return user, nil
}

func (s *Store) LoginUser(user types.CreateLoginPayload) (*types.User, error) {
	queries := db.New(s.db)
	ctx := context.Background()

	dbUser, err := queries.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	if !auth.ComparePasswords(dbUser.Password, []byte(user.Password)) {
		return nil, fmt.Errorf("senha inválida")
	}

	loginUserParams := db.LoginUserParams{
		Email:    user.Email,
		Password: dbUser.Password,
	}

	loggedUser, err := queries.LoginUser(ctx, loginUserParams)
	if err != nil {
		fmt.Println("Erro ao fazer login do usuário:", err)
		return nil, err
	}
	convertedUser := convertDBUserToUser(loggedUser)

	return convertedUser, nil
}

func convertDBUserToUser(dbUser db.User) *types.User {
	user := &types.User{
		ID:        dbUser.ID,
		Email:     dbUser.Email,
		Password:  dbUser.Password,
		IsActive:  dbUser.IsActive,
		UserType:  dbUser.UserType,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}
	return user
}
