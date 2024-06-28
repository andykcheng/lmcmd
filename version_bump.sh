#!/bin/bash

# Read the current version from the VERSION file
current_version=$(cat VERSION)
IFS='.' read -r -a version_parts <<< "$current_version"

# Increment the patch version
version_parts[2]=$((version_parts[2] + 1))

# Create the new version string
new_version="${version_parts[0]}.${version_parts[1]}.${version_parts[2]}"

# Write the new version back to the VERSION file
echo "$new_version" > VERSION

# Output the new version
echo "New version: $new_version"
