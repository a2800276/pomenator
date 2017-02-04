package pomenator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type UploadBundleResponse struct {
	RepositoryUris []string `json:"repositoryUris"`
}

func (u *UploadBundleResponse) StagedRepositoryID() string {
	// is probably https://oss.sonatype.org/content/repositories/dekuriositaet-1029
	// need the bit after the last '/'
	if len(u.RepositoryUris) < 1 {
		panic("sonatype is returning madness...")
	}
	ru := u.RepositoryUris[0]
	i := strings.LastIndex(ru, "/")
	if i == -1 {
		panic("unexpected response: " + ru)
	}
	return ru[i+1:]
}

const sonatype = "https://oss.sonatype.org/service/local/staging/bundle_upload"

// -> file
// <- {"repositoryUris":["https://oss.sonatype.org/content/repositories/dekuriositaet-1029"]}
const drop = "https://oss.sonatype.org/service/local/staging/bulk/drop"

// -> {data: {description: "", stagedRepositoryIds: ["dekuriositaet-1029"]}}
// <- 201 : empty

const release = "https://oss.sonatype.org/service/local/staging/bulk/promote"
const releaseTemplate = `{"data":{"autoDropAfterRelease":true,"description":"bumsi!","stagedRepositoryIds":["%s"]}}`

// some time passes between the upload and the time the sonatype mainframe
// runs the batch job to set the status from open to closed (whatever that means)
// since their API is undocumented, we just need to keep trying until we on longer
// an error. This is is the number of times to try.
const releaseMaxTries = 3
const releaseWaitSeconds = 10

// -> {"data":{"autoDropAfterRelease":true,"description":"bumsi!","stagedRepositoryIds":["dekuriositaet-1030"]}}
// <- 201 : empty

const uid = "a2800276"
const pwd = "PEZCUWXu%Kb5"

func UploadBundle(fn string) (repoid string, err error) {
	fmt.Printf("uploading: %s to %s\n", fn, sonatype)
	file, err := os.Open(fn)
	if err != nil {
		return
	}
	defer file.Close()

	req, err := http.NewRequest(http.MethodPost, sonatype, file)
	if err != nil {
		return
	}
	req.SetBasicAuth(uid, pwd)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		io.Copy(os.Stdout, resp.Body)
		return "", fmt.Errorf("unexpected response: %v", resp)
	}

	decoder := json.NewDecoder(resp.Body)
	var ubr UploadBundleResponse
	if err = decoder.Decode(&ubr); err != nil {
		return
	}
	return ubr.StagedRepositoryID(), nil
}

func ReleaseRepo(repo string) error {
	return releaseRepo(repo, 1)
}

func releaseRepo(repo string, count int) error {
	if count > releaseMaxTries {
		return fmt.Errorf("giving up releasing %s, please go to oss.sonatype.org to release manually", repo)
	}

	waitingPeriod := releaseWaitSeconds * count
	fmt.Printf("attempting to release %s after %ds. Attempt (%d/%d)\n", repo, waitingPeriod, count, releaseMaxTries)
	time.Sleep(time.Second * time.Duration(waitingPeriod))

	releaseJson := fmt.Sprintf(releaseTemplate, repo)

	req, err := http.NewRequest(http.MethodPost, release, strings.NewReader(releaseJson))
	if err != nil {
		return err
	}
	req.SetBasicAuth(uid, pwd)
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 201:
		return nil

	case 500: // check
		// 500 {"errors":[{"id":"*","msg":"Unhandled: Repository: dekuriositaet-1034 has invalid state: open"}]}
		// if we come across the above, just keep looping I guess. Sheeeessh!
		// then ...
		// 500 {"errors":[{"id":"*","msg":"Unhandled: Staging repository is already transitioning: dekuriositaet-1041"}]}
		if body, err := bodyAsString(resp); err != nil {
			return err
		} else if !strings.Contains(body, "has invalid state: open") && !strings.Contains(body, "Staging repository is already transitioning") {
			return fmt.Errorf("unexpected response: %d %v", resp.StatusCode, body)
		} else {
			fmt.Printf("got: %s, retrying\n", body)
		}
		return releaseRepo(repo, count+1)

	default:
		io.Copy(os.Stdout, resp.Body)
		return fmt.Errorf("unexpected response: %d %v", resp.StatusCode, resp)
	}
}

func bodyAsString(resp *http.Response) (body string, err error) {
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		return
	}
	return buf.String(), err
}
