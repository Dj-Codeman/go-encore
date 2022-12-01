package functions

import (
	"encoding/json"
	enc "encore/encrypt"
	sys "encore/system"
	"io/ioutil"
	"strconv"
	"time"

	cnf "encore/config"
	"fmt"
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
	sys.Help("encore [-w] write [-r] read [-d] destroy [-i] initialize [--update] [--uninstall] [-v] version ")
	sys.Help("encore -w -i FILENAME [name] [owner]")
	sys.Help("encore -r name owner ")
	sys.Help("encore -d name owner ")
	sys.Help("encore -i **WARNING THIS WILL DELETE ANY STORED DATA AND KEYS** ")
	sys.Help("encore update performs system wellness test then downloads the lates version of encore ")
	sys.Help("Uninstall will delete all stored data and binaries associated with encore")
	// fmt.Println("\n")

}

func Generate_keys() {
	sys.Pass("Regenerating keys and indexs")

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

	// Creating the JSON
	var index = map[string]string{
		"version":  Version(),
		"number":   "0",
		"location": cnf.Systemkey,
		"parent":   "-",
		"hash":     hash,
	}

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

		// Creating the JSONs
		var index = map[string]string{
			"version":  Version(),
			"number":   strconv.Itoa(i),
			"location": key_path,
			"parent":   cnf.Systemkey,
			"hash":     hash,
		}

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
	type key_index struct {
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
	var any_key_index key_index

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

func Write() {
	// this is the part that writes stuff
	var data string = "Hello world"
	var key string = enc.Create_key()
	sys.Pass(enc.Encrypt(data, key))
	sys.Pass(key)
}

func Destroy() {
	Relazy()
}

func Version() string {
	var ver = "G0.00"
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
