# FACEIT Discord Bot

This project is a Discord bot that integrates with the FACEIT API to provide player statistics and other features. 

## Features
- Fetch player statistics from the FACEIT API.
- Store and manage player data in a SQLite database.
- Interact with users via Discord commands.

## Prerequisites
To run this project locally, you will need:

1. **Go**: Install Go from [https://golang.org/dl/](https://golang.org/dl/).
2. **SQLite**: Ensure SQLite is installed on your system.
3. **Discord Bot Token**: Create a bot on the [Discord Developer Portal](https://discord.com/developers/applications) and copy the bot token.
4. **FACEIT API Key**: Obtain an API key from the [FACEIT Developer Portal](https://developers.faceit.com/).

## Setup Instructions

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd faceit-stats-bot
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Set up environment variables:
   Create a `.env` file in the root directory and add the following:
   ```env
   DISCORD_BOT_TOKEN=<your-discord-bot-token>
   FACEIT_API_KEY=<your-faceit-api-key>
   ```

4. Run database migrations:
   ```bash
   make migrate
   ```

5. Start the bot:
   ```bash
   make run
   ```

## License
This project is licensed under the MIT License.
