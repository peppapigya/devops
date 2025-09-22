package util

import (
	"k8s-platform-go/internal/dal/redis"

	"github.com/mojocn/base64Captcha"
)

func GetCaptcha(redisClient *redis.Client) (string, string, error) {
	digit := base64Captcha.DefaultDriverDigit
	id, base64String, _, err := base64Captcha.NewCaptcha(digit, redisClient).Generate()
	return id, base64String, err
}

func VerifyCaptcha(redisClient *redis.Client, id, answer string, clear bool) bool {
	captcha := base64Captcha.NewCaptcha(base64Captcha.DefaultDriverDigit, redisClient)
	return captcha.Verify(id, answer, clear)
}
