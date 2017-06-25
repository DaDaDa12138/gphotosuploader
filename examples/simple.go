package main

import (
	"gphotosuploader/api"
	"gphotosuploader/auth"
	"time"
)

// Simple example which consist in the upload of a single image
func main () {
	// Load cookie for credentials from a json file
	credentials, err := auth.NewCookieCredentialsFromFile("../cookies.json")
	if err != nil {
		panic(err)
	}

	// Get a new API token using the TokenScraper from the api package
	token, err := api.NewTokenScraper(credentials).ScrapeNewToken()
	if err != nil {
		panic(err)
	}

	// Add the token to the credentials
	credentials.SetAPIToken(token)


	// Create an UploadOptions object that describes the upload.
	options := api.UploadOptions{
		FileToUpload: "path/to/file.png", // This field is required

		// Below fields are optional
		Name: "logo.png",
		Timestamp: time.Now().Unix(),
	}


	// Create an upload using the NewUpload method from the api package
	upload, err := api.NewUpload(&options, credentials)
	if err != nil {
		panic(err)
	}

	// Finally upload the image
	if err = upload.TryUpload(); err != nil {
		panic(err)
	}

	// Image uploaded!
}