![Logo](http://i.imgur.com/xnYSYVn.png)
# Chappy

Chappy is the simplest way to deploy websites using GitHub webhooks.

Define a project and deployment script to run and Chappy will listen for changes to your GitHub repo and deploy automatically. It's perfect for small scale deployments like a single DigitalOcean server ❄️

## Quick Start Guide

### Download

| Operating System                                                                                 |
|--------------------------------------------------------------------------------------------------|
| [Linux x64](https://github.com/danbovey/Chappy/raw/master/build/chappy-linux-amd64)              |
| [MaxOS (10.7 and above)](https://github.com/danbovey/Chappy/raw/master/build/chappy-macos-amd64) |
| [Windows x64](https://github.com/danbovey/Chappy/raw/master/build/chappy-windows-amd64.exe)      |
| [Windows x32](https://github.com/danbovey/Chappy/raw/master/build/chappy-windows-386.exe)        |

There are more options for different types of operating system in the `build` folder.

### Installation

To get started, download the correct executable for your operating system, and move it to a directory in your `$PATH`. For example on Linux x64:

```bash
$ mv chappy-linux-amd64 chappy
$ sudo mv chappy /usr/bin/local/chappy
```

❓ Questions? You can read the step by step installation instructions for [Linux](https://github.com/danbovey/Chappy/wiki/Installing-Chappy-on-Linux), [Mac](https://github.com/danbovey/Chappy/wiki/Installing-Chappy-on-Mac) and [Windows](https://github.com/danbovey/Chappy/wiki/Installing-Chappy-on-Windows).

### Creating a project

The next step is to create a projects file which will define the webhooks you want to serve for one or more repos. First, create an empty `projects.json` file in your home or www directory. Let's define a project named`MyWebsite` that will run a deploy script located in `/var/www/MyWebsite/deploy.sh`.

```json
[
  {
    "name": "MyWebsite",
    "branch": "master",
    "script": "/var/www/MyWebsite/deploy.sh",
    "secret": "<SECRET>"
  }
]
```

🔐 To make sure only GitHub can run your webhook, each project should have a unique secret string. You can quickly generate a random 32 character string by running `chappy secret`, or use a random password generator - either way, make sure to replace `<SECRET>`.

### Creating a deploy script

Your deploy script should be an executable script (make sure to `chmod +x deploy.sh`). At it's most basic form, it should run `git pull`, to update the repo with the latest changes. The example below installs any new dependencies with composer and npm and rebuilds assets using gulp.

```bash
#!/bin/bash
git pull
composer install
npm install
gulp --production
```

Arguments with the event details are passed to the script, which can be used to run commands dynamically. There are some more advanced deploy script examples on the [Deploy Script page](https://github.com/danbovey/Chappy/wiki/Deploy-script) that show how to use this feature.

You can now start Chappy using

```bash
$ chappy start
```

⚙ Check the [CLI page](https://github.com/danbovey/Chappy/wiki/CLI) to see a list of commands available, how to configure the IP and port that Chappy runs on and enable other settings like hot reloading the projects file or serving over HTTPS.

### Creating the webhook

Add a new Webhook to your GitHub repo, which can be found in Settings -> Webhooks -> Add webhook.

- By default, the payload URL will be your server IP, port 9000 and then your project name. i.e. `http://123.456.0.1:9000/MyWebsite`.
- A content type of `application/json` is recommended but it can be any.
- The secret should be the secret string you defined in `projects.json`.
- The event you need to listen to is just the `push` event.

### Testing

⏰ Coming Soon - Chappy will log to a file letting you know if anything's wrong - and listen for the GitHub `ping` event to skip having to make a test commit.

To test everything runs successfully, make a test commit or pull request to the main branch (`"Beep, Boop! - Testing Chappy 🤖"` will do just fine).

If the webhook finishes without errors and your script runs correctly, then congrats 🎉! If there are errors, please read the [Troubleshooting page](https://github.com/danbovey/Chappy/wiki/Troubleshooting) or submit an issue. 🕷

### Running forever

To run Chappy in the background, follow the step by step guide for:

- [Linux](https://github.com/danbovey/Chappy/wiki/Installing-Chappy-on-Linux)
- [Mac](https://github.com/danbovey/Chappy/wiki/Installing-Chappy-on-Mac)
- [Windows](https://github.com/danbovey/Chappy/wiki/Installing-Chappy-on-Windows)

## License

- A lot of webhook logic taken from the [webhook](https://github.com/adnanh/webhook) library by [adnanh](https://github.com/adnanh).
- Logo by [Arsenty](https://thenounproject.com/arsenty/) from the Noun Project.

---

Who's a good boy?

![Chappy is.](http://i.imgur.com/jceU3mv.gif)
