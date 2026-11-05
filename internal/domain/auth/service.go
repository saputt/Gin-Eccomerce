package auth

import (
	"context"
	"errors"
	"fmt"
	"orders/internal/domain/user"
	iuser "orders/internal/domain/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserService    UserServiceContract
	JwtSecret      string
	JwtExpiresHour int
}

func NewAuthService(JwtSecret string, JwtExpiresHour int, userService UserServiceContract) *AuthService {
	return &AuthService{
		JwtSecret:      JwtSecret,
		JwtExpiresHour: JwtExpiresHour,
		UserService:    userService,
	}
}

type JWTClaims struct {
	Id    string        `json:"id"`
	Email string        `json:"email"`
	Role  user.UserRole `json:"role"`

	jwt.RegisteredClaims
}

func (s *AuthService) generateToken(user *user.User) (string, error) {
	claims := JWTClaims{
		Id:    user.Id,
		Email: user.Email,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Id,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.JwtExpiresHour) * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.JwtSecret))
}

func (s *AuthService) ValidateToken(token string) (*JWTClaims, error) {
	claims := &JWTClaims{}

	parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.JwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (s *AuthService) Login(ctx context.Context, dto *LoginRequest) (LoginResponse, error) {
	fmt.Println("masuk")

	user, err := s.UserService.FindByEmail(ctx, dto.Email)
	if err != nil {
		if errors.Is(err, iuser.ErrUserNotFound) {
			return LoginResponse{}, ErrInvalidCredential
		}
		return LoginResponse{}, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password)); err != nil {
		return LoginResponse{}, ErrInvalidCredential
	}

	token, err := s.generateToken(user)

	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		Token: token,
		User: UserResponse{
			Id:    user.Id,
			Name:  user.Name,
			Email: user.Email,
			Role:  string(user.Role),
		},
	}, nil
}

func (s *AuthService) Register(ctx context.Context, dto *RegisterRequest) (RegisterResponse, error) {
	fmt.Println("masuk")

	hashed, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return RegisterResponse{}, err
	}

	user := user.User{
		Id:       uuid.NewString(),
		Name:     dto.Name,
		Email:    dto.Email,
		Password: string(hashed),
		Role:     user.UserRole(dto.Role),
	}

	createdUser, err := s.UserService.CreateUser(ctx, &user)
	if err != nil {
		if errors.Is(err, iuser.ErrEmailAlreadyExist) {
			return RegisterResponse{}, ErrInvalidCredential
		}
		return RegisterResponse{}, err
	}

	return RegisterResponse{
		User: UserResponse{
			Id:    createdUser.Id,
			Name:  createdUser.Name,
			Email: createdUser.Email,
			Role:  string(createdUser.Role),
		},
	}, nil
}
