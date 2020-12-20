    ▄▄▄▄·             • ▌ ▄ ·.  ▄▄▄▄      ▄▄▄▄▄▄▄
    ▐█ ▀█▪ ▄█▀▄  ▄█▀▄ ·██ ▐███▪▐█ ▀█▪ ▄█▀▄ •██
    ▐█▀▀█▄▐█▌.▐▌▐█▌.▐▌▐█ ▌▐▌▐█·▐█▀▀█▄▐█▌.▐▌ ▐█.▪
    ██▄▪▐█▐█▌.▐▌▐█▌.▐▌██ ██▌▐█▌██▄▪▐█▐█▌.▐▌ ▐█▌·
    ·▀▀▀▀  ▀█▄▀▪ ▀█▄▀▪▀▀  █▪▀▀▀·▀▀▀▀  ▀█▄▀▪ ▀▀▀

## **_The ultimate interactive and intuitive music streaming Discord bot_**

## **Master**-![Tests](https://github.com/aplomBomb/boombot/workflows/Tests/badge.svg)

## **Dev**-![Tests|Dev](https://github.com/aplomBomb/boombot/workflows/Tests/badge.svg?branch=dev)

<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-4%25-brightgreen.svg?longCache=true&style=flat)</a>
[![Go Report Card](https://goreportcard.com/badge/github.com/aplombomb/boombot)](https://goreportcard.com/report/github.com/aplombomb/boombot)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Description

Boombot is a music streaming bot for group listening and sharing in Discord voice channels. The idea behind Boombot is to make music streaming simple, accessible and interactive enough for every Discord user, not just the hip ones. This is Boombot currently supports Youtube videos/playlist links and custom searching from any text channel.

## Getting Started

### Boombot is not currently setup for easy plug and play **_quite_** yet, but if you're eager...

### Dependencies

- Docker
- FFMPEG
- Clone the repo and add these .env files into root dir with their respective key=values
  - `boombot.env`
    - `AWS_ACCESS_KEY=yourawsaccesskeyhere`
    - `AWS_SECRET_ACCESS_KEY=yourawssecretaccesskeyhere`
    - `ENV=container` - this is the literal value that needs to be here
    - `BOT_PREFIX=yourprefferedcommandprefix` - the symbol you'd like to be used for invoking commands
  - `db.env`
    - `POSTGRES_PASSWORD=yourpgpasswordhere`
    - `POSTGRES_USER=yourpguserhere`
- In order to be able to fetch the bot's client and youtube tokens, you'll need to register your own applications and store the tokens under the secret name: `boombot_creds` with the keys `BOT_TOKEN` and `YOUTUBE_TOKEN`

### Building\Launching

- After installing your dependencies and setting up your secrets, run `docker-compose up` at root dir

# How to use Boombot via commands

- Commands
  - ### [_prefix_]**play** + _youtubeLink_ **OR** _a string search query_
    This will invoke Boombot to fetch your request from youtube and put it into your own queue. Boombot will then join you in the voice channel you're currently residing in and start playing.
  - ### [_prefix_]**next**
    This will skip the currently playing song, moving on to your next item in your queue, unless multiple users have queues, in that case Boombot will alternate between the different queues to avoid a listening monopoly for a better group listening experience.
  - ### [_prefix_]**purge**
    This will stop all playback for your queue and remove it completely. This also hapens if you leave voice chat entirely, preventing users from trolling by requesting large playlists, then disappearing.
  - ### [_prefix_]**shuffle**
    This will shuffle your queue, randomizing playback order.
  - ### [_prefix_]**pause**
    Pauses playback
  - ### [_prefix_]**play**
    Resumes paused playback

# How to use Boombot via jukebox

### Boombot uses a special text channel to display the currently playing song info inside an embedded message, these messages have auto-generated reaction emojis attached to them that represent all the different command controls, all you have to do is click them and Boombot will respond accordingly!
