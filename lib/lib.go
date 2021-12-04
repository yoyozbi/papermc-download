
package lib

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)
type PaperDownloadApi struct {
	BaseURL string
	Project string
	Version string
	latestBuild string `default:""`
	latestFilename string `default:""`
	latestHash string `default:""`
}

func (a PaperDownloadApi) GetJarFile(output string) error{
	var filename string
	var build string
	if a.latestBuild == "" || a.latestFilename == ""{
		err, b, f,_ := a.GetLatestBuild()
		if err != nil {return err}
		build = b
		filename = f
	}else {
		build = a.latestBuild
		filename = a.latestFilename
	}

	if output == "" {output = "."}
	url := a.BaseURL + `/projects/` + a.Project + `/versions/` + a.Version + `/builds/` + build + `/downloads/` + filename

	err, body := get(url)
	if err != nil {return err}

	err = ioutil.WriteFile(filepath.Join(output,"/paper.jar"), body,0644)
	return err
}

// FileIsLatest
/*
  know if the local paper file is the latest by checking the hash
*/
func (a PaperDownloadApi) FileIsLatest(output string) bool {
	err, hash := sha256OfFile(filepath.Join(output,"/paper.jar"))
	if err != nil {return false}

	err, _, _, lHash := a.SaveLatestBuild(false)
	if err != nil {return false}

	return lHash == hash
}
func sha256OfFile(output string) (error, string) {
	f, err := os.Open(output)
	if err != nil {return err,""}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return nil, hex.EncodeToString(h.Sum(nil))
}
func (a PaperDownloadApi) SaveLatestBuild(force bool) (error, string, string, string){
	if a.latestFilename == "" || a.latestBuild == "" || force {
		err, build, filename,hash := a.GetLatestBuild()
		if err != nil {
			return err, "", "", ""
		}
		a.latestBuild = build
		a.latestFilename = filename
		a.latestHash = hash
		return nil, build, filename, hash
	}else {
		return nil,a.latestBuild,a.latestFilename, a.latestHash
	}
}
func (a PaperDownloadApi) GetLatestBuild() (error, string, string,string) {
	url := a.BaseURL + `/projects/` + a.Project + `/version_group/` + a.Version + `/builds`
	err, resp := get(url)
	if err != nil {return err, "", "",""}

	var v map[string]interface{}
	json.Unmarshal(resp, &v)
	builds := v["builds"].([]interface{})
	var numBuilds []float64
	for _, build := range builds {
		num := build.(map[string]interface{})["build"].(float64)
		numBuilds = append(numBuilds, num)
	}
	pos := biggestFloatPos(numBuilds)
	latestFileName := builds[pos].(map[string]interface{})["downloads"].(map[string]interface{})["application"].(map[string]interface{})["name"].(string)
	latestBuild := int(numBuilds[pos])
	latestHash := builds[pos].(map[string]interface{})["downloads"].(map[string]interface{})["application"].(map[string]interface{})["sha256"].(string)
	return nil,strconv.Itoa(latestBuild),latestFileName, latestHash
}
func biggestFloatPos(input []float64) int {
	var biggest float64
	var index int
	for i, num := range input {
		if num > biggest {
			biggest = num
			index = i
		}
	}
	return index
}
func findInArray(array []string, value string) int {
	for i, v := range array {
		if v == value {
			return i
		}
	}
	return -1
}
func get(url string) (error, []byte) {
	resp, err := http.Get(url)
	if err != nil {
		return err, nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, nil
	}
	return nil, body

}
func decodeString(input string) map[string]interface{} {
	var output map[string]interface{}

	json.Unmarshal([]byte(input), &output)
	return output
}