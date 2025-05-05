# Parts-Bin

`parts-bin` is a CLI tool for tracking physical parts — like electronic components, tools, or other items — stored in nested containers (bins). Inspired by a Unix-like command structure, it allows you to manage parts and bins with commands like `ls`, `rm`, and `mv`.

## Features

- REPL and single-command execution modes
- Hierarchical storage system (`/bin1/bin2/part`)
- Unique SKU generation for parts
- Move, rename, or delete parts and bins
- Persistent storage with PostgreSQL (SQLite version coming soon)

## Installation
### Clone Repo
```Bash
git clone https://github.com/itsMe-ThatOneGuy/parts-bin.git
cd parts-bin
```
### Build
```bash
sudo go build -o /usr/local/bin/parts-bin 
```
### Set up PostgreSQL / run Migrations
Ensure PostgreSQL is installed and accessible.
```bash
createdb parts_bin
goose postgres "user=username dbname=parts_bin sslmode=disable" up
```
*Replace `username` with your actual PostgreSQL username.*

## Usage
The first time you run the app, a default config file will be created in `~/.config/parts-bin`.
You will need to update the config file with your PostgreSQL database url.
Alternatively, you can store the DB url in a .env file located anywhere on your machine.
You would then need to update the config file's env_location

### Single-command Execution
```bash
parts-bin <cmd> [options] path
```
A good place to start is the `help` command. This will give you a list of commands

### REPL
```bash
parts-bin
```
You can type `exit` to leave REPL mode.

### Available Commands
| Command | Description                         |
|---------|-------------------------------------|
| help    | Show general help and available commands |
| ls      | List contents of a bin              |
| rm      | Delete part(s) or bin(s)            |
| mv      | Rename a part/bin, or move part/bin to a new location  |
| mkbin   | Create a new bin                    |
| mkprt   | Create a new part                   |
| exit   | Leave REPL mode                   |

Use the flag `-h` with command to see usage details
```bash
parts-bin ls -h
```

## Contributing
If you'd like to contribute, please fork the repo and open a pull request to the `main` branch.
