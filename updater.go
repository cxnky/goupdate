package goupdate

import (
	"net/http"
	"time"
	"io/ioutil"
	"encoding/json"
	"github.com/cxnky/goupdate/utils"
	"os"
	"path/filepath"
	"log"
	"github.com/cxnky/goupdate/errors"
	"fmt"
	"strconv"
	"io"
)

/**

 * Created by cxnky on 23/08/2018 at 14:55
 * goupdate
 * https://github.com/cxnky/

**/

const (

	VERSION = "1.0.0.0"

)

type Updater struct {

	VersionURL			string
	CurrentVersion		string
	CheckFrequency		int
	ShowProgress		bool

	isInitialised		bool
	shouldAutoCheck		bool
	webClient			http.Client
	logger				*log.Logger

}

type updateResponse struct {

	Version		string   `json:"version"`
	Changelog	string 	 `json:"changelog"`
	FileURL		string 	 `json:"url"`
	SHA256Hash 	string	 `json:"hash"`

}

var (

	updateResp = updateResponse{}

)

// CreateUpdater creates and returns a constructed Updater object using a url, current version and check frequency (0 for off)
func CreateUpdater(url, currentVersion string, checkFrequency int, showProgress bool) Updater {

	updater := Updater{isInitialised:true}

	if checkFrequency == 0 {

		updater.shouldAutoCheck = false

	} else {

		updater.shouldAutoCheck = true
		updater.CheckFrequency = checkFrequency

	}

	updater.VersionURL = url
	updater.CurrentVersion = currentVersion

	updater.webClient = http.Client{Timeout: 10 * time.Second}
	updater.ShowProgress = showProgress

	updater.logger = log.New(os.Stderr, "goupdate: ", log.LstdFlags | log.Lshortfile)
	return updater

}

// PerformUpdate downloads the updated file, unzips it and displays the changelog
func (u Updater) PerformUpdate() error {

	fileName := ""
	osys := utils.GetOSName()

	if osys == "windows" {

		fileName = filepath.Base(os.Args[0])

	} else if osys == "linux" {

		fileName = os.Args[0][2:]

	}

	pwd := utils.GetPWD()
	fullPath := pwd + "\\" + fileName

	if osys == "windows" {

		// move file to <name>.bak (windows supports renaming executables whilst running)
		os.Rename(fullPath, fullPath + ".bak")

	} else if osys == "linux" {

		// replace the file (linux supports replacing running executables)
		os.Rename(fullPath, fullPath + "-bak")

	}

	// actually download the update file
	u.downloadUpdate(pwd)

	// validate the checksum of the downloaded file
	if !utils.ValidateChecksum(pwd + "\\update.zip", updateResp.SHA256Hash) {

		return errors.NewError("checksum did not match.")

	}

	return nil

}

func printDownloadPercent(done chan int64, path string, total int64) {

	var stop bool = false

	for {

		select {

			case <-done:
				stop = true

				default:

					file, err := os.Open(path)

					if err != nil {

						errors.NewError(err.Error())
						return

					}

					fi, err := file.Stat()

					if err != nil {

						errors.NewError(err.Error())
						return

					}

					size := fi.Size()

					if size == 0 {

						size = 1

					}

					var percent float64 = float64(size) / float64(total) * 100

					fmt.Printf("%.0f", percent)
					fmt.Println("%")

		}

		if stop {

			break

		}

		time.Sleep(time.Second)

	}

}

// DownloadUpdate will download the zip file from the given download URL and only returns an error if something went wrong
func (u Updater) downloadUpdate(directory string) error {

	if u.ShowProgress {

		fmt.Printf("Downloading update from %s\n", updateResp.FileURL)

		start := time.Now()

		out, err := os.Create(directory + "\\update.zip")

		defer out.Close()

		headResp, err := http.Head(updateResp.FileURL)

		if err != nil {

			panic(errors.NewError(err.Error()))

		}

		defer headResp.Body.Close()

		size, err := strconv.Atoi(headResp.Header.Get("Content-Length"))

		if err != nil {

			panic(errors.NewError(err.Error()))

		}

		done := make(chan int64)

		go printDownloadPercent(done, directory + "\\update.zip", int64(size))

		resp, err := http.Get(updateResp.FileURL)

		if err != nil {

			panic(errors.NewError(err.Error()))

		}

		defer resp.Body.Close()

		n, err := io.Copy(out, resp.Body)

		if err != nil {

			panic(errors.NewError(err.Error()))

		}

		done <- n

		elapsed := time.Since(start)
		fmt.Printf("Update download completed in %s", elapsed)

	} else {

		fmt.Printf("Downloading update from %s\n", updateResp.FileURL)

		start := time.Now()

		out, err := os.Create(directory + "\\update.zip")

		defer out.Close()

		headResp, err := http.Head(updateResp.FileURL)

		if err != nil {

			panic(errors.NewError(err.Error()))

		}

		defer headResp.Body.Close()

		size, err := strconv.Atoi(headResp.Header.Get("Content-Length"))

		if err != nil {

			panic(errors.NewError(err.Error()))

		}

		fmt.Println("Download length: " + strconv.Itoa(size) + " bytes")

		resp, err := http.Get(updateResp.FileURL)

		if err != nil {

			panic(errors.NewError(err.Error()))

		}

		defer resp.Body.Close()

		n, err := io.Copy(out, resp.Body)

		println(n)

		if err != nil {

			panic(errors.NewError(err.Error()))

		}

		elapsed := time.Since(start)
		fmt.Printf("Update download completed in %s", elapsed)

	}

	return nil

}

// CheckForUpdate manually performs an update check and returns a bool based on whether an update is available or not
func (u Updater) CheckForUpdate() (available bool, err error) {

	if !u.isInitialised {
		return false, errors.NewError("updater has not been initialised")
	}

	r, err := http.NewRequest(http.MethodGet, u.VersionURL, nil)

	if err != nil {

		return false, err

	}

	r.Header.Set("User-Agent", "GoUpdate v" + VERSION + " (https://github.com/cxnky/goupdate)")

	res, err := u.webClient.Do(r)

	if err != nil {

		return false, err

	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {

		return false, err

	}

	updateResponse := updateResp
	jsonErr := json.Unmarshal(body, &updateResponse)
	updateResp = updateResponse

	if jsonErr != nil {

		return false, err

	}

	if u.CurrentVersion != updateResponse.Version {

		return true, nil

	} else {

		return false, nil

	}


}