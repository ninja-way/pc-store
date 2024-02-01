package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/ninja-way/pc-store/internal/models"
	"math/rand"
	"strconv"
	"time"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UsersRepository interface {
	CreateUser(context.Context, models.User) error
	GetUser(context.Context, string, string) (models.User, error)
}

type SessionsRepository interface {
	CreateToken(ctx context.Context, token models.RefreshSession) error
	GetToken(ctx context.Context, token string) (models.RefreshSession, error)
}

type Users struct {
	repo         UsersRepository
	sessionsRepo SessionsRepository

	hasher     PasswordHasher
	hmacSecret []byte
	tokenTTL   time.Duration
}

func NewUsers(repo UsersRepository, sessionsRepo SessionsRepository, hasher PasswordHasher, secret []byte, ttl time.Duration) *Users {
	return &Users{
		repo:         repo,
		sessionsRepo: sessionsRepo,
		hasher:       hasher,
		hmacSecret:   secret,
		tokenTTL:     ttl,
	}
}

func (u *Users) SignUp(ctx context.Context, inp models.SignUp) error {
	password, err := u.hasher.Hash(inp.Password)
	if err != nil {
		return err
	}

	user := models.User{
		Name:         inp.Name,
		Email:        inp.Email,
		Password:     password,
		RegisteredAt: time.Now(),
	}

	return u.repo.CreateUser(ctx, user)
}

func (u *Users) SignIn(ctx context.Context, inp models.SignIn) (string, string, error) {
	password, err := u.hasher.Hash(inp.Password)
	if err != nil {
		return "", "", err
	}

	user, err := u.repo.GetUser(ctx, inp.Email, password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", models.ErrUserNotFound
		}
		return "", "", err
	}

	return u.generateTokens(ctx, user.ID)
}

func (u *Users) ParseToken(ctx context.Context, token string) (int64, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return u.hmacSecret, nil
	})
	if err != nil {
		return 0, err
	}

	if !t.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return 0, errors.New("invalid subject")
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		return 0, errors.New("invalid subject")
	}

	return int64(id), nil
}

func (u *Users) generateTokens(ctx context.Context, userID int) (string, string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.Itoa(userID),
		IssuedAt:  &jwt.NumericDate{Time: time.Now()},
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(u.tokenTTL)},
	})

	accessToken, err := t.SignedString(u.hmacSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}

	if err := u.sessionsRepo.CreateToken(ctx, models.RefreshSession{
		UserID:    int64(userID),
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	}); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func newRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func (u *Users) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	session, err := u.sessionsRepo.GetToken(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	if session.ExpiresAt.Unix() < time.Now().Unix() {
		return "", "", models.ErrRefreshTokenExpired
	}

	return u.generateTokens(ctx, int(session.UserID))
}
