package tokenizer

import (
	"TaskTracker/internal/models"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"log/slog"
	"time"
)

type Params struct {
	fx.In

	Config Config
	Logger *slog.Logger
}

type Tokenizer struct {
	cfg Config
	log *slog.Logger
}

func New(p Params) *Tokenizer {
	if p.Logger == nil || p.Config.KeyJWT == nil {
		panic("failed to initialize Tokenizer: missing dependencies")
	}

	p.Logger.Info("Tokenizer created successfully")
	return &Tokenizer{
		cfg: p.Config,
		log: p.Logger,
	}
}

func (t *Tokenizer) GenerateJWT(payload *models.TokenPayload) (*models.TokenResponse, error) {
	t.log.Debug("where")
	t.log.Debug(time.Now().String())
	t.log.Debug(t.cfg.AccessExpirationTime.String())
	t.log.Debug(time.Now().Add(t.cfg.AccessExpirationTime).String())
	expTime := time.Now().Add(t.cfg.AccessExpirationTime)
	t.log.Debug("GenerateJWT", "exp", expTime)

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": payload.UserID,
		"exp": expTime.Unix(),
	})
	tokenStr, err := token.SignedString(t.cfg.KeyJWT)
	if err != nil {
		return nil, err
	}

	tokenResponse := &models.TokenResponse{
		Token: tokenStr,
		Exp:   expTime,
	}

	return tokenResponse, nil
}

func (t *Tokenizer) ValidateJWT(tokenString string) (*models.TokenPayload, error) {
	t.log.Debug(tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(`unexpected signing method`)
		}

		return t.cfg.KeyJWT, nil
	})
	if err != nil {
		t.log.Error("parsing token", "error", err)
		return nil, errors.New("invalid token")
	}

	payload, err := parseClaims(token)
	if err != nil {
		t.log.Error("parsing token claims", "error", err)
		return nil, errors.New("invalid token")
	}

	if payload.Exp.Before(time.Now()) {
		t.log.Error("token expired")
		return nil, errors.New("token expired")
	}

	return payload, nil
}

//	func (t *Tokenizer) GeneratePairToken(payload *models.TokenPayload) (*models.PairToken, error) {
//		pair := &models.PairToken{}
//		payload.Exp = time.Now().Add(t.cfg.AccessExpirationTime)
//		pair.ExpAccessToken = payload.Exp
//		accessToken, err := t.GenerateJWT(payload)
//		if err != nil {
//			t.log.Error("generating access token", "error", err)
//			return nil, err
//		}
//		pair.AccessToken = accessToken
//
//		payload.Exp = time.Now().Add(t.cfg.RefreshExpirationTime)
//		pair.ExpRefreshToken = payload.Exp
//		refreshToken, err := t.GenerateJWT(payload)
//		if err != nil {
//			t.log.Error("generating refresh token", "error", err)
//			return nil, err
//		}
//		pair.RefreshToken = refreshToken
//
//		return pair, nil
//	}
func parseClaims(token *jwt.Token) (*models.TokenPayload, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}

	userID, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		return nil, errors.New("invalid userID in token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("invalid exp in token claims")
	}
	expTime := time.Unix(int64(exp), 0)

	return &models.TokenPayload{
		UserID: userID,
		Exp:    expTime,
	}, nil
}
