package utils

import (
	"crypto/sha256"
	"io/ioutil"
	"encoding/hex"
	"github.com/cxnky/goupdate/errors"
	"archive/zip"
	"os"
	"path/filepath"
	"io"
)

/**

 * Created by cxnky on 24/08/2018 at 18:13
 * utils
 * https://github.com/cxnky/
 
**/

func Unzip(src, dest string) error {

	r, err := zip.OpenReader(src)

	if err != nil {

		return errors.NewError(err.Error())

	}

	defer func() {

		if err := r.Close(); err != nil {

			panic(err)

		}

	}()

	os.MkdirAll(dest, 0755)

	extractAndWriteFile := func(f *zip.File) error {

		rc, err := f.Open()

		if err != nil {

			return errors.NewError(err.Error())

		}

		defer func() {

			if err := rc.Close(); err != nil {

				panic(err)

			}

		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return errors.NewError(err.Error())
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return errors.NewError(err.Error())
			}
		}
		return nil

	}

	for _, f := range r.File {

		err := extractAndWriteFile(f)

		if err != nil {

			return errors.NewError(err.Error())

		}

	}

	return nil

}

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