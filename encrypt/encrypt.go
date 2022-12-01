package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	sys "encore/system"
	"math/rand"
	"strings"
	"time"
)

func Create_key() string {
	rand.Seed(time.Now().UnixNano())
	var charather_range string = "abcdefghijklmnopqrstuvwxyz123456789012345678901324569870"

	var key_bytes []byte = make([]byte, 32)

	for i := range key_bytes {
		key_bytes[i] = charather_range[rand.Int63()%int64(len(charather_range))]
	}

	var key string = string(key_bytes)

	return key

}

func Create_iv() []byte {
	// generating initial vector
	rand.Seed(time.Now().UnixNano())
	var charather_range string = "abcdefghijklmnopqrstuvwxyz1234567890"
	// for len in iv number pick randon byte charathers
	var iv_bytes []byte = make([]byte, 16)
	for i := range iv_bytes {
		iv_bytes[i] = charather_range[rand.Intn(len(charather_range))]
	}
	var iv_legnth int = len(iv_bytes)

	// just making sure im not dumb
	if iv_legnth <= 15 && iv_legnth >= 17 {
		sys.Break("IV is the wrong size ?")
	}

	return iv_bytes
}

func Encrypt(input string, key string) string {

	var iv_bytes []byte = Create_iv()
	var key_bytes []byte = []byte(key)
	// the file will have to be read into a variable ?
	// file size / ram implications ?
	var input_bytes []byte = PKCS5Padding([]byte(input), aes.BlockSize, len(input))
	block, err := aes.NewCipher(key_bytes)
	if err != nil {
		sys.Handle_err(err, "break")
	}
	// Turning the input into bytes
	var cipher_text []byte = make([]byte, len(input_bytes))
	// getting kinda oopish
	// make better comments
	mode := cipher.NewCBCEncrypter(block, iv_bytes)
	mode.CryptBlocks(cipher_text, input_bytes)

	// converting the bytes to hex sctring
	var hex_cipher_text string = hex.EncodeToString(cipher_text)

	// adding IV to the end of the file
	// Should be equivalent to hex.EncodeToString(iv_bytes)
	var iv string = string(iv_bytes)

	hex_cipher_text += iv

	// generate Hmac
	// 64 charathers

	var hash = hmac.New(sha256.New, []byte(hex_cipher_text))
	var hmac string = hex.EncodeToString(hash.Sum(nil))
	// appending hmac to file end
	hex_cipher_text += hmac

	return hex_cipher_text
}

func Decrypt(input string, key string) string {

	// getting the hmac
	var old_hmac string = input[len(input)-64:]

	// removing the hmac from the file
	var cipher_text_iv string = strings.TrimSuffix(input, old_hmac)

	// regenerating new hmac and verifing
	hash := hmac.New(sha256.New, []byte(cipher_text_iv))
	var new_hmac string = hex.EncodeToString([]byte(hash.Sum(nil)))

	// Eventually find a better comparison then if
	if new_hmac == old_hmac {
		// seperating iv from ciphertext
		var iv string = cipher_text_iv[len(cipher_text_iv)-16:]
		var cipher_text string = strings.TrimSuffix(cipher_text_iv, iv)

		// ciphertextdecoded, err := hex.DecodeString(ciphertext)
		plain_text, _ := hex.DecodeString(cipher_text)

		byte_key := []byte(key)
		byte_iv := []byte(iv)

		if len(byte_iv) <= 15 && len(iv) >= 17 {
			sys.Break("Invalid IV size")
		}

		block, err := aes.NewCipher(byte_key)
		if err != nil {
			sys.Handle_err(err, "break")
		}

		// CBC mode always works in whole blocks.
		if len(plain_text)%aes.BlockSize != 0 {
			sys.Break("ciphertext is not a multiple of the block size")
		}

		mode := cipher.NewCBCDecrypter(block, byte_iv)
		mode.CryptBlocks(plain_text, plain_text)

		data := string(PKCS5UnPadding(plain_text))

		return data

	} else {
		sys.Red("Error, unable to drcrypt file. It might have been encrypted")
		sys.Break("by a diffrent key, or it's been tampered with")
	}
	return string("What in the firery hell happened here")
}

func PKCS5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

// -aes-256-cbc
