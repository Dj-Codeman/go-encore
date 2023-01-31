package main

import (
	"encoding/json"
	cnf "encore/config"
	enc "encore/encrypt"
	def "encore/functions"
	ins "encore/install"
	sys "encore/system"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
)

// import "flag"
func min_arguments(min int) {
	var current int = sys.Count_Positional_Vars()
	if current < min {
		sys.Invalid_Op()
	}

}

// The user key isnot kept stored on the disk
// there will be a challenge. some data encrypted with the user key
// when the user inputs the pbdkf passowrd the data will run throught the
// decryption steps with the key given
// if the decryption is sucessfull the key is stored till the program has
// finished running

func main() {

	min_arguments(1)
	switch os.Args[1] {
	case "-d":
		min_arguments(3)
		var object_owner string = os.Args[2]
		var object_name string = os.Args[3]

		def.Destroy(object_owner, object_name)

	case "--debug":
		min_arguments(2)
		sys.Pass("Debugging \n")

		// sub switch statement
		switch os.Args[2] {
		case "--current":

			var userkey string = def.Authenthicate()
			var randkey = def.Fetch_keys(strconv.Itoa(rand.Intn(cnf.Key_max - cnf.Key_cur + 1)))

			sys.Dump(def.Fetch_writting_key(randkey, userkey))

		case "--auth":
			sys.Pass("Auth Testing :\n")

			var d1 string = def.Authenthicate()
			sys.Dump(d1)

		case "--map":
			sys.Pass("Unlocking map \n")
			min_arguments(4)

			var object_owner string = os.Args[3]
			var object_name string = os.Args[4]
			sys.Dump(map_dump(object_owner, object_name))

		case "--index":
			sys.Pass("Key pair manipulation \n")
			min_arguments(4)
			var key string = os.Args[4]

			switch os.Args[3] {
			case "--rehash":
				sys.Pass("Rehasing key \n")
				if safe_key(key) {
					key, _ := strconv.Atoi(os.Args[4])
					if rehash(key) {
						sys.Pass("DONE \n")
					}
				} else if key == "userkey" {
					if rehash_userkey() {
						sys.Pass("DONE \n")
					}
				}

			case "--regen":
				if safty_prompt() {
					if safe_key(key) {
						// Re write key
						var key_data string = enc.Create_key()
						var key_path string = cnf.Keydir + "/" + string(key) + ".dk"
						sys.WriteToFile(key_data, key_path, "write")
						key_int, _ := strconv.Atoi(key)
						if rehash(key_int) {
							sys.Pass("DONE \n")
						}
					}
				}

			case "--divert":
				if safe_key(key) {
					def.Relazy()
				}
			}

		case "--new-sanity":
			bytes, _ := ioutil.ReadFile("/opt/go-encryption-core/test.file")
			sanity_data := string(bytes)
			sys.Dump(enc.Encrypt(sanity_data, "j416wlr6345t331a74sp2iua69660886"))

		case "--dump-userkey":
			sys.Pass("Dumping Public Secret")
			sys.Dump(def.Authenthicate())

		default:
			def.Relazy()
			// Make a debug help adding auth where possible
		}

	case "-h":
		min_arguments(1)
		def.Show_help()

	case "-i":
		min_arguments(1)
		def.Initialize()

	case "-r":
		min_arguments(3)
		var object_owner string = os.Args[2]
		var object_name string = os.Args[3]

		if def.Read(object_owner, object_name) {
			sys.Pass("DONE")
		}

	case "-s": // Sanity test
		_ = def.Authenthicate()
		var key_data string = "j416wlr6345t331a74sp2iua69660886"
		var msg_data string = "21a4bfc92ec6978593228ca65588a48f4b3dd7c0059109780c45d59197a690939dc14cbb5031ea3d09821053025412a2afff636afc0a61ce8508911b4a6b11d8f0df9382ece846559091584f5c85889bb8b5e16be76d25a99c3db030e32f1e59d165c4337e574cdcc528a200bf726bcfe923fb789e4ee6a1f6cfe09fe142058550064de82780b596a7ba50b86f919df95b0fbe70c3d719ac4593a4c6c8655dfde976b63b8b850f0d503f9a8d3f06e1d9084a671caf92bf517eccf7813c7b9f5f3d6cd8102b510abd857eab489c879aeba0e53892bd9c2480c4d3a40c411ade2984b73f13dca064ebae4596ac337cda5083983f5e6f704d14fabc6f7e2edaf82e44837da19f463528ab24ec851ef9d2af3342c074751b29079b44b41962deccc6c31b96045b537de99c4ab9f3f2a259ee63f47ded95ceaa0338008cc0aea445e1d869430ee6e3f70dab5713e783fa9c2cxzbyoq5uvrte5hyk1f9046371f7ff7b5e91f30b79a14007c4125cc839ecfa5fb6c8117ab1b2ec0902292d98f507598444069a27268db529ba2b4376c8bcfd58ca3363890aba72a97"
		sys.Warning(enc.Decrypt(msg_data, key_data))
		sys.Help("See your not crazy")

		// LEGACY cipher text
		// the new ciphter text has a 512 bit hmac
		// var msg_data string = "8c9df5863be41519bc915585451b5c77acf646dbb83112737c518e155fb113f24821f34687293baa4fb2a57257aba5d3c2123e5afc666f2e87bda91a1536b054e5fa95945b33f66bdd9ab94010813572a5f84e5053d41766535da6ef3744882e05d77c667d4c5f32420c4c07c5f63a25fe326c4a20b4c3356f74c2e78fdd83b70cdb25e34bc96af5a94c8abf1bd12d050f6971be707ebae3b5124f98fe9b2a8095c74b72556483a488f8ac2c76059d4d308ac09190819f91fa1c072ab32d51b40e8de28478d04e419f6185e216eaacbb87f55b171821d53f1c7cf12bca63520c7b4a1dfe255b306581983b9d4435bfed03cff80ce1c0338c3cfabb662bc9944bdece94aedf344d1e49dbaef09da327915e70beaac2a1401778d6c947ac7900e0919766b18df61945a70b7340977959c8036422700ec8c7c15afbe4ebc7be2204608aa6cccf122241f0ed8ccbb8717028mqyuja3r23kwltkdced7be29ecf7b515c6f579e503b3b8691022f91937670c5e65a3992e4642f85e"

	case "-t":
		// right now this is for developemnt but it'll have an extended list of system tests
		_ = def.Authenthicate()
		sys.Warning("This will take 50+ mins. On lower spec machines this might not run")
		sys.Warning("If you don't intend on encrypting files larger than 500mb you don't need to run this")
		sys.Warning("If you are encrypting bigger files just sit back and wait for this the finish")
		enc.Larger_test()

	case "--update":
		min_arguments(1)
		if sys.Count_Positional_Vars() == 2 {
			if os.Args[2] == "--force" {
				var force string = os.Args[2]
				var status string = ins.Update(force)
				sys.Warning(status)
			} else {
				def.Update_Help()
			}
		} else {
			var force string = "nil"
			var status string = ins.Update(force)
			sys.Warning(status)
		}

	case "--uninstall":
		min_arguments(1)
		if sys.Count_Positional_Vars() == 2 {
			if os.Args[2] == "--force" {
				var force string = os.Args[2]
				var status string = ins.Uninstall(force)
				sys.Warning(status)
			} else {
				def.Uninstall_Help()
			}
		} else {
			var force string = "nil"
			var status string = ins.Uninstall(force)
			sys.Warning(status)
		}

	case "-v":
		sys.Pass(def.Version())

	case "-w":
		min_arguments(4)
		var dirty_object_path string = os.Args[2]
		var dirty_object_owner string = os.Args[3]
		var dirty_object_name string = os.Args[4]
		object_path, object_owner, object_name := def.Write_preperation(dirty_object_path, dirty_object_owner, dirty_object_name)

		if def.Write(object_path, object_owner, object_name) {
			sys.Pass("DONE")
		} else {
			sys.Break("An error has occoured")
		}

	default:
		sys.Pass(def.Version())
		def.Show_help()
		sys.Break("")
	}
}

func safty_prompt() bool {
	var msg0 string = "WARNING THIS MAY MAKE DATA UNREADABLE !!!!!!!!"
	sys.Warning(msg0)
	var msg1 string = "TO CONTINUE TYPE: 'BIM BAM COMPUTER GO HAM !**!' "
	sys.Fail(msg1)

	awnser := sys.Input_normal("> ")
	if awnser == "BIM BAM COMPUTER GO HAM !**!" {
		return true
	} else {
		return false
	}
}

func safe_key(key string) bool {
	if key == "systemkey" {
		return false
	} else if key == "userkey" {
		return false
	} else {
		return true
	}
}

func rehash(key int) bool {

	var path string = cnf.Plnjson + "/" + strconv.Itoa(key) + ".json"

	bytes, _ := ioutil.ReadFile(path)
	// creating structure
	var any_key_index def.Key_Index

	msg := json.Unmarshal(bytes, &any_key_index)
	if msg != nil {
		sys.Handle_err(msg, "break")
	}

	hash, msg := sys.Hash_file_md5(any_key_index.Location)

	new_key_index := new(def.Key_Index)

	new_key_index.Key_version = def.Version()
	new_key_index.Number = any_key_index.Number
	new_key_index.Location = any_key_index.Location
	new_key_index.Parent = any_key_index.Parent
	new_key_index.Hash = hash

	new_bytes, _ := json.MarshalIndent(new_key_index, "", "  ")
	sys.WriteToFile(string(new_bytes), path, "write")

	return true
}

func rehash_userkey() bool {

	var path string = cnf.Plnjson + "/userkey.json"

	bytes, _ := ioutil.ReadFile(path)
	// creating structure
	var any_key_index def.Key_Index

	msg := json.Unmarshal(bytes, &any_key_index)
	if msg != nil {
		sys.Handle_err(msg, "break")
	}

	hash, msg := sys.Hash_file_md5(any_key_index.Location)

	new_key_index := new(def.Key_Index)

	new_key_index.Key_version = def.Version()
	new_key_index.Number = any_key_index.Number
	new_key_index.Location = any_key_index.Location
	new_key_index.Parent = any_key_index.Parent
	new_key_index.Hash = hash

	new_bytes, _ := json.MarshalIndent(new_key_index, "", "  ")
	sys.WriteToFile(string(new_bytes), path, "write")

	return true
}

func map_dump(object_owner string, object_name string) string {
	var log string = "Debugging map :"
	log += object_owner + "-" + object_name
	def.Write_log(log)

	// Creating the path to the encrypted json file
	var encrypted_json_path string = cnf.Encjson + "/" + object_owner + "-" + object_name + ".json"

	// decrypting the json data
	// getting the ciphertext from the json
	encrypted_json_bytes, _ := ioutil.ReadFile(encrypted_json_path)
	var encrypted_json_data string = string(encrypted_json_bytes)

	// getting the systemkey
	def.Write_log("Auth request")
	var userkey_data string = def.Authenthicate()

	// Getting the plaintext json
	var decrypted_json_data string = enc.Decrypt(encrypted_json_data, userkey_data)
	// to unmarshall json data the format must be in bytes when passed to the function
	var decrypted_json_bytes []byte = []byte(decrypted_json_data)

	// initializing new strut for the data
	var decryption_index def.Secret_Data_Index

	// unpacking the data to the strut
	msg := json.Unmarshal(decrypted_json_bytes, &decryption_index)
	if msg != nil {

		//! make this cleaner by re writting error handeler
		var log string = msg.Error()
		def.Write_log(log)
		sys.Handle_err(msg, "break")
	}

	Dump := fmt.Sprint(decryption_index)
	return Dump
}
