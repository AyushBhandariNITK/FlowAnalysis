#!/bin/bash

# Directory containing the tar files
input_dir="./tar-images"

# Check if the input directory exists
if [[ ! -d "$input_dir" ]]; then
  echo "Error: Directory $input_dir does not exist."
  exit 1
fi

# Change to the tar-images directory
cd "$input_dir" || { echo "Failed to enter directory $input_dir"; exit 1; }

# Check if there are any .tar files in the directory
tar_files=(*.tar)

if [[ ${#tar_files[@]} -gt 0 ]]; then
  # Loop through each tar file and load it
  for tar_file in "${tar_files[@]}"; do
    echo "Loading Docker image from $tar_file..."
    docker load -i "$tar_file"

    if [[ $? -eq 0 ]]; then
      echo "Successfully loaded image from $tar_file."
    else
      echo "Failed to load image from $tar_file." >&2
    fi
  done
else
  echo "No .tar files found in $input_dir."
fi

echo "All images in $input_dir have been processed."
