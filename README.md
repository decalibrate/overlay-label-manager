# Overlay Label Manager

Overlay Label Manager is a templating tool for small text files. Its aim is to provide a single place to manage all of the Text Sources that you may want to place on your stream overlay.


The app is essentially a web server which exposes end points to view and update the contents of simple text files. These text files can be as simple as plain text labels or include complex logic via an implementation of [Shopify's Liquid templating language](https://shopify.github.io/liquid/).


Whilst a Text Source is the least memory intensive element to add to an overlay it doesn't provide much customisation and only updates every few seconds (every 60 seconds by defualt on OBS). This app makes each label available through a web url which can be used in a stream Browser Source. The browser view supports custom html layout files which makes the content of these infinitely customisable.


Through a tunneling service like [ngrok](https://ngrok.com) the end points of the web server and UI can be shared with external users and tools. Use caution though as this is directly giving the outside world control of the text on your overlay and any file in the labels directory.

---

## Getting Started

### Run the program

When you run the program without any arguments it will create a folder `.overlaylabelmanager` this folder contains the app configuration and copies of all the variables and label templates that you create.

The app will try to save it in your `/home/UserName/` directory, or if that can't be found it will just create the folder in the directory where the app is executed from.

A second folder is created `labels` in the same directory. This folder will contain the rendered text files you create from the templating and variable system.

Once the initial load process is complete the web UI will be available. The URL will be output in the console window, by default it will be

```http://localhost:9144/```

Visit the UI and check the configuration page, theres a couple of settings you can change there. Alternatively, look in the `.overlaylabelmanager/configuration.json` file and change them directly.

If you're happy with everything there take a look at the Documentation page which has an example of setting up a simple Label template and Variable.


### Additional Runtime Arguments

There are a number of arguments you can pass when you run the app that change the configurations.

| Arg | Example | Description |
| --- | ---     | ----        |
| `-cfg` | `-cfg="/home/user/.configs.json"` | The path for a configuration file, use this parameter if you want your configuration/templates/variables json files to be saved in an alternative location |
| `-p` | `-p=9144` | The port you want the web UI/API to run on, defaults to 9144, this can be changed directly in the configuration.json file or in the Web UI |
| `-dev` | `-dev` | This flag is useful if you're making edits to the UI html/js. It indicates that a local [/static](/static) folder should be used to serve all the static web content. You will need to be running the program within a copy of this repo. The entire static folder is embedded when the binary is built. |

---

## Building the binaries

Prebuilt binaries are avaialble on the [Releases](https://github.com/decalibrate/overlay-label-manager/releases) page


The main program is written in Golang so once you have golang installed building the app is as simple as

```go build ```

This should build the overlay-label-manager.exe file if you're on windows, or the relevant binary if you're on a different OS.


---

This is currently a work in progress. Let me know if you like it though. 👍
