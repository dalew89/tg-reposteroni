# tg-reposteroni

## A (very) simple repost checker for Telegram group chats

This bot runs in the background, and saves URLs to a local SQLite database. If a URL already exists in the database the bot will flag the message and alert the reposter.
A list of reposters is also stored and can be triggered by using the command `/repoststats` in the chat where the bot is active. Reposters will be shown in descending order.

To build:

`go build -o tg-reposteroni`

Required env variables:

`BOT_TOKEN` Your bot token. Learn how to create a new token here: https://core.telegram.org/bots#6-botfather

`DATABASE_PATH` is the location in which the bot's URL and repost database is stored, e.g:

`DATABASE_PATH=data/logs.db`

### Docker

A Docker Image of this bot can be found under the name dockerd89/tg-reposteroni

To run in Docker:

`docker run -d -v /your/local/dir/:/app/data -P -e "DATABASE_PATH=data/reposts.db" -e "BOT_TOKEN=<YOUR_BOT_TOKEN>" dockerd89/tg-reposteroni`