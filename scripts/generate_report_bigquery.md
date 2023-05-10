# generate_report_bigquery.sh

This script generates bond, unbond, and withdraw reports for a Google Cloud Project using the BigQuery API. The script accepts a mandatory flag `--gcp` followed by the Google Cloud Project ID.

## Usage

```bash
./generate_report_bigquery.sh --gcp YOUR_GOOGLE_CLOUD_PROJECT_ID
```

## Features
- The script checks if the --gcp flag is provided and exits with an error message if not.
- Sets the GOOGLE_CLOUD_PROJECT_ID environment variable.
- Runs the main.go command for bond, unbond, and withdraw with different flags (default, --confirmed, --pending).
- Creates a folder named with the current date and time in the /scripts directory and moves the generated CSV files into it.
- Add, commit, and push the generated reports to a remote repository.

## Requirements
- The script should be placed in the /scripts folder.
- The main.go file should be in the parent directory of the /scripts folder.

## Troubleshooting

If you experience issues with the script, please ensure that:
- You have the correct folder structure and the script is placed in the /scripts folder.
- The GOOGLE_CLOUD_PROJECT_ID environment variable is set correctly.
- The necessary dependencies for running the main.go file are installed and working properly.