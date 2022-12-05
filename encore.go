package main

import (
	enc "encore/encrypt"
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
		sys.Break("Not implemented")

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
		var key_data string = "j416wlr6345t331a74sp2iua69660886"
		var msg_data string = "8c9df5863be41519bc915585451b5c77acf646dbb83112737c518e155fb113f24821f34687293baa4fb2a57257aba5d3c2123e5afc666f2e87bda91a1536b054e5fa95945b33f66bdd9ab94010813572a5f84e5053d41766535da6ef3744882e05d77c667d4c5f32420c4c07c5f63a25fe326c4a20b4c3356f74c2e78fdd83b70cdb25e34bc96af5a94c8abf1bd12d050f6971be707ebae3b5124f98fe9b2a8095c74b72556483a488f8ac2c76059d4d308ac09190819f91fa1c072ab32d51b40e8de28478d04e419f6185e216eaacbb87f55b171821d53f1c7cf12bca63520c7b4a1dfe255b306581983b9d4435bfed03cff80ce1c0338c3cfabb662bc9944bdece94aedf344d1e49dbaef09da327915e70beaac2a1401778d6c947ac7900e0919766b18df61945a70b7340977959c8036422700ec8c7c15afbe4ebc7be2204608aa6cccf122241f0ed8ccbb8717028mqyuja3r23kwltkdced7be29ecf7b515c6f579e503b3b8691022f91937670c5e65a3992e4642f85e"
		sys.Warning(enc.Decrypt(msg_data, key_data))
		sys.Help("See your not crazy")

	case "-t":
		// right now this is for developemnt but it'll have an extended list of system tests
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
