package middleware

import (
	"context"
	"errors"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/pbuser/server/jwt"
)

var UserKey string

type TokenInfo struct {
	UserId uint `json:"user_id"`
}

func AuthInterceptor(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	//fmt.Printf("token is %v\n", token)

	jwtManager := jwt.NewJwtManager()

	claims, err := jwtManager.ParseToken(token)
	if err != nil {
		return nil, err
	}

	//fmt.Println(claims.UserId)
	return context.WithValue(ctx, UserKey, &TokenInfo{UserId: claims.UserId}), nil
}

func GetUserInfo(ctx context.Context) (*TokenInfo, error) {

	user, ok := ctx.Value(UserKey).(*TokenInfo)
	if !ok {
		return nil, errors.New("get user info fail")
	}

	return user, nil
}
