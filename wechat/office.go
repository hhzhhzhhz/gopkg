package wechat

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

const (
	RandomStr                     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	WechatCryptorEncryptMsgFormat = `
<xml>
<Encrypt><![CDATA[%s]]></Encrypt>
<MsgSignature><![CDATA[%s]]></MsgSignature>
<TimeStamp>%s</TimeStamp>
<Nonce><![CDATA[%s]]></Nonce>
</xml>
`
)

type OfficePostBody struct {
	XMLName    xml.Name `xml:"xml"`
	ToUserName string   `xml:"ToUserName"`
	AppId      string   `xml:"AppId"`
	Encrypt    string   `xml:"Encrypt"`
}

// OfficeCrypto 微信公众号加解密
type OfficeCrypto struct {
	appId  string
	token  string
	aesKey []byte
}

func NewOfficeCrypto(appId, token, AESKey string) (*OfficeCrypto, error) {
	aes, err := base64.StdEncoding.DecodeString(AESKey + "=")
	if err != nil {
		return nil, fmt.Errorf("")
	}
	return &OfficeCrypto{
		appId:  appId,
		token:  token,
		aesKey: aes,
	}, nil
}

// 随机生成16位字符串
func (o *OfficeCrypto) random16Str() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := []byte(RandomStr)
	var result []byte
	for i := 0; i < 16; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func (o *OfficeCrypto) EncryptMsg(msg, timeStamp, nonce string) (string, error) {
	encrypt, sign, timeStamp, nonce, err := o.EncryptMsgContent(msg, timeStamp, nonce)
	if nil != err {
		return "", err
	}
	return fmt.Sprintf(WechatCryptorEncryptMsgFormat, encrypt, sign, timeStamp, nonce), nil
}

func (o *OfficeCrypto) DecryptMsg(msgSign, timeStamp, nonce, postData string) (string, error) {
	postBody := OfficePostBody{}
	err := xml.Unmarshal([]byte(postData), &postBody)
	if nil != err || 0 == len(postBody.Encrypt) {
		return "", fmt.Errorf("parse xml error: %s", err.Error())
	}

	return o.DecryptMsgContent(msgSign, timeStamp, nonce, postBody.Encrypt)
}

func (o *OfficeCrypto) EncryptMsgContent(msg, timeStamp, nonce string) (string, string, string, string, error) {
	encrypt, err := o.Encrypt(o.random16Str(), msg)
	if nil != err {
		return "", "", "", "", err
	}

	if 0 == len(timeStamp) {
		timeStamp = fmt.Sprint(time.Now().Unix())
	}

	sign := o.sha1(o.token, timeStamp, nonce, encrypt)
	return encrypt, sign, timeStamp, nonce, nil
}

func (o *OfficeCrypto) DecryptMsgContent(msgSign, timeStamp, nonce, encrypt string) (string, error) {
	sign := o.sha1(o.token, timeStamp, nonce, encrypt)
	if msgSign != sign {
		return "", fmt.Errorf("validate signature")
	}

	return o.Decrypt(encrypt)
}

// 对明文进行加密
func (o *OfficeCrypto) Encrypt(randomStr, text string) (string, error) {
	randomBytes := []byte(randomStr)
	textBytes := []byte(text)
	networkBytes := o.buildNetworkBytesOrder(len(textBytes))
	appIdBytes := []byte(o.appId)
	var unencrypted []byte
	unencrypted = append(unencrypted, randomBytes...)
	unencrypted = append(unencrypted, networkBytes...)
	unencrypted = append(unencrypted, textBytes...)
	unencrypted = append(unencrypted, appIdBytes...)
	encrypted, err := o.encrypt(unencrypted, o.aesKey)
	if nil != err {
		return "", fmt.Errorf("encrypt ase error %s", err.Error())
	}
	return encrypted, nil
}

// 对密文进行解密
func (o *OfficeCrypto) Decrypt(text string) (string, error) {
	original, err := o.decrypt(text, o.aesKey)
	if nil != err {
		return "", fmt.Errorf("decrypt error %s", err.Error())
	}
	networkBytes := original[16:20]
	textLen := o.recoverNetworkBytesOrder(networkBytes)
	textBytes := original[20 : 20+textLen]
	appIdBytes := original[20+textLen:]
	if o.appId != string(appIdBytes) {
		return "", fmt.Errorf("validate appid")
	}
	return string(textBytes), nil
}

func (o *OfficeCrypto) encrypt(rawData, key []byte) (string, error) {
	data, err := o.AesCBCEncrypt(rawData, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func (o *OfficeCrypto) decrypt(rawData string, key []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(rawData)
	if err != nil {
		return nil, err
	}
	dnData, err := o.cbcDecrypt(data, key)
	if err != nil {
		return nil, err
	}
	return dnData, nil
}

// AesCBCEncrypt 加密，填充秘钥key的16位
func (o *OfficeCrypto) AesCBCEncrypt(rawData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// 填充原文
	rawData = o.PKCS7Padding(rawData)
	cipherText := make([]byte, len(rawData))
	// 初始向量IV
	iv := key[:16]

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, rawData)

	return cipherText, nil
}

// aes解密
func (o *OfficeCrypto) cbcDecrypt(encryptData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	// 初始向量IV
	iv := key[:16]
	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(encryptData, encryptData)
	// 解填充
	encryptData = o.PKCS7UnPadding(encryptData)
	return encryptData, nil
}

func (o *OfficeCrypto) PKCS7Padding(ciphertext []byte) []byte {
	amountToPad := 32 - (len(ciphertext) % 32)
	if 0 == amountToPad {
		amountToPad = 32
	}
	padChr := (byte)(amountToPad & 0xFF)
	result := make([]byte, len(ciphertext), len(ciphertext)+amountToPad)
	copy(result, ciphertext)
	for i := 0; i < amountToPad; i++ {
		result = append(result, padChr)
	}
	return result
}

func (o *OfficeCrypto) PKCS7UnPadding(origData []byte) []byte {
	pad := (int)(origData[len(origData)-1])
	if pad < 1 || pad > 32 {
		pad = 0
	}
	return origData[:len(origData)-pad]
}

// 生成4个字节的网络字节序
func (o *OfficeCrypto) buildNetworkBytesOrder(number int) []byte {
	return []byte{
		(byte)(number >> 24 & 0xFF),
		(byte)(number >> 16 & 0xF),
		(byte)(number >> 8 & 0xFF),
		(byte)(number & 0xFF),
	}
}

// 还原4个字节的网络字节序
func (o *OfficeCrypto) recoverNetworkBytesOrder(orderBytes []byte) int {
	var number = 0
	for i := 0; i < 4; i++ {
		number <<= 8
		number |= (int)(orderBytes[i] & 0xff)
	}
	return number
}

func (o *OfficeCrypto) sha1(token, timestamp, nonce, encrypt string) string {
	array := []string{token, timestamp, nonce, encrypt}
	sort.Strings(array)
	str := strings.Join(array, "")

	hash := sha1.New()
	hash.Write([]byte(str))
	sum := hash.Sum(nil)
	sumHex := make([]byte, hex.EncodedLen(len(sum)))
	hex.Encode(sumHex, sum)
	return string(sumHex)
}
