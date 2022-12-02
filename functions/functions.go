package functions

import (
	"encoding/base64"
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

func Relazy() {
	fmt.Println("")
}

func Timestamp() string {
	currentTime := time.Now()
	var timestamp string = currentTime.Format("01-02-2006 3:4:5")
	return timestamp
}

func Start_log() {
	sys.Pass("New log started")
	var msg string = "LOG START @ " + Timestamp() + "\n\n"
	sys.WriteToFile(msg, cnf.Logdir, "write")
}

func Write_log(data string) {
	var timestamp string = Timestamp()
	var log_data string = data + " @ " + timestamp + "\n\n"

	sys.WriteToFile(log_data, cnf.Logdir, "append")
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
	sys.Red("We've all heard this before but with great power comes ")
	sys.Break("great responsibility. Use this wisley.")
}

func Uninstall_Help() {
	sys.Red("Make sure all your data has been read from this program")
	sys.Red("Uninstall will indiscreminantly delete all data, keys, maps")
	sys.Break("and anything else that has been created by the application")
}

func Generate_keys() {
	sys.Pass("Regenerating keys and indexs")
	// creating a slut to contain the read json data
	type Key_Index struct {
		Hash        string `json:"hash"`
		Parent      string `json:"parent"`
		Location    string `json:"location"`
		Number      string `json:"number"`
		Key_version string `json:"version"`
	}

	var master_json_directory string = cnf.Plnjson + "/" + "master.json"

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

	// sys.Handle_err(err, "break")

	sys.WriteToFile(string(bytes), master_json_directory, "write")

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

	}
	sys.Pass("DONE")
}

func find_de_way(key string) string {
	if key == "systemkey" {
		var key_index_path string = cnf.Plnjson + "/master.json"
		return key_index_path
	} else {
		var key_index_path string = cnf.Plnjson + "/" + key + ".json"
		return key_index_path
	}
}

func Fetch_keys(key string) string {
	// creating a slut to contain the read json data
	type Key_Index struct {
		Hash        string `json:"hash"`
		Parent      string `json:"parent"`
		Location    string `json:"location"`
		Number      string `json:"number"`
		Key_version string `json:"version"`
	}

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
		sys.Warning("Mismatch version warning: The version of encore used by this key is not the same")
		sys.Warning("as the application version. Extract this data and re-initialize to garuntee")
		sys.Warning("Data safety. Or don't ¯\\_(ツ)_/¯")
	}

	// checking the hashes
	new_hash, err := sys.Hash_file_md5(any_key_index.Location)
	if err != nil {
		//  This is a warning because there will be an option to ignore checking md5 sums
		sys.Handle_err(err, "warn")
	}

	if new_hash != any_key_index.Hash {
		sys.Red("HASH FAULT: The key hash associated with the key does")
		sys.Break("not match with the current hash. KEYS HAVE BEEN TAMPERED WITH")
	}

	key_bytes, _ := ioutil.ReadFile(any_key_index.Location)
	return string(key_bytes)

}

func Read(data string, key string) bool {
	// This is the part that will read stuff
	// do some validation

	// create the temporary file in datadir

	// decrypt the data

	//  write to the temp file

	// depending on config file replace the original file

	// return bool for status

	var real_shit string = enc.Decrypt(data, key)
	sys.Pass(real_shit)
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

	type Secret_Data_Index struct {
		Version string `json:"version"`
		Name    string `json:"object_name"`
		Owner   string `json:"object_owner"`
		Key     string `json:"key"`
		Uid     string `json:"uid"`
		Path    string `json:"origin_path"`
		Dir     string `json:"data_path"`
		// maybe later
		// Hash    string `json:"hash"`
	}

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
	fmt.Print(data_index)

	// sys.Pass(enc.Encrypt(data, key))
	// sys.Pass(key)
	return true
}

func Destroy() {
	Relazy()
}

func Version() string {
	var ver = "G0.02"
	return ver
}

func Initialize() {
	sys.Pass("Running Initialization")

	Start_log()

	Write_log("Started initialization")

	status, msg := enc.Test()
	if status == "Pass" {
		sys.Pass("DONE")
	} else {
		Write_log(msg)
		sys.Break(msg)
	}

	Generate_keys()

	sys.Pass("Testing key fetch functionality")
	// make this a rand int
	var msg1 string = "Random key fetched : " + Fetch_keys("5")
	sys.Pass(msg1)
	sys.Pass("DONE")

	Write_log("Finished initialization")
	sys.Pass("Initialization Finished")
}
