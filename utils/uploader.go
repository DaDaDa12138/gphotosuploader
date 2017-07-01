package utils

import (
	"gopkg.in/headzoo/surf.v1/errors"
	"github.com/simonedegiacomi/gphotosuploader/auth"
	"github.com/simonedegiacomi/gphotosuploader/api"
	"sync"
	"os"
)

// Simple client used to implement the tool that can upload multiple photos or videos at once
type ConcurrentUploader struct {
	credentials       auth.Credentials

	// Buffered channel to limit concurrent uploads
	concurrentLimiter chan bool

	// Map of uploaded files (used as a set)
	uploadedFiles     map[string]bool

	// Waiting group used for the implementation of the Wait method
	waitingGroup      sync.WaitGroup

	// Flag to indicate if the client is waiting for all the upload to finish
	waiting           bool


	// Channel that is used to communicate CompletedUploads
	CompletedUploads  chan string

	// Channel that is used to communicate IgnoredUploads (ex: a file is not an image/video)
	IgnoredUploads    chan string

	// Channel that is used to communicate errors
	Errors            chan error
}

// Creates a new ConcurrentUploader using the specified credentials. The second argument is the maximum number
// of concurrent uploads (which must not be 0).
func NewUploader (credentials auth.Credentials, maxConcurrentUploads int) (*ConcurrentUploader, error) {
	if maxConcurrentUploads <= 0 {
		return nil, errors.New("maxConcurrentUploads must be greather than zero")
	}

	return &ConcurrentUploader{
		credentials: credentials,

		concurrentLimiter: make(chan bool, maxConcurrentUploads),

		uploadedFiles: make(map[string]bool),

		CompletedUploads: make(chan string),
		IgnoredUploads: make(chan string),
		Errors: make(chan error),
	}, nil
}

// Add files to the list of already uploaded files
func (u *ConcurrentUploader) AddUploadedFiles (files ...string) {
	for _, name := range files {
		u.uploadedFiles[name] = true
	}
}

// Enqueue a new upload. You must not call this method while waiting for some uploads to finish (The method return an
// error if you try to do it).
// Due to the fact that this method is asynchronous, if nil is return, it doesn't mean the the upload was completed,
// for that check the Errors and CompletedUploads channels
func (u *ConcurrentUploader) EnqueueUpload(filePath string) error {
	if u.waiting {
		return errors.New("Can't add new uploads when waiting")
	}
	if _, uploaded := u.uploadedFiles[filePath]; uploaded {
		u.IgnoredUploads <- filePath
		return nil
	}

	// Check if the file is an image or a video
	if valid, err := IsImageOrVideo(filePath); err != nil {
		u.Errors <- err
		return nil
	} else if !valid {
		u.IgnoredUploads <- filePath
		return nil
	}

	u.waitingGroup.Add(1)
	go u.uploadFile(filePath)

	return nil
}

func (u *ConcurrentUploader) decrementLimit () {
	<- u.concurrentLimiter
}

func (u *ConcurrentUploader) uploadFile(filePath string) {
	u.concurrentLimiter <- true
	defer u.decrementLimit()

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		u.Errors <- err
	}
	defer file.Close()

	// Create options
	options, err := api.NewUploadOptionsFromFile(file)
	if err != nil {
		u.Errors <- err
	}

	// Create a new upload
	upload, err := api.NewUpload(options, u.credentials)
	if err != nil {
		panic(err)
	}

	// Try to upload the image
	if err := upload.TryUpload(); err != nil {
		u.Errors <- err
	} else {
		u.uploadedFiles[filePath] = true
		u.CompletedUploads <- filePath
	}

	u.waitingGroup.Done()
}

// Blocks the goroutine until all the upload are completed. You can not add uploads when a goroutine call this method
func (u *ConcurrentUploader) WaitUploadsCompleted () {
	u.waiting = true
	u.waitingGroup.Wait()
	u.waiting = false
}