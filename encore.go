package main

import (
	def "encore/functions"
	ins "encore/install"
	sys "encore/system"
	"os"
)

// import "flag"
func min_arguments(min int) {
	var current int = sys.Count_Positional_Vars()
	if current < min {
		sys.Invalid_Op()
	}

}

func main() {
	var uid int = os.Geteuid()

	if uid != 0 {
		var message = "You shall not pass."
		sys.Break(message)
	}

	min_arguments(1)
	switch os.Args[1] {
	case "-d":
		min_arguments(3)

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

		def.Read(object_owner, object_name)

	case "-t":
		// sys.Test1()
		// enc.Encrypt()
		// sys.Pass(enc.Encrypt("This is a very secret message, Treat this with care"))
		// ciphertest, key := enc.Encrypt(" 0dg3edetyhtiyhgzovglekukcpqy2ird5qpk7o1getdkmbjt659oivqs8z7un0y5220nx0cto4dug0hf18xh8ohiloc9zb342mfisub35ai1300agujx5bwxwqpnguw3")
		// sys.Pass(enc.Decrypt(ciphertest, key))
		// 1 def.Write()
		// var data string = "e7fdc50d0142fdb10453b642d7ab5f687vvaw84zhdfu7nczbf639016179fa5fa2f1296ba2873f4f67428c387573a6723aeb4ed4ef0ec7d22"
		// var key string = "j416wlr6345t331a74sp2iua69660886"
		// def.Read(data, key)
		// def.Start_log()
		// def.Write_log("hello world")
		// sys.Warning(enc.Test())
		// enc.Decrypt("2a07e90227936dc7e5b5b43d193aad955b6df34eecc4478393c1d70b3b3520586c343867386c616b6632346b3832316893c0c054c3e35a451cd3067e04decd0ba1789407da0d074061ee62a36f7e3c95", "9134425d9lc6e8t4sg7egm0135trx2w9")
		sys.Pass(def.Fetch_keys("99"))

		// right now this is for developemnt but it'll have an extended list of system tests

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
