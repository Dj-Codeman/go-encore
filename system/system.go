package system

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	ColorBlack  = "\u001b[30m"
	ColorRed    = "\u001b[31m"
	ColorGreen  = "\u001b[32m"
	ColorYellow = "\u001b[33m"
	ColorBlue   = "\u001b[34m"
	ColorReset  = "\u001b[0m"
)

func Warning(message string) {
	fmt.Println(string(ColorYellow), message, string(ColorReset))
}

func Break(message string) {
	fmt.Println(string(ColorRed), message, string(ColorReset))
	os.Exit(0)
}

func Pass(message string) {
	fmt.Println(string(ColorGreen), message, string(ColorReset))

}

func Help(message string) {
	fmt.Println(string(ColorBlue), message, string(ColorReset))
}

func Red(message string) {
	fmt.Println(string(ColorRed), message, string(ColorReset))
}

func Count_pos() int {
	var arg_len int = len(os.Args[1:])
	return arg_len
}

// help := "-h"
// version := "-v"
// initialize := "-i"
// test -t
// read := "-r"
// write := "-w"
// destroy := "-d"
// uninstall := "--uninstall"
// update := "--update"

func Encrypt_san(arguments int, filename string, object_name string, object_owner string) string {

	Warning(filename)
	Warning(object_name)
	Warning(object_owner)
	// Checking if the file name exists

	return "Valid"
}

func Invalid_op() {
	Break("Invalid option or number of arguments given run encore -h for help")
}

func WriteToFile(data string, location string) {
	// preping data
	var d = []byte(data)
	// checking if file exists
	if Existence(location) == true {
		// default is to overwrite the file
		if DeleteFile(location) == true { // file was deleted

			file, err := os.OpenFile(location, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0400)
			if err != nil {
				Break("Failed to create file")
			}
			if _, err := file.Write(d); err != nil {
				file.Close() // ignore error; Write error takes precedence
				log.Fatal(err)
			}
			if err := file.Close(); err != nil {
				Break("File stream incorrectly terminated")
			}
		} else {
			Break("Im lazy the file could not be deleted")
		}

	} else {
		// Nothing needs to be deleted just write it
		file, err := os.OpenFile(location, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0400)
		if err != nil {
			Break("Failed to create file")
		}
		if _, err := file.Write(d); err != nil {
			file.Close() // ignore error; Write error takes precedence
			log.Fatal(err)
		}
		if err := file.Close(); err != nil {
			Break("File stream incorrectly terminated")
		}
	}

}

func DeleteFile(filename string) bool {
	if Existence(filename) == true {
		del := os.Remove(filename)
		if del != nil {
			Warning("File not deleted")
			return false
		}
		return true

	} else {
		var msg string = "File not found : "
		msg += string("'" + filename + "'")
		Warning(msg)
		return true
	}
}

func Existence(filename string) bool {
	_, error := os.Stat(filename)

	// check if error is "file not exists"
	if os.IsNotExist(error) {
		return false
	} else {
		return true
	}
}

// Thanks https://mrwaggel.be/post/generate-md5-hash-of-a-file-in-golang
func Hash_file_md5(filePath string) (string, error) {
	var returnMD5String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}
	defer file.Close()
	var hash = md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}
	var hashInBytes = hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil
}

func Test1() {
	for i := 1; i <= 43000; i++ {
		Pass("Generating key")

		rand.Seed(time.Now().UnixNano())
		letters := "abcdefghijklmnopqrstuvwxyz123456789012345678901324569870"

		keyr := make([]byte, 128)

		for i := range keyr {
			keyr[i] = letters[rand.Int63()%int64(len(letters))]
		}

		key := string(keyr)
		key += "\n"

		Warning(key)

		file, err := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := file.Write([]byte(key)); err != nil {
			file.Close() // ignore error; Write error takes precedence
			log.Fatal(err)
		}
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}
}
