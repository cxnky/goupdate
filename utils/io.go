package utils

import (
	"crypto/sha256"
	"io/ioutil"
	"encoding/hex"
	"github.com/cxnky/goupdate/errors"
)

/**

 * Created by cxnky on 24/08/2018 at 18:13
 * utils
 * https://github.com/cxnky/
 
**/

func ValidateChecksum(filePath, expectedChecksum string) bool {

	hasher := sha256.New()
	s, err := ioutil.ReadFile(filePath)
	hasher.Write(s)

	if err != nil {

		errors.NewError("unable to validate checksum of file")
		return false

	}

	actualChecksum := hex.EncodeToString(hasher.Sum(nil))

	return actualChecksum == expectedChecksum

}