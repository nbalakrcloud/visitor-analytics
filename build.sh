#!/bin/bash

# Clean any existing binaries and modules cache
go clean -modcache

# Build the Go application
go build -o visitor-analytics ./cmd

# Check if the build was successful
if [ $? -eq 0 ]; then
    echo "Build successful. Executable: ./myapp"
else
    echo "Build failed"
fi