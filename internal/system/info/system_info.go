package systemInfo

import "os/user"

func GetSystemInfo() {

}
func getUser() any {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	homeDir := usr.HomeDir
	userName := usr.Username

	return map[string]string{
		"homeDir":  homeDir,
		"userName": userName,
	}

}
