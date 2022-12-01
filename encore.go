package main

import (
	def "encore/functions"
	ins "encore/install"
	sys "encore/system"
	"os"
)

//import "flag"

func main() {
	uid := os.Geteuid()

	if uid != 0 {
		var message = "You shall not pass."
		sys.Break(message)
	}

	arguments := sys.Count_pos()

	if arguments == 0 {
		sys.Invalid_op()

	} else if arguments == 1 {
		flag1 := os.Args[1]

		if flag1 == "-i" {
			def.Initialize()

		} else if flag1 == "-v" {
			version := "Version: "
			version += def.Version()
			sys.Pass(version)

		} else if flag1 == "-h" {
			def.Show_help()

		} else if flag1 == "--update" {
			sys.Warning(ins.Update("none"))

		} else if flag1 == "--uninstall" {
			sys.Break(ins.Uninstall("none"))

		} else if flag1 == "-t" {
			// sys.Test1()
			// enc.Encrypt()
			// sys.Pass(enc.Encrypt("This is a very secret message, Treat this with care"))
			// ciphertest, key := enc.Encrypt(" 0dg3edetyhtiyhgzovglekukcpqy2ird5qpk7o1getdkmbjt659oivqs8z7un0y5220nx0cto4dug0hf18xh8ohiloc9zb342mfisub35ai1300agujx5bwxwqpnguw3")
			// sys.Pass(enc.Decrypt(ciphertest, key))
			// 1 def.Write()
			// var data string = "e7fdc50d0142fdb10453b642d7ab5f687vvaw84zhdfu7nczbf639016179fa5fa2f1296ba2873f4f67428c387573a6723aeb4ed4ef0ec7d22"
			// var key string = "j416wlr6345t331a74sp2iua69660886"
			// def.Read(data, key)
			def.Start_log()
			def.Write_log("hello world")
			// enc.Decrypt("2a07e90227936dc7e5b5b43d193aad955b6df34eecc4478393c1d70b3b3520586c343867386c616b6632346b3832316893c0c054c3e35a451cd3067e04decd0ba1789407da0d074061ee62a36f7e3c95", "9134425d9lc6e8t4sg7egm0135trx2w9")
		} else {
			sys.Invalid_op()
		}

	} else if arguments == 2 {
		flag1 := os.Args[1]
		flag3 := os.Args[2]

		if flag3 != "--force" {
			sys.Break("If you aren't using --force this is  not for you")

		} else {

			if flag1 == "--update" {
				sys.Warning(ins.Update(flag3))

			} else if flag1 == "--uninstall" {
				sys.Break(ins.Uninstall(flag3))

			} else {
				sys.Invalid_op()

			}

		}

	} else if arguments == 3 {
		flag1 := os.Args[1]

		if flag1 == "-r" {
			def.Relazy()
		} else if flag1 == "-d" {
			def.Relazy()
		} else {
			sys.Invalid_op()
		}

	} else if arguments == 4 {
		flag1 := os.Args[1]

		if flag1 == "-w" {

			filename := os.Args[1]
			object_path := os.Args[2]
			object_name := os.Args[3]

			san_result := sys.Encrypt_san(arguments, filename, object_path, object_name)

			if san_result == "invalid" {
				sys.Break("Write test failed !")
			} else {
				sys.Pass("Write test passed")
				// enc.Encrypt("Hello World", dump string, "m09558ahlh891066f02l050ozwt0m30pnva8smz6i1la8374hx279q90t19m9a045gam6d72e242m3nv6784hl5un5lk08fe1ig4r32ok7m4do3jwxv3lx8t7w2tg55d" )
			}

		} else {
			sys.Invalid_op()

		}

	} else {
		sys.Invalid_op()

	}

	os.Exit(0)
}
