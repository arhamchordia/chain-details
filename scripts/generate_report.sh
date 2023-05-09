#!/bin/sh

# Checkout the reports branch or create if does nt exist
git checkout branch/Reports 2>/dev/null || git checkout -b branch/Reports

run_command_variants() {
    base_command=$1
    go run main.go $base_command
    go run main.go $base_command --confirmed
    go run main.go $base_command --pending
}

# Run all command variants for bond, unbond, and withdraw
run_command_variants "bigquery bond"
run_command_variants "bigquery unbond"
run_command_variants "bigquery withdraw"

# Create a new folder and move generated reports inside that folder
TODAY_DATETIME=$(date +"%Y-%m-%d_%H-%M-%S")
mkdir -p $TODAY_DATETIME/reports
mv *.csv $TODAY_DATETIME/reports

# Git add, commit, and push
git add $TODAY_DATETIME/reports
git commit -m "Generated reports on $TODAY_DATETIME"
git push

# Just a final feedback about operation successful
echo "All the reports have been generated and pushed"
