package main

import (
	"flag"
	"log"
	"paper/lib"
)

func main() {
	baseURL := "https://api.papermc.io/v2"
	defProject := "paper"
	defVersion := "1.19"

	var output string
	flag.StringVar(&output, "o", ".", "path of output file")

	var version string
	flag.StringVar(&version, "v", defVersion, "version of minecraft")

	flag.Parse()
	api := lib.PaperDownloadApi{BaseURL: baseURL, Project: defProject, Version: version}

	if !api.FileIsLatest(output) {

		err, build, _, _ := api.SaveLatestBuild(false)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("downloading build " + build + " from " + version)

		err = api.GetJarFile(output)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("downloaded paper.jar inside " + output)

	} else {
		log.Println("latest version is installed")
	}
}
