#!/bin/sh

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

# Checkout the reports branch or create if does not exist TODO uncomment below once agreed with the team
#git checkout branch/Reports 2>/dev/null || git checkout -b branch/Reports

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

# Git add, commit, and push TODO uncomment below once agreed with the team
#git add $TODAY_DATETIME
#git commit -m "Generated reports on $TODAY_DATETIME"
#git push

# Just a final feedback about operation successful
echo "All the reports have been generated and pushed"
