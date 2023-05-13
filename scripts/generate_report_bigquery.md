# generate_report_bigquery.sh

This script generates bond, unbond, and withdraw reports for a Google Cloud Project using the BigQuery API. The script
accepts a mandatory flag --gcp followed by the Google Cloud Project ID. Once the report is generated, it gets zipped and
uploaded to a Slack channel.

## Usage

```bash
./generate_report_bigquery.sh --gcp YOUR_GOOGLE_CLOUD_PROJECT_ID
```

## Features

- The script checks if the --gcp flag is provided and exits with an error message if not.
- Sets the GOOGLE_CLOUD_PROJECT_ID environment variable.
- Runs the main.go command for bond, unbond, and withdraw with different flags (default, --confirmed, --pending).
- Creates a folder named with the current date and time in the /scripts directory and moves the generated CSV files into
  it.
- Compresses the report folder into a .zip file.
- Reads the Slack OAuth token from a .env file.
- Uploads the .zip file to a specific Slack channel.

## Requirements

- The script should be placed in the /scripts folder.
- The main.go file should be in the parent directory of the /scripts folder.
- A .env file containing SLACK_OAUTH_TOKEN should be in the same directory as the script. This file should not be
  committed to your repository.

## Troubleshooting

If you experience issues with the script, please ensure that:

- You have the correct folder structure and the script is placed in the /scripts folder.
- The GOOGLE_CLOUD_PROJECT_ID environment variable is set correctly.
- The necessary dependencies for running the main.go file are installed and working properly.
- The .env file exists and has the correct Slack OAuth token.
- The Slack OAuth token has the correct scopes and the bot is invited to the channel where the report will be uploaded.

## Slack App and Bot

This script uses a Slack App and Bot to upload the reports to a Slack channel. You need to create a Slack App, add a Bot
to it, and install it in your workspace. The Bot should be added to the channel where the reports will be uploaded. The
OAuth token should have the `files:write` scopes and being invited to the channel via `/invite @YourBotName`.

## Security Note

Remember to keep your Slack OAuth token secure. Do not commit the .env file containing the token to your repository. You
may want to add .env to your .gitignore file to prevent it from being accidentally committed.