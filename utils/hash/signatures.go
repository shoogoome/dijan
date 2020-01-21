package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"dijan/utils"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

// 生成cookie
func CookieSignature(mes string) string {
	return encryptionString(mes)
}

// 加密字符串
func encryptionString(payloadString string) string {
	ghmac := hmac.New(sha256.New, []byte(utils.GlobalSystemConfig.Server.HashSalt))
	ghmac.Write([]byte(payloadString + utils.GlobalSystemConfig.Server.HashSalt))
	return hex.EncodeToString(ghmac.Sum([]byte(nil)))
}

// 生成token
func GenerateToken(payload interface{}, encryption ...bool) string {
	if payload == nil {
		return ""
	}
	payloadStr, err := json.Marshal(payload)
	if err != nil {
		return ""
	}
	payloadString := base64.StdEncoding.EncodeToString([]byte(payloadStr))
	if len(encryption) > 0 && encryption[0] {
		a := fmt.Sprintf("%s.%s", payloadString, encryptionString(payloadString))
		a = strings.Replace(a, "+", "-", -1)
		a = strings.Replace(a, "/", "*", -1)
		return a
	}
	return payloadString
}

func DecodeToken(token string, object interface{}) {
	token = strings.Replace(token, "-", "+", -1)
	token = strings.Replace(token, "*", "/", -1)
	if len(token) == 0 {
		panic("")
	}
	// 格式验证
	userToken := strings.Split(token, ".")
	if len(userToken) == 2 {
		// 密钥验证
		if CookieSignature(userToken[0]) != userToken[1] {
			panic("")
		}
	}
	payloadStr, err := base64.StdEncoding.DecodeString(userToken[0])
	if err != nil {
		panic("")
	}
	// 反序列化
	err = json.Unmarshal(payloadStr, &object)
	if err != nil {
		panic("")
	}
}