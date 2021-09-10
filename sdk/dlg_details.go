package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type DlgDetails struct {
	DownloadDetails     []DownloadDetails   `json:"downloadFiles"`
	EligibilityResponse EligibilityResponse `json:"eligibilityResponse"`
	EulaResponse        EulaResponse        `json:"eulaResponse"`
}

type DownloadDetails struct {
	FileName       string `json:"fileName"`
	Sha1Checksum   string `json:"sha1checksum"`
	Sha256Checksum string `json:"sha256checksum"`
	Md5Checksum    string `json:"md5checksum"`
	Build          string `json:"build"`
	ReleaseDate    string `json:"releaseDate"`
	FileType       string `json:"fileType"`
	Description    string `json:"description"`
	FileSize       string `json:"fileSize"`
	Title          string `json:"title"`
	Version        string `json:"version"`
	Status         string `json:"status"`
	UUID           string `json:"uuid"`
	Header         bool   `json:"header"`
	DisplayOrder   int    `json:"displayOrder"`
	Relink         bool   `json:"relink"`
	Rsync          bool   `json:"rsync"`
}

type EligibilityResponse struct {
	EligibleToDownload bool `json:"eligibleToDownload"`
}
type EulaResponse struct {
	EulaAccepted bool   `json:"eulaAccepted"`
	EulaURL      string `json:"eulaURL"`
}

type FoundDownload struct {
	DownloadDetails    DownloadDetails
	EulaAccepted       bool
	EligibleToDownload bool
}

const (
	dlgDetailsURLAuthenticated = baseURL + "/channel/api/v1.0/dlg/details"
	dlgDetailsURLPublic        = baseURL + "/channel/public/api/v1.0/dlg/details"
)

var ErrorDlgDetailsInputs = errors.New("dlgDetails: downloadGroup or productId invalid")
var ErrorMultipleFileGlob = errors.New("dlgDetails: file glob invalid. can only contain a single *")
var ErrorNoFileGlob = errors.New("dlgDetails: fileNameGlob must contain a *")
var ErrorNoMatchingFiles = errors.New("dlgDetails: no files match provided glob")
var ErrorMultipleMatchingFiles = errors.New("dlgDetails: more than 1 file matches glob")
var ErrorEulaUnaccepted = errors.New("dlgDetails: EULA needs to be accepted for this version")
var ErrorNotEntitled = errors.New("dlgDetails: user is not entitled to download this file")

// curl "https://my.vmware.com/channel/public/api/v1.0/dlg/details?downloadGroup=VMTOOLS1130&productId=1073" |jq
func (c *Client) GetDlgDetails(downloadGroup, productId string) (data DlgDetails, err error) {
	err = c.EnsureLoggedIn()
	// Use public URL when user is not logged in
	// This will not return entitlement or EULA sections
	var dlgDetailsURL string
	if err != nil {
		dlgDetailsURL = dlgDetailsURLPublic
	} else {
		dlgDetailsURL = dlgDetailsURLAuthenticated
	}

	search_string := fmt.Sprintf("?downloadGroup=%s&productId=%s", downloadGroup, productId)
	var res *http.Response
	res, err = c.HttpClient.Get(dlgDetailsURL + search_string)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode == 400 {
		err = ErrorDlgDetailsInputs
		return
	} else if res.StatusCode == 401 {
		err = ErrorNotAuthenticated
		return
	}

	err = json.NewDecoder(res.Body).Decode(&data)

	return
}

func (c *Client) FindDlgDetails(downloadGroup, productId, fileNameGlob string) (data FoundDownload, err error) {
	globCount := strings.Count(fileNameGlob, "*")
	if globCount == 0 {
		err = ErrorNoFileGlob
		return
	} else if globCount > 1 {
		err = ErrorMultipleFileGlob
		return
	}

	if err = c.EnsureLoggedIn(); err != nil {
		return
	}

	var dlgDetails DlgDetails
	dlgDetails, err = c.GetDlgDetails(downloadGroup, productId)
	if err != nil {return}

	// Search for file which matches the single glob pattern
	splitString := strings.Split(fileNameGlob, "*")
	foundFiles := 0
	var foundDownload DownloadDetails
	for _, download := range dlgDetails.DownloadDetails {
		fn := download.FileName
		if strings.HasPrefix(fn, splitString[0]) && strings.HasSuffix(fn, splitString[1]) {
			foundFiles++
			foundDownload = download
		}
	}

	if foundFiles == 0 {
		err = ErrorNoMatchingFiles
	} else if foundFiles > 1 {
		err = ErrorMultipleMatchingFiles
	} else {
		data = FoundDownload{
			DownloadDetails:    foundDownload,
			EulaAccepted:       dlgDetails.EulaResponse.EulaAccepted,
			EligibleToDownload: dlgDetails.EligibilityResponse.EligibleToDownload,
		}

	}
	return
}
