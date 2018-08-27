package utils

import (
	"runtime"
	"path/filepath"
	"os"
	"github.com/cxnky/goupdate/errors"
)

/**

 * Created by cxnky on 23/08/2018 at 15:29
 * utils
 * https://github.com/cxnky/
 
**/

// GetOSName returns the name of the current operating system
func GetOSName() string {

	switch runtime.GOOS {

	case "darwin":
		return "osx"

	case "windows":
		return runtime.GOOS

	case "linux":
		return runtime.GOOS

	default:
		return runtime.GOOS

	}

}

// GetPWD returns the present working directory (for updating binaries)
func GetPWD() string {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {

		panic(errors.NewError("FATAL: " + err.Error() + " when getting PWD"))

	}

	return dir

}