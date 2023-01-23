package functions

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	cnf "encore/config"
	enc "encore/encrypt"
	sys "encore/system"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Secret_Data_Index struct {
	Version string `json:"version"`
	Name    string `json:"object_name"`
	Owner   string `json:"object_owner"`
	Key     string `json:"key"`
	Uid     string `json:"uid"`
	Path    string `json:"origin_path"`
	Dir     string `json:"secret_path"`
	// maybe later
	// Hash    string `json:"hash"`
}

// creating a slut to contain the read json data
type Key_Index struct {
	Hash        string `json:"hash"`
	Parent      string `json:"parent"`
	Location    string `json:"location"`
	Number      string `json:"number"`
	Key_version string `json:"version"`
}

func Relazy() {
	fmt.Println("")
	// debugging function for if statents
}

func Timestamp() string {
	currentTime := time.Now()
	var timestamp string = currentTime.Format("01-02-2006 3:4:5")
	return timestamp
}

func Start_log() {
	sys.Pass("New log started \n")
	var log_dir string = cnf.Logdir + "/general"
	var msg string = "LOG START @ " + Timestamp() + "\n\n"
	sys.WriteToFile(msg, log_dir, "write")
}

func Write_log(data string) {
	var timestamp string = Timestamp()
	var log_dir string = cnf.Logdir + "/general"
	var log_data string = data + " @ " + timestamp + "\n\n"

	sys.WriteToFile(log_data, log_dir, "append")
}

func Show_help() {
	// Help just makes things colored blue
	sys.Help("encore [-w] write [-r] read [-d] destroy [-t] additional tests [-i] initialize [--update] [--uninstall] [-v] version ")
	sys.Help("encore -w -i FILENAME [name] [owner]")
	sys.Help("encore -r name owner ")
	sys.Help("encore -d name owner ")
	sys.Warning("encore -t Run additional tests *Normal users <= than 1000 keys* ")
	sys.Warning("and files <= 500mb don't need this option")
	sys.Help("encore -i **WARNING THIS WILL DELETE ANY STORED DATA AND KEYS** ")
	sys.Help("encore update performs system wellness test then downloads the lates version of encore ")
	sys.Help("Uninstall will delete all stored data and binaries associated with encore")
	// fmt.Println("\n")
	os.Exit(0)
}

func Update_Help() {
	sys.Help("The only additional option for update is --force")
	sys.Warning("Using --force may delete data or break this intallating")
	sys.Fail("We've all heard this before but with great power comes ")
	sys.Break("great responsibility. Use this wisley.")
}

func Uninstall_Help() {
	sys.Fail("Make sure all your data has been read from this program")
	sys.Fail("Uninstall will indiscreminantly delete all data, keys, maps")
	sys.Break("and anything else that has been created by the application")
}

func Generate_userkeys() {

	sys.Pass("Setting up authentication")
	Write_log("Creating user key and hash test")

	var userkey_json_directory string = cnf.Plnjson + "/userkey.json"

	Create_password()

	// Getting integrity
	hash, err := sys.Hash_file_md5(cnf.Userkey)
	if err != nil {
		//  This is a warning because there will be an option to ignore checking md5 sums
		sys.Handle_err(err, "warn")
	}

	// Creating the JSON with a strut
	index := new(Key_Index)
	index.Key_version = Version()
	index.Number = "00"
	index.Location = cnf.Userkey
	index.Parent = "master"
	index.Hash = hash

	// write master json
	bytes, _ := json.MarshalIndent(index, "", "  ")

	sys.WriteToFile(string(bytes), userkey_json_directory, "write")
	Write_log("Userkey created")
}

func Generate_keys() {
	sys.Pass("Regenerating keys and indexs")
	Write_log("Recreating keys and jsons")

	var master_json_directory string = cnf.Plnjson + "/master.json"

	// add part to generate systemkey
	var key string = enc.Create_key()
	sys.WriteToFile(key, cnf.Systemkey, "write")

	// Getting integrity
	hash, err := sys.Hash_file_md5(cnf.Systemkey)
	if err != nil {
		//  This is a warning because there will be an option to ignore checking md5 sums
		sys.Handle_err(err, "warn")
	}

	// Creating the JSON with a strut
	index := new(Key_Index)
	index.Key_version = Version()
	index.Number = "0"
	index.Location = cnf.Systemkey
	index.Parent = "-"
	index.Hash = hash

	// write master json
	bytes, _ := json.MarshalIndent(index, "", "  ")

	sys.WriteToFile(string(bytes), master_json_directory, "write")
	Write_log("Systemkey created: id:0")

	for i := cnf.Key_cur; i <= cnf.Key_max; i++ {
		// Delete keys
		var key_path string = cnf.Keydir + "/" + strconv.Itoa(i) + ".dk"
		var index_path string = cnf.Plnjson + "/" + strconv.Itoa(i) + ".json"

		// Recreating
		var key string = enc.Create_key()
		sys.WriteToFile(key, key_path, "write")

		// Getting integrity
		hash, msg := sys.Hash_file_md5(key_path)
		if msg != nil {
			sys.Handle_err(err, "warn")
		}

		// Creating the JSON with a strut
		index := new(Key_Index)
		index.Key_version = Version()
		index.Number = strconv.Itoa(i)
		index.Location = key_path
		index.Parent = cnf.Systemkey
		index.Hash = hash

		// write indexdir
		// two space seperationg
		bytes, _ := json.MarshalIndent(index, "", "  ")
		sys.WriteToFile(string(bytes), index_path, "write")

		var log string = "Key pair : "
		log += index.Number
		log += " created"

		Write_log(log)

	}
	sys.Pass("DONE \n")
}

func find_de_way(key string) string {
	if key == "systemkey" {
		var key_index_path string = cnf.Plnjson + "/master.json"
		return key_index_path

	} else if key == "userkey" {
		var key_index_path string = cnf.Plnjson + "/userkey.json"
		return key_index_path

	} else {
		var key_index_path string = cnf.Plnjson + "/" + key + ".json"
		return key_index_path
	}
}

func Fetch_keys(key string) string {

	var schrodingers_path string = find_de_way(key)

	// Loading the json file
	bytes, _ := ioutil.ReadFile(schrodingers_path)
	// creating structure
	var any_key_index Key_Index

	msg := json.Unmarshal(bytes, &any_key_index)
	if msg != nil {
		sys.Handle_err(msg, "break")
	}

	if any_key_index.Key_version != Version() {
		sys.Warning("\n Mismatch version warning: The version of encore used by this key is not the same")
		sys.Warning("as the application version. Extract this data and re-initialize to garuntee")
		sys.Warning("Data safety. Or don't ¯\\_(ツ)_/¯ \n")

		Write_log("KEY FETCH WARNING: VERSION MISMATCH SAVE DATA AND REINITIALIZE")
	}

	// checking the hashes
	new_hash, err := sys.Hash_file_md5(any_key_index.Location)
	if err != nil {
		//  This is a warning because there will be an option to ignore checking md5 sums
		sys.Handle_err(err, "warn")
	}

	if new_hash != any_key_index.Hash {
		var log string = "BREAKING FAULT: MISMATCHED HASH ON KEY: "
		log += key
		log += " Regenerate hash with INSERT DEBUG OPTION or delete key"
		Write_log(log)

		sys.Fail("HASH FAULT: The hash stored with this key does not match the")
		sys.Break("not match the current file hash. KEYS HAVE BEEN TAMPERED WITH")
	}

	key_bytes, _ := ioutil.ReadFile(any_key_index.Location)
	return string(key_bytes)

}

func Read(object_owner string, object_name string) bool {

	var log string = "Reading :"
	log += object_owner + "-" + object_name
	Write_log(log)

	// Creating the path to the encrypted json file
	var encrypted_json_path string = cnf.Encjson + "/" + object_owner + "-" + object_name + ".json"

	// decrypting the json data
	// getting the ciphertext from the json
	encrypted_json_bytes, _ := ioutil.ReadFile(encrypted_json_path)
	var encrypted_json_data string = string(encrypted_json_bytes)

	// getting the systemkey
	Write_log("Auth request")
	var userkey_data string = Authenthicate()

	// Getting the plaintext json
	var decrypted_json_data string = enc.Decrypt(encrypted_json_data, userkey_data)
	// to unmarshall json data the format must be in bytes when passed to the function
	var decrypted_json_bytes []byte = []byte(decrypted_json_data)

	// initializing new strut for the data
	var decryption_index Secret_Data_Index

	// unpacking the data to the strut
	msg := json.Unmarshal(decrypted_json_bytes, &decryption_index)
	if msg != nil {

		//! make this cleaner by re writting error handeler
		var log string = msg.Error()
		Write_log(log)
		sys.Handle_err(msg, "break")
	}

	// getting variables

	// key data
	var index_key string = decryption_index.Key
	var index_key_data string = Fetch_keys(index_key)

	// file data
	encrypted_data_bytes, _ := ioutil.ReadFile(decryption_index.Dir)
	var index_encrypted_data string = string(encrypted_data_bytes)

	// decrypting data
	var decrypted_secret_data string = enc.Decrypt(index_encrypted_data, index_key_data)

	// depending on config file replace the original file
	if cnf.Re_place {
		var index_plain_directory string = decryption_index.Path
		sys.WriteToFile(decrypted_secret_data, index_plain_directory, "write")

		//logging
		var log string = "Decrypted :"
		log += object_owner + "-" + object_name + "to: " + index_plain_directory
		Write_log(log)

		// add function to save and retrive the permissions and ownership
		// set_uid(index_plain_directory)
	} else {
		var index_plain_directory string = cnf.Datadir + "/" + decryption_index.Owner + "-" + decryption_index.Name
		sys.WriteToFile(decrypted_secret_data, index_plain_directory, "write")
		// set_uid(index_plain_directory)

		//logging
		var log string = "Decrypted :"
		log += object_owner + "-" + object_name + "to: " + index_plain_directory
		Write_log(log)

	}

	return true
}

func Filename_Sanatization(filename string) string {
	// Checking if the file name exists
	if !filepath.IsAbs(filename) {
		if strings.Contains(filename, "./") {
			var new_file_string string = strings.ReplaceAll(filename, "./", "")
			// getting the current working folder
			working_directory, _ := os.Getwd()
			// this should be the path
			var object_path string = working_directory + "/" + new_file_string
			// cheking if the path we created is valid
			if sys.Existence(object_path) {
				return object_path
			} else {
				sys.Warning("Path doesn't exist : " + object_path)
				return "nil"
			}

		} else {
			// getting the current working folder
			working_directory, _ := os.Getwd()

			// tack current working direcroy to file name
			var object_path string = working_directory + "/" + filename
			if sys.Existence(object_path) {
				return object_path
			} else {
				sys.Warning("Path doesn't exist : " + object_path)
				return "nil"
			}

		}
	} else {
		if sys.Existence(filename) {
			return filename
		} else {
			sys.Warning("Path doesn't exist : " + filename)
			return "nil"
		}

	}
}

func Write_preperation(dirty_object_path string, dirty_object_owner string, dirty_object_name string) (object_path string, object_owner string, object_name string) {

	// checking if the filename has been validated
	var clean_object_path string = Filename_Sanatization(dirty_object_path)
	if clean_object_path != "nil" {
		// testing if the map exists
		var map_test string = cnf.Encjson + "/" + dirty_object_owner + "_" + dirty_object_name + ".json"
		if sys.Existence(map_test) {
			sys.Break("Choose a diffrent name")
		} else {
			return clean_object_path, dirty_object_owner, dirty_object_name
		}

	} else {
		sys.Break("Invalid filename given")
		return "", "", ""
	}
	sys.Break("Never thought i'd get this far")
	return "", "", ""
}

func Write(object_path string, object_owner string, object_name string) bool {
	rand.Seed(time.Now().UnixNano())

	//logging
	var log string = "Writting :"
	log += object_path + "to: " + object_owner + object_name
	Write_log(log)

	// turn this into a checksum ???
	var key_int int = rand.Intn(cnf.Key_max - cnf.Key_cur + 1)
	var key_data string = Fetch_keys(strconv.Itoa(key_int))
	var uid_bytes []byte = []byte(key_data)[0:9]
	var uid_data string = base64.StdEncoding.EncodeToString(uid_bytes)
	var encrypted_data_path string = cnf.Datadir + "/" + uid_data

	data_index := new(Secret_Data_Index)
	data_index.Version = Version()
	data_index.Name = object_name
	data_index.Owner = object_owner
	data_index.Key = strconv.Itoa(key_int)
	data_index.Uid = uid_data
	data_index.Path = object_path
	data_index.Dir = encrypted_data_path

	// var plain_json string = cnf.Encjson + "/" + object_owner + "-" + object_name + ".jn"
	var encrypted_json string = cnf.Encjson + "/" + object_owner + "-" + object_name + ".json"

	// writing the index file
	plain_json_bytes, _ := json.MarshalIndent(data_index, "", "  ")
	var plain_file_string string = string(plain_json_bytes)

	// Creating ciphertext from the file we read
	var cipher_plain string = enc.Encrypt(plain_file_string, key_data)
	sys.WriteToFile(cipher_plain, encrypted_data_path, "write")
	if !cnf.Soft_move {
		if !sys.DeleteFile(object_path) {
			sys.Break("File wasn't deleted idk how tf you got here")
		}
	}

	// Generating the data in the correct formats
	var plain_json_data string = string(plain_json_bytes)
	var userkey_data string = Authenthicate()

	// Creating ciphertext from the plain json map
	var cipher_json string = enc.Encrypt(plain_json_data, userkey_data)
	sys.WriteToFile(cipher_json, encrypted_json, "write")

	Write_log("Write successful")
	return true

}

func Destroy(object_owner string, object_name string) bool {

	//logging
	var log string = "Deleting :"
	log += object_owner + "-" + object_name
	Write_log(log)

	// Creating the path to the encrypted json file
	var encrypted_json_path string = cnf.Encjson + "/" + object_owner + "-" + object_name + ".json"
	// decrypting the json data
	// getting the ciphertext from the json
	encrypted_json_bytes, _ := ioutil.ReadFile(encrypted_json_path)
	var encrypted_json_data string = string(encrypted_json_bytes)

	// getting the systemkey
	var userkey_data string = Authenthicate()

	// Getting the plaintext json
	var decrypted_json_data string = enc.Decrypt(encrypted_json_data, userkey_data)
	// to unmarshall json data the format must be in bytes when passed to the function
	var decrypted_json_bytes []byte = []byte(decrypted_json_data)

	// initializing new strut for the data
	var decryption_index Secret_Data_Index

	// unpacking the data to the strut
	msg := json.Unmarshal(decrypted_json_bytes, &decryption_index)
	if msg != nil {
		sys.Handle_err(msg, "break")
	}

	if cnf.Leave_in_peace {
		if Read(object_owner, object_name) {
			sys.Pass("FILE READ !\n")
			Write_log("Data saved before the purge")
		} else {
			sys.Fail("FILE WAS NOT READ !\n")
			Write_log("BREAKING FAULT: File was not read ")
			return false
		}
	}

	var file_location string = decryption_index.Dir

	// DELETING files
	if sys.DeleteFile(file_location) {
		Write_log("Deleting the files")
		if sys.DeleteFile(encrypted_json_path) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func Version() string {
	var ver = "G1.00"
	return ver
}

func Initialize() {
	sys.Pass("Running Initialization \n")

	Start_log()

	sys.Pass("Creating Folders \n")
	// making folders for first time run
	folders := [6]string{cnf.Datadir, cnf.Encjson, cnf.Plnjson, cnf.Keydir, cnf.Logdir, "/tmp/encore"}
	for i := 0; i <= len(folders)-1; i++ {
		sys.MakeFolder(folders[i])
	}

	Write_log("Running Encryption Tests")

	status, msg := enc.Test()
	if status == "Pass" {
		sys.Pass("DONE \n")
	} else {
		Write_log(msg)
		sys.Break(msg)
	}

	sys.Pass("Generating keys")
	// The normal system key and n others
	Generate_keys()

	// The user auth key
	Generate_userkeys()
	// this fuction handles the ouutput

	sys.Pass("Running Key Fetch Functionallity Test")

	var Random_Key int = rand.Intn(cnf.Key_max - cnf.Key_cur + 1)
	var msg1 string = "Random key fetched : " + Fetch_keys(fmt.Sprint(Random_Key))

	sys.Pass(msg1)
	sys.Pass("DONE \n")

	Write_log("Finished initialization")
	var t string = "Initialization Finished @ " + Timestamp() + "\n"
	sys.Pass(t)
}

func Create_password() {

	// Gathering and validating password
	var psswd_1 string = sys.Input_secret(" Pick a passowrd for the locker : ")
	sys.Pass("DONE")

	if len(psswd_1) <= 1 {
		sys.Break("Password must be between 1 - 255 charathers")
	} else if len(psswd_1) >= 256 {
		sys.Break("Password must be between 1 - 255 charathers")
	} else {
		Relazy()
	}

	var psswd_2 string = sys.Input_secret(" Please re-type password : ")

	if psswd_1 == psswd_2 {
		sys.Pass("DONE \n")
	} else {
		sys.Break("\n\n Password Invalid ! \n")
	}

	// Converting to bytes to d to k from the p
	var password_bytes []byte = []byte(psswd_1)

	//running the password derived from key function
	var pdk []byte = enc.Pbkdf(password_bytes, []byte(Fetch_keys("systemkey")), 12400000, 16, sha512.New)

	// Converting to hex to keep the storage format the same
	var password_key = hex.EncodeToString(pdk)

	// Making the secret to check aginst
	var password_check string = enc.Encrypt("The hotdog man isn't real !?", password_key)

	//Storing the data
	sys.WriteToFile(password_check, cnf.Userkey, "write")
	Write_log("Password and hash created")
}

func Check_password() string {
	var password string = sys.Input_secret("Input the locker password : ")

	var password_bytes []byte = []byte(password)

	//running the password derived from key function
	var pdk []byte = enc.Pbkdf(password_bytes, []byte(Fetch_keys("systemkey")), 12400000, 16, sha512.New)
	var password_key = hex.EncodeToString(pdk)

	// use the fetch key function
	var verification_ciphertext string = Fetch_keys("userkey")
	var verification_string string = "The hotdog man isn't real !?"

	var verification_result string = enc.Decrypt(string(verification_ciphertext), password_key)

	if verification_result == verification_string {
		sys.Pass("Locker unlocked \n")
		return password_key
	} else {
		return "nil"
	}
}

func Authenthicate() string {

	// Having set uid for the program
	var uid int = os.Geteuid()

	// Verify the user to the most high
	if uid != 0 {
		var message string = "Invalid UID your not that guy\n"
		sys.Break(message)
		return ""

	} else {
		var message string = "Valid user !\n"
		sys.Pass(message)
		var userkey string = Check_password()

		if userkey != "nil" {
			sys.Pass("User authenticated ! \n")
			return userkey
		} else {
			sys.Break("Authentication Failed ! \n")
			return ""
		}
	}
}
