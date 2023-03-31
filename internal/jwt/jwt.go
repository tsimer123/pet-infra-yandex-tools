package jwt

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"encoding/base64"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"github.com/tsimer123/pet-infra-yandex-tools/internal/env"
)

type JWT struct {
	keyID            string
	serviceAccountID string
	key              string
}

func NewJWT(o *env.Options) *JWT {
	return &JWT{
		keyID:            o.YaJWTkeyID,
		serviceAccountID: o.YaJWTserviceAccountID,
		key:              o.YaJWTkey,
	}
}

func (t *JWT) signedToken() string {
	claims := jwt.RegisteredClaims{
		Issuer:    t.serviceAccountID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Audience:  []string{"https://iam.api.cloud.yandex.net/iam/v1/tokens"},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)
	token.Header["kid"] = t.keyID

	privateKey := t.loadPrivateKey()
	signed, err := token.SignedString(privateKey)
	if err != nil {
		logrus.Errorf("Error signing token: %s", err)
	}

	logrus.Infof("Signed token: %s", signed)

	return signed
}

func (t *JWT) loadPrivateKey() *rsa.PrivateKey {
	key, err := base64.RawStdEncoding.DecodeString(t.key)
	if err != nil {
		logrus.Errorf("Error decoding private key: %s", err)
	}

	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil {
		logrus.Errorf("Error parsing private key: %s", err)
	}

	logrus.Infof("Private key: %s", "<sensitive value>")

	return rsaPrivateKey
}

func (t *JWT) GetIAMToken() string {
	jot := t.signedToken()
	fmt.Println(jot)
	resp, err := http.Post(
		"https://iam.api.cloud.yandex.net/iam/v1/tokens",
		"application/json",
		strings.NewReader(fmt.Sprintf(`{"jwt":"%s"}`, jot)),
	)
	if err != nil {
		logrus.Errorf("Error getting IAM token: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logrus.Errorf("Error getting IAM token: %s", body)
	}
	var data struct {
		IAMToken string `json:"iamToken"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		logrus.Errorf("Error decoding response: %s", err)
	}

	logrus.Infof("Got IAM token: %s", "<sensitive value>")

	return data.IAMToken
}
