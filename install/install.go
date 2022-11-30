package install

func Update(keyword string) string {
	if keyword == "--force" {
		return "Updating with --FORCE"
	}
	return "Updating"

}

func Uninstall(keyword string) string {
	if keyword == "--force" {
		return "Uninstalling with --FORCE"
	}
	return "Uninstalling"

}
