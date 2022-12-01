package config

const (
	// LOCATIONS
	//  Data this is where the finished and encrypted files live
	//  When keys are regenerated this folder will be emptyied
	//  default /var/encore/data

	Datadir string = "/var/encore/data"

	//  JSON This is where plan text maps will live
	//  these are generated along side the keys
	//  default /var/encore/indexs

	Plnjson string = "/var/encore/indexs"

	//  This is where the encrypted jsons for written file
	//  will live. The json debug tool should be used to decrypt
	//  and modify these files

	encjson string = "/var/encore/maps"

	//  KEY These are the random encryption keys
	//  128 bit strings for use with the encrypt script
	//  https://www.fastsitephp.com/fr/documents/file-encryption-bash
	//  default /opt/encore/keys

	Keydir string = "/var/encore/keyz"

	//  SYSTEM KEY JSON file that contain location and key information
	//  are encrypted using this key
	//  if this key is missing on script call all file in:
	//  $datadir will be illegible
	//  IF THIS KEY IS DELETED ALL DATA IS CONSIDERED LOST
	//  default /opt/encore/keys/systemkey.dk

	Systemkey string = "/etc/systemkey1.dk"

	// log dir

	Logdir string = "/var/log/encore/general"

	//  key_max the limit of keys to generate
	//  default=1000

	Key_max int = 1

	//  Works like a key min value
	//  by key_cur and key_max the range from which keys are picked
	//  can be changed

	Key_cur int = 0

	//  soft moving
	//  set 1 to use cp instead of mv when gatheing files to encrypt
	//  default = false

	Soft_move bool = false

	//  re-place file
	//  the original path of files are stored when encrypted
	//  if set files will be re placed back in there original
	//  directory
	//  default= true

	Re_place bool = true

	//  save on destroy
	//  if you want the destroy function to recover the file before deleting
	//  the encrypted copy set this to 1
	//  default=1

	Leave_in_peace bool = true
)
