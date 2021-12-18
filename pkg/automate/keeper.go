package automate

import "os"

func getKeeperTagName() string {
	k, ok := os.LookupEnv("KEEPER_TAG_NAME")
	if !ok {
		k = "com.thorsten-hans.keeper"
	}
	return k
}
