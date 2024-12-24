# XenoxAssistant
A Telegram bot that logs conversations and handles contact information.

## Prerequisites
- Docker installed on your system
- Telegram Bot Token

## Building and Running with Docker

### 1. Build the Docker image
```bash
docker build -t xenox-assistant .
```

### 2. Run the container
There are two ways to run the container:

#### Using environment variable directly:
```bash
docker run -d \
  --name xenox-bot \
  -e TELEGRAM_BOT_TOKEN=your_token_here \
  xenox-assistant
```

#### Using .env file:
Create a .env file with:
```env
TELEGRAM_BOT_TOKEN=your_token_here
```

Then run:
```bash
docker run -d \
  --name xenox-bot \
  --env-file .env \
  xenox-assistant
```

### 3. Check container status
```bash
# View logs
docker logs xenox-bot

# Check if container is running
docker ps

# Stop container
docker stop xenox-bot
```

## Bot Usage
1. Start chat with `/start`
2. Bot will request contact information
3. All conversations are logged in 

conversations.json


```