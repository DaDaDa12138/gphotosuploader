# G Photos Uploader - Beta
[![Build Status](https://travis-ci.org/simonedegiacomi/gphotosuploader.svg?branch=master)](https://travis-ci.org/simonedegiacomi/gphotosuploader)

## Why? What is this?
Google Photos does not have a desktop uploader for Linux, neither an API to upload a photo programmatically ... So here
they are!

G Photos Uploader let's you upload photos from Linux (and, in theory, any OS for which you can compile a Go program)
specifying the file name or watching a directory for changes.
Furthermore, the project can also be used as a library that you can include in other Go programs.

## Disclaimer
G Photos Uploader is an unofficial tool, I (and any possible contributor) don't guarantee any result. Any security or
other kind of issues are at your own risk.

## How can i use it?
### Standalone tool
To use G Photos Uploader as a standalone tool you need to get be authenticated. Authentication in handled with a
JSON file that contains your cookies and your user Id.

#### Authentication
Every time your run the tool the program will check for the auth file. If the file is not found or the cookies seems to
be expired the tool will ask you if you want to run a wizard to get new cookies.
The authentication wizard uses the WebDrivers protocol, which is usually used to perform automation tests, that allows
G Photos Upload to control a browser and read the cookies from it. To use the WebDrivers Protocol you need to install a
web driver:

To install the Chromium Chrome Web Driver on Ubuntu use:
```sh
sudo apt-get install chromium-chromedriver

# Create a link to launch the driver just typing 'chromedriver'
sudo ln -s /usr/lib/chromium-browser/chromedriver /usr/bin/chromedriver
```

While on Mac Os X (using [Homebrew](https://brew.sh/)):
```sh
brew install chromedriver
```

And then launch it using:
```sh
chromedriver
```

When the Driver starts it will print the address at which it is listening.
Once you enter the name of the browser and the address of the web driver on the wizard a new browser window will appear
with the Google Photos Login page. Then you can login with your account just like you always do. When you're logged in
the tool will read the cookies from the browser, save them into the auth file and close the browser window.
(You can now stop the web driver server)


#### Upload a photo or watch a directory
Once you have the auth file you're ready to go, just launch the program. For example, to upload a file named image.png:
```sh
go run main.go --upload ./image.png
```

Or to watch a directory:
```sh
go run main.go --watch path/to/photos --maxConcurrent 4
```

You can even upload all the photos of a directory and then start to watch another one:
```sh
go run main.go --upload /path/to/old/photos --upload /downloads/cat.png --watch path/to/new/photos
```


The tool crates a file (which the default name is uploaded.txt) which is a list of uploaded files, which will not be
re-uploaded. You can specify your own file using the uploadedList argument.
To see all the available arguments, use --help.

### Library
You can read a simple example [here](examples/simple.go) or get the documentation [here](http://godoc.org/github.com/simonedegiacomi/gphotosuploader).

## Used libreries
* [fsnotify](https://github.com/fsnotify/fsnotify): To watch for file system events;
* [Selenium](https://github.com/tebeka/selenium): To authenticate using a browser;


## Contributors:
* [simonedegiacomi](https://github.com/simonedegiacomi)
* [alessiofaieta](https://github.com/alessiofaieta)