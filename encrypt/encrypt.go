package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	sys "encore/system"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func Create_key() string {
	rand.Seed(time.Now().UnixNano())
	var charather_range string = "abcdefghijklmnopqrstuvwxyz123456789012345678901324569870"
	// letters := "123456789012345678901324569870"

	var key_bytes []byte = make([]byte, 32)

	for i := range key_bytes {
		key_bytes[i] = charather_range[rand.Int63()%int64(len(charather_range))]
	}

	var key string = string(key_bytes)
	// key += "\n"

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
		panic(err)
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
	iv := string(iv_bytes)

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
	old_hmac := input[len(input)-64:]

	// removing the hmac from the file
	etiv := strings.TrimSuffix(input, old_hmac)

	// regenerating new hmac and verifing
	hash := hmac.New(sha256.New, []byte(etiv))
	new_hmac := hex.EncodeToString([]byte(hash.Sum(nil)))

	// Eventually find a better comparison then if
	if new_hmac == old_hmac {
		// seperating iv from ciphertext
		iv := etiv[len(etiv)-16:]
		ciphertext := strings.TrimSuffix(etiv, iv)

		// ciphertextdecoded, err := hex.DecodeString(ciphertext)
		ciphertextdecoded, err := hex.DecodeString(ciphertext)

		byte_key := []byte(key)

		byte_iv := []byte(iv)
		fmt.Println(byte_iv)

		if len(byte_iv) <= 15 && len(iv) >= 17 {
			sys.Break("Invalid IV size")
		}

		block, err := aes.NewCipher(byte_key)
		if err != nil {
			sys.Warning("a thing")
		}

		// CBC mode always works in whole blocks.
		if len(ciphertextdecoded)%aes.BlockSize != 0 {
			sys.Break("ciphertext is not a multiple of the block size")
		}

		fmt.Println(len(ciphertextdecoded))

		mode := cipher.NewCBCDecrypter(block, byte_iv)
		mode.CryptBlocks(ciphertextdecoded, ciphertextdecoded)

		data := string(PKCS5UnPadding(ciphertextdecoded))

		return data

	} else {
		sys.Red("Error, unable to drcrypt file. It might have been encrypted")
		sys.Break("by a diffrent key, or it's been tampered with")
	}

	return "I fucked up somewhere"
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
