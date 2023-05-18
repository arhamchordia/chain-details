#!/bin/sh

# Load the .env file
if [ -f .env ]; then
    source ./.env
else
    echo "No .env file found. Please create one with the SLACK_API_XOXB_TOKEN variable set."
    exit 1
fi

# Initialize the gcp_project_id variable
gcp_project_id=""

# Parse command-line arguments
if [ "$#" -ge 2 ] && [ "$1" == "--gcp" ]; then
    gcp_project_id="$2"
    shift 2
else
    echo "The --gcp flag is mandatory. Please provide the Google Cloud Project ID."
    exit 1
fi

# Export the GOOGLE_CLOUD_PROJECT_ID environment variable
export GOOGLE_CLOUD_PROJECT_ID="$gcp_project_id"

#Defining reusable command function
run_command_variants() {
    for base_command in "$@"; do
        go run ../main.go $base_command &
        go run ../main.go $base_command --confirmed &
        go run ../main.go $base_command --pending &
    done

    # Check if any of the go run commands failed
    for job in $(jobs -p); do
        wait $job || {
            echo "Error: One or more go run commands failed."
            exit 1
        }
    done
}
# TODO delete this once withdraw supports variants as well
run_command() {
    for base_command in "$@"; do
        go run ../main.go $base_command &
    done

    # Check if any of the go run commands failed
    for job in $(jobs -p); do
        wait $job || {
            echo "Error: One or more go run commands failed."
            exit 1
        }
    done
}

# Run all command variants for bond and unbond in parallel
run_command_variants "bigquery bond" "bigquery unbond"
run_command "bigquery withdraw" # TODO remove this as well for the reason above

# Create a new folder and move generated reports inside that folder
TODAY_DATETIME=$(date +"%Y-%m-%d_%H-%M-%S")
mkdir -p $TODAY_DATETIME
mv ./output/*.csv ./$TODAY_DATETIME
rm -rf output

# Create a zip archive of the reports
zip -r ${TODAY_DATETIME}.zip ${TODAY_DATETIME}

echo ${TODAY_DATETIME}
echo ${SLACK_OAUTH_TOKEN}

# Upload the zip archive to Slack
curl -F file=@${TODAY_DATETIME}.zip \
     -F "initial_comment=Generated reports on ${TODAY_DATETIME}" \
     -F channels='#monitor-vault-reports' \
     -H "Authorization: Bearer ${SLACK_OAUTH_TOKEN}" \
     https://slack.com/api/files.upload

rm -rf ${TODAY_DATETIME} ${TODAY_DATETIME}.zip

# Just a final feedback about operation successful
echo "All the reports have been generated and pushed"