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
filter_csv() {
    local file="$1"
    local tmp_file="${file}.tmp"

    head -n 1 "$file" > "$tmp_file"

    local now=$(date +%s)

    perl -MTime::Piece -e '
        my $now = shift;
        while (<>) {
            next if $. == 1;
            my @F = split /,/;
            my $timestamp_str = $F[-1];
            $timestamp_str =~ s/ \+0000 UTC//;
            my $timestamp = eval { Time::Piece->strptime($timestamp_str, "%Y-%m-%d %H:%M:%S") };
            if (!$@ && $timestamp->epoch >= $now - 24*60*60) {
                print $_;
            }
        }
    ' "$now" "$file" >> "$tmp_file"

    mv "$tmp_file" "$file"
}


# Run all command variants for bond and unbond in parallel
run_command_variants "bigquery bond" "bigquery unbond"
run_command "bigquery withdraw" # TODO remove this as well for the reason above

# Create a new folder and move generated reports inside that folder
TODAY_DATETIME=$(date +"%Y-%m-%d_%H-%M-%S")
mkdir -p $TODAY_DATETIME
mv ./output/*.csv ./$TODAY_DATETIME
rm -rf output

# Filter each CSV file
for csv_file in ./$TODAY_DATETIME/*.csv; do
    filter_csv "$csv_file"
done

# Initialize counters
total_bonds=0
confirmed_bonds=0
pending_bonds=0
total_unbonds=0
confirmed_unbonds=0
pending_unbonds=0
total_withdraws=0

# Filter each CSV file and update counters
for csv_file in ./$TODAY_DATETIME/*.csv; do
    filter_csv "$csv_file"

    # Update counters based on file name
    case "$csv_file" in
    *_bond_confirmed*.csv)
        confirmed_bonds=$(($(wc -l < "$csv_file") - 1))
        ;;
    *_bond_pending*.csv)
        pending_bonds=$(($(wc -l < "$csv_file") - 1))
        ;;
    *_bond*.csv)
        total_bonds=$(($(wc -l < "$csv_file") - 1))
        ;;
    *_unbond_confirmed*.csv)
        confirmed_unbonds=$(($(wc -l < "$csv_file") - 1))
        ;;
    *_unbond_pending*.csv)
        pending_unbonds=$(($(wc -l < "$csv_file") - 1))
        ;;
    *_unbond*.csv)
        total_unbonds=$(($(wc -l < "$csv_file") - 1))
        ;;
    *_withdraw*.csv)
        total_withdraws=$(($(wc -l < "$csv_file") - 1))
        ;;
    esac
done

# Create a zip archive of the reports
zip -r ${TODAY_DATETIME}.zip ${TODAY_DATETIME}

# Generate the summary message
summary_message="Summary for the last 24 hours:
Bonds:
- Total: $total_bonds
- Confirmed: $confirmed_bonds
- Pending: $pending_bonds
Unbonds:
- Total: $total_unbonds
- Confirmed: $confirmed_unbonds
- Pending (including pending for bonding period, ignore this): $pending_unbonds
Withdraws:
- Total: $total_withdraws"

echo $summary_message
# Upload the zip archive to Slack
curl -F file=@${TODAY_DATETIME}.zip \
     -F "initial_comment=$summary_message" \
     -F channels='#monitor-vault-reports' \
     -H "Authorization: Bearer ${SLACK_OAUTH_TOKEN}" \
     https://slack.com/api/files.upload

rm -rf ${TODAY_DATETIME} ${TODAY_DATETIME}.zip

# Just a final feedback about operation successful
echo "All the reports have been generated and pushed"
