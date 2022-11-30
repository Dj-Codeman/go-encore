package functions

import (
	"encoding/json"
	enc "encore/encrypt"
	sys "encore/system"
	"strconv"

	cnf "encore/config"
	"fmt"
)

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

func generate_keys() {
	sys.Warning("Regenerating keys and indexs")

	masterdir := cnf.Plnjson + "/" + "master.json"

	// Deleting systemkey
	sys.DeleteFile(cnf.Systemkey)
	sys.DeleteFile(masterdir)

	// add part to generate systemkey
	key := enc.Create_key()
	sys.WriteToFile(key, cnf.Systemkey)

	// Getting integrity
	hash, msg := sys.Hash_file_md5(cnf.Systemkey)
	if msg != nil {
		sys.Warning(msg.Error())
	}

	// Creating the JSON
	index := map[string]string{
		"version":  Version(),
		"number":   "0",
		"location": cnf.Systemkey,
		"hash":     hash,
	}

	// write master json
	bytes, _ := json.MarshalIndent(index, "", "\t")
	sys.WriteToFile(string(bytes), masterdir)

	for i := cnf.Key_cur; i <= cnf.Key_max; i++ {
		// Delete keys
		keydir := cnf.Keydir + "/" + strconv.Itoa(i) + ".dk"
		indexdir := cnf.Plnjson + "/" + strconv.Itoa(i) + ".json"

		sys.DeleteFile(keydir)
		sys.DeleteFile(indexdir)

		// Recreating
		key := enc.Create_key()
		sys.WriteToFile(key, keydir)

		// Getting integrity
		hash, msg := sys.Hash_file_md5(keydir)
		if msg != nil {
			sys.Warning(msg.Error())
		}

		// Creating the JSONs
		index := map[string]string{
			"version":  Version(),
			"number":   strconv.Itoa(i),
			"location": keydir,
			"parent":   cnf.Systemkey,
			"hash":     hash,
		}

		// write indexdir
		// two space seperationg
		bytes, _ := json.MarshalIndent(index, "", "  ")
		sys.WriteToFile(string(bytes), indexdir)

	}
}

func Read() {
	// This is the part that will read stuff
	var data string = "e7fdc50d0142fdb10453b642d7ab5f687vvaw84zhdfu7nczbf639016179fa5fa2f1296ba2873f4f67428c387573a6723aeb4ed4ef0ec7d22"
	var key string = "j416wlr6345t331a74sp2iua69660886"
	var real_shit string = enc.Decrypt(data, key)
	sys.Pass(real_shit)
}

func Write() {
	// this is the part that writes stuff
	var data string = "Hello world"
	var key string = enc.Create_key()
	sys.Pass(enc.Encrypt(data, key))
	sys.Pass(key)
}

func Version() string {
	var ver = "G0.00"
	return ver
}

func Initialize() {
	sys.Pass("Running Initialization")

	generate_keys()

	Relazy()
}

func Relazy() {
	fmt.Println("")

}
