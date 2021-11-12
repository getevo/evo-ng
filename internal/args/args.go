package args

import "os"

func Get(sw string) string {
	for i := 0; i < len(os.Args); i++{
		if os.Args[i] == sw{
			if i+1 < len(os.Args){
				return os.Args[i+1]
			}
		}
	}
	return ""
}

func Exists(sw string) bool {
	for i := 0; i < len(os.Args); i++{
		if os.Args[i] == sw {
			return true
		}
	}
	return false
}