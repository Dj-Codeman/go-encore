package system

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
)

const (
	ColorBlack  = "\u001b[30m"
	ColorRed    = "\u001b[31m"
	ColorGreen  = "\u001b[32m"
	ColorYellow = "\u001b[33m"
	ColorBlue   = "\u001b[34m"
	ColorBold   = "\x1B[1m"
	ColorReset  = "\u001b[0m"
)

func Warning(message string) {
	fmt.Println(string(ColorYellow), message, string(ColorReset))
}

func Break(message string) {
	fmt.Println(string(ColorRed), message, string(ColorReset))
	os.Exit(1)
}

func Pass(message string) {
	fmt.Println(string(ColorGreen), message, string(ColorReset))
}

func Fail(message string) {
	fmt.Println(string(ColorRed), message, string(ColorReset))
}

func Help(message string) {
	fmt.Println(string(ColorBlue), message, string(ColorReset))
}

func Dump(message string) {
	fmt.Println(string(ColorBlue), string(ColorBold), message, string(ColorReset))
	os.Exit(0)
}

func Count_Positional_Vars() int {
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

func Invalid_Op() {
	Break("Invalid option or number of arguments given run encore -h for help")
}

func WriteToFile(data string, location string, append string) {
	// preping data
	var d = []byte(data)
	// checking if file exists
	if append == "write" {
		if Existence(location) {
			// default is to overwrite the file
			if DeleteFile(location) { // file was deleted

				file, err := os.OpenFile(location, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0400)
				if err != nil {
					Handle_err(err, "break")
				}
				if _, err := file.Write(d); err != nil {
					file.Close() // ignore error; Write error takes precedence
					Handle_err(err, "break")
				}
				if err := file.Close(); err != nil {
					Handle_err(err, "break")
				}
			} else {
				Break("Im lazy the file could not be deleted")
			}

		} else {
			// Nothing needs to be deleted just write it
			file, err := os.OpenFile(location, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0400)
			if err != nil {
				Handle_err(err, "break")
			}
			if _, err := file.Write(d); err != nil {
				file.Close() // ignore error; Write error takes precedence
				Handle_err(err, "break")
			}
			if err := file.Close(); err != nil {
				Break("File stream incorrectly terminated")
			}
		}

	} else if append == "append" {
		// Nothing needs to be deleted just write it
		file, err := os.OpenFile(location, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0400)
		if err != nil {
			Handle_err(err, "break")
		}
		if _, err := file.Write(d); err != nil {
			file.Close() // ignore error; Write error takes precedence
			Handle_err(err, "break")
		}
		if err := file.Close(); err != nil {
			Break("File stream incorrectly terminated")
		}

	} else {
		Warning("Invalid option given")
	}

}

func DeleteFile(filename string) bool {
	if Existence(filename) {
		del := os.Remove(filename)
		if del != nil {
			Handle_err(del, "warn")
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

func MakeFolder(path string) bool {
	if !Existence(path) {
		// folder doesn't exist make one
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			Break("An error has occoured can't make folder : " + path)
			return false
		} else {
			return true
		}
	} else {
		Warning("Path exists : " + path)
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

func Handle_err(msg error, action string) {

	var error_message string = msg.Error()
	if action == "break" {
		Break(error_message)
	} else if action == "warn" {
		Warning(error_message)
	} else {
		Warning(error_message)
	}

}

func Input_normal(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

func Input_secret(prompt string) string {
	fmt.Print(prompt)

	// Common settings and variables for both stty calls.
	attrs := syscall.ProcAttr{
		Dir:   "",
		Env:   []string{},
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
		Sys:   nil}
	var ws syscall.WaitStatus

	// Disable echoing.
	pid, err := syscall.ForkExec(
		"/bin/stty",
		[]string{"stty", "-echo"},
		&attrs)
	if err != nil {
		panic(err)
	}

	// Wait for the stty process to complete.
	_, err = syscall.Wait4(pid, &ws, 0, nil)
	if err != nil {
		panic(err)
	}

	// Echo is disabled, now grab the data.
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	// Re-enable echo.
	pid, err = syscall.ForkExec(
		"/bin/stty",
		[]string{"stty", "echo"},
		&attrs)
	if err != nil {
		panic(err)
	}

	// Wait for the stty process to complete.
	_, err = syscall.Wait4(pid, &ws, 0, nil)
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(text)
}
