# Slackbot for LLMs in Go

Uses `botman` to switch between various LLMs.

Create the following environment variables, either with a `.env` file or for real.

```conf
SLACK_BOT_TOKEN=your-bot-token-here
SLACK_APP_TOKEN=your-app-token-here

# can be: claude, openai, fireworksai
BOTMAN_LLM=claude

# Base System Prompt. Note that LLM specific prompts get added onto the base prompt.
# BOTMAN_PROMPT=You are a helpful chat model

BOTMAN_CLAUDE_API_KEY=your-claude-api-key
BOTMAN_CLAUDE_MODEL=claude-3-5-sonnet-20240620
#BOTMAN_CLAUDE_PROMPT=

#Open AI specific settings (required when LLM is set to openai)
#BOTMAN_OPENAI_API_KEY=you-openai-key-here
#BOTMAN_OPENAI_MODEL=gpt-4o
#BOTMAN_OPENAI_PROMPT=

#Fireworks AI specific settings (required when LLM is set to fireworks)
#BOTMAN_FIREWORKS_API_KEY=
#BOTMAN_FIREWORKS_MODEL=
#BOTMAN_FIREWORKS_PROMPT=
```

## Slackbot Manifest

Create your own app [here](https://api.slack.com/apps).

Add the manifest below, and then get a bot token and app-level token.

```yaml
display_information:
  name: Yappidy Yapyap
  description: An LLM Chatbot
  background_color: "#9c2005"
features:
  bot_user:
    display_name: Yappidy Yapyap
    always_online: true
oauth_config:
  scopes:
    bot:
      - app_mentions:read
      - channels:history
      - chat:write
      - commands
      - groups:history
      - im:history
      - im:read
      - im:write
      - reactions:write
settings:
  event_subscriptions:
    bot_events:
      - app_mention
      - message.channels
      - message.groups
      - message.im
  interactivity:
    is_enabled: true
  org_deploy_enabled: false
  socket_mode_enabled: true
  token_rotation_enabled: false
```
