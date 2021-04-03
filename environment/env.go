package environment

import "os"

var (
	EVERY8D_USERNAME = os.Getenv("EVERY8D_USERNAME")
	EVERY8D_PASSWORD = os.Getenv("EVERY8D_PASSWORD")
)
