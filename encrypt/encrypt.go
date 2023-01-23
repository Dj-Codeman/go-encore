package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	sys "encore/system"
	"hash"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func Create_key() string {
	rand.Seed(time.Now().UnixNano())
	var charather_range string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz123456789012345678901324569870"

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
	var charather_range string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz123456789012345678901324569870"
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
	// NOW WITH 64 more bytes

	var hash = hmac.New(sha512.New, []byte(hex_cipher_text))
	var hmac string = hex.EncodeToString(hash.Sum(nil))
	// appending hmac to file end
	hex_cipher_text += hmac

	return hex_cipher_text
}

func Decrypt(input string, key string) string {

	// getting the hmac
	var old_hmac string = input[len(input)-128:]

	// removing the hmac from the file
	var cipher_text_iv string = strings.TrimSuffix(input, old_hmac)

	// regenerating new hmac and verifing
	hash := hmac.New(sha512.New, []byte(cipher_text_iv))
	var new_hmac string = hex.EncodeToString([]byte(hash.Sum(nil)))

	// Eventually find a better comparison then if
	if new_hmac == old_hmac {
		// seperating iv from ciphertext
		var iv string = cipher_text_iv[len(cipher_text_iv)-16:]
		var cipher_text string = strings.TrimSuffix(cipher_text_iv, iv)

		// ciphertextdecoded, err := hex.DecodeString(ciphertext)
		plain_text, _ := hex.DecodeString(cipher_text)

		var byte_key []byte = []byte(key)
		var byte_iv []byte = []byte(iv)

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

		var data string = string(PKCS5UnPadding(plain_text))

		return data

	} else {
		sys.Fail("Error, unable to drcrypt file. It might have been encrypted")
		sys.Break("by a diffrent key, or it's been tampered with")

	}
	return string("What in the firery hell happened here")
}

func PKCS5Padding(ciphertext []byte, blockSize int, after int) []byte {
	var padding int = (blockSize - len(ciphertext)%blockSize)
	var padtext []byte = bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(src []byte) []byte {
	var length int = len(src)
	var unpadding int = int(src[length-1])
	return src[:(length - unpadding)]
}

func Test() (string, string) {

	// internal test
	sys.Pass("\n Running internal testing")
	// File sizes we're checking for
	file_array := [4]int32{9518, 9518000, 9518000, 9518000}
	// file_array := [4]int32{9518, 10000000, 100000000, 500000000}

	// Any test bigger than 500mb takes literall minutes to run
	// I want to keep the normal initializating quick
	// maybe add a larger file test option

	// 1024000000 1gb
	// 3000000000 3gb
	var charather_range string = "abcdefghijklmnopqrstuvwxyz1234567890"
	for i := 0; i < len(file_array); i++ {
		var Big_Data_Bytes []byte = make([]byte, file_array[i])

		var k int
		for k = range Big_Data_Bytes {
			Big_Data_Bytes[k] = charather_range[rand.Intn(len(charather_range))]
		}
		var Big_Data_Legnth int = len(Big_Data_Bytes) / 1024
		var Big_Data_Name string = "/tmp/encore/" + hex.EncodeToString([]byte(strconv.Itoa(Big_Data_Legnth))) + ".tmp"
		// write bytes to name
		// make folder if not exist /tmp/encore
		sys.WriteToFile(string(Big_Data_Bytes), Big_Data_Name, "write")

		fileBytes, _ := ioutil.ReadFile(Big_Data_Name)
		// sys.Handle_err(err, "break")

		var key string = Create_key()
		var cipher_text string = Encrypt(string(fileBytes), key)

		var Decrypted_Data string = Decrypt(cipher_text, key)

		if Decrypted_Data != string(fileBytes) {
			var msg string = "Error validating file :" + Big_Data_Name

			return "Failed", msg
		}
	}

	// sys.Handle_err(err, "warn")

	return string("Pass"), ""
}

func Larger_test() (string, string) {
	//  rework this large file test RUN TIME IS FUCKING

	file_array := [3]int64{9518, 1024000000, 3000000000}
	// ideally nobody running a FAT fs wouldnever run this because large file sizes
	// but just incase the largest file is 3gb not 4
	// OOM This doesnt finish find a better way to do it

	var charather_range string = "abcdefghijklmnopqrstuvwxyz1234567890"

	for i := 0; i < len(file_array); i++ {
		var Big_Data_Bytes []byte = make([]byte, file_array[i])

		var k int
		for k = range Big_Data_Bytes {
			Big_Data_Bytes[k] = charather_range[rand.Intn(len(charather_range))]
		}
		var Big_Data_Legnth int = len(Big_Data_Bytes) / 1024
		var Big_Data_Name string = "/tmp/encore/" + hex.EncodeToString([]byte(strconv.Itoa(Big_Data_Legnth))) + ".tmp"
		// write bytes to name
		// make folder if not exist /tmp/encore
		sys.WriteToFile(string(Big_Data_Bytes), Big_Data_Name, "write")

		fileBytes, _ := ioutil.ReadFile(Big_Data_Name)
		// sys.Handle_err(err, "break")

		var key string = Create_key()
		var cipher_text string = Encrypt(string(fileBytes), key)

		var Decrypted_Data string = Decrypt(cipher_text, key)

		if Decrypted_Data != string(fileBytes) {
			var msg string = "Error validating file :" + Big_Data_Name

			return "Failed", msg
		}
	}

	return string("Passed"), ""
}

// Copied pbkdf copied from the go site :: put link
func Pbkdf(password, salt []byte, iter, keyLen int, h func() hash.Hash) []byte {
	prf := hmac.New(h, password)
	hashLen := prf.Size()
	numBlocks := (keyLen + hashLen - 1) / hashLen

	var buf [4]byte
	dk := make([]byte, 0, numBlocks*hashLen)
	U := make([]byte, hashLen)
	for block := 1; block <= numBlocks; block++ {
		// N.B.: || means concatenation, ^ means XOR
		// for each block T_i = U_1 ^ U_2 ^ ... ^ U_iter
		// U_1 = PRF(password, salt || uint(i))
		prf.Reset()
		prf.Write(salt)
		buf[0] = byte(block >> 24)
		buf[1] = byte(block >> 16)
		buf[2] = byte(block >> 8)
		buf[3] = byte(block)
		prf.Write(buf[:4])
		dk = prf.Sum(dk)
		T := dk[len(dk)-hashLen:]
		copy(U, T)

		// U_n = PRF(password, U_(n-1))
		for n := 2; n <= iter; n++ {
			prf.Reset()
			prf.Write(U)
			U = U[:0]
			U = prf.Sum(U)
			for x := range U {
				T[x] ^= U[x]
			}
		}
	}
	return dk[:keyLen]
}
