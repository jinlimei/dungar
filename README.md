# Dungar

Dungar is a Go-based Markov chat bot with Discord support, Postgres-backed message storage, and a growing set of trigger-driven behaviors.

The current runtime is Discord-first. Slack code still exists, but the actively maintained path is Discord.

## Features

- Markov-based responses trained from stored Discord messages
- Discord history backfill for the last 30 days in configured channels
- Channel allowlist for both speaking and learning
- Manic minute trigger system with persistent cooldowns and admin controls
- User simulation: `@Dungarmatic simulate @user`
- User blocklist and deletion of stored messages by author

## Requirements

- Go `1.21`
- Docker and Docker Compose for the local runtime stack
- A Discord bot token
- A PostgreSQL database

## Configuration

Start from the example files:

```sh
cp settings.ini.example settings.ini
cp secrets.ini.example secrets.ini
```

Important settings:

```ini
[base]
mode=discord

[discord]
allowed_output_channel_ids=362144221148479489,471817271749115905
```

```ini
[discord]
token=YOUR_BOT_TOKEN

[pgsql]
host=postgres
user=postgres
pass=postgres
data=dungar
```

`allowed_output_channel_ids` restricts both Discord output and Markov learning. Discord learning also only considers the last 30 days of recorded rows.

## Discord Setup

In the Discord Developer Portal, enable:

- `MESSAGE CONTENT INTENT`
- `SERVER MEMBERS INTENT`

The runtime only needs the bot token. The Application ID is used for invite URLs, not for `dungar run`.

## Build and Run

Build locally:

```sh
make build
./bin/dungar help
```

Run the local stack with Docker:

```sh
docker compose up --build
```

That starts:

- `postgres` on an internal Compose network
- `dungar` with `settings.ini` and `secrets.ini` mounted read-only

Fresh Docker volumes load [infra/db/structure.sql](/home/mcp/projects/dungar/infra/db/structure.sql) automatically.

## Existing Database Volumes

If you already have a Postgres volume, apply SQL migrations in [infra/db/migrations](/home/mcp/projects/dungar/infra/db/migrations) manually, in filename order. Example:

```sh
docker compose exec -T postgres psql -U postgres -d dungar < infra/db/migrations/20260609_add_manic_minute_runtime_state_cooldown.sql
```

## Useful Commands

- `dungar run` starts the bot
- `dungar learn` exercises DB-backed learning in a standalone process
- `dungar backfill-discord` imports the last 30 days from the configured Discord channels
- `dungar manic-word` prints the current manic word, current chance, and cooldown remaining
- `dungar manic-word rotate` picks a different random persisted manic word
- `dungar manic-word set <word>` forces a specific persisted manic word
- `dungar bot-info` shows protocol-specific bot metadata

If `dungar run` is already running, restart it after `manic-word rotate` or `manic-word set` so the live bot reloads the new trigger word.

## Discord Bot Commands

- `!manic status`
- `!manic stats`
- `!manic rotate`
- `!manic stop`
- `!manic test [word]`
- `!blocklist add @nick`

Examples of directed triggers:

- `@Dungarmatic tacos or burgers?`
- `@Dungarmatic how fried is xyz?`
- `@Dungarmatic simulate @Kevin`

## Testing

Run the default suite:

```sh
make test
```

Targeted builds and tests:

```sh
make build
go run ./cmd/dungar help
```
