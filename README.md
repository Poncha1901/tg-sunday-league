# TG Sunday League

## Table of Contents

- [Introduction](#introduction)
- [Project Structure](#project-structure)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Introduction

The TG Sunday League project is a Telegram bot designed to help manage football game arrangements. It allows users to create games, add players, track attendees, and manage other game details, all within a Telegram chat interface.

## Project Structure

The project is organized into the following directories:

- `bot/`: Contains the bot logic, including handling Telegram commands, user interactions, and message delivery.
- `config/`: Loads environment variables, such as API tokens and database configurations, using `.env` files.
- `db/`: Contains the SQLite database initialization and connection handling for persistent data storage.
- `models/`: Defines the data models used in the bot, representing entities like games, users, and player statuses.
- `repositories/`: Contains repository interfaces and implementations for accessing and managing database entries.
- `services/`: Encapsulates the main logic of the bot, including handling game creation, player management, and other core functionalities.

## Installation

To set up the TG Sunday League bot, follow these steps:

1. **Clone the repository**:
    ```sh
    git clone https://github.com/yourusername/tg-sunday-league.git
    ```

2. **Navigate to the project directory**:
    ```sh
    cd tg-sunday-league
    ```

3. **Set up the environment**:
   - Create a `.env` file in the project root directory and add the necessary configuration values. For example:
     ```plaintext
     TELEGRAM_BOT_TOKEN=your_bot_token_here
     DATABASE_URL=sqlite://db/tg_sunday_league.db
     ```
   - Replace `your_bot_token_here` with your actual Telegram bot token.

4. **Run the bot**:
    ```sh
    go run main.go
    ```

## Usage

After installation, start the bot using the `/new` command in Telegram to create a new game. Use the following format:

```plaintext
/new (YYYY-MM-DD, HH:MM, Location, Opponent, Optional[Price])
