#!/bin/bash

# This is a disaster waiting to happen. Testing in bash 
# is bound to have a crazy amount of edge cases
#
# Even so, this is a good quick and dirty check to make sure
# the lambda functionality is working. If it becomes
# too unstable make it a proper test.

set -euo pipefail

SAM_PID=""
function cleanup {
    echo "exit trap received: cleaning up" >&2
    if [ -n "$SAM_PID" ]; then
        echo "killing sam (pid=$SAM_PID)" >&2
        kill -9 "$SAM_PID"
    fi
}

trap cleanup EXIT

# build the test image. uses whatever image URI is defined in template.yaml
function build_test_image {
        IMAGEURI=$(cat template.yaml | yq -r '.Resources.LambdaFunction.Properties.ImageUri')
        echo "Building image $IMAGEURI" >&2

        # we need to be in the root of the project to build the image
        TESTDIR="$(pwd)"
        cd .. && docker build --platform linux/amd64 -t "$IMAGEURI" . >/dev/null 2>&1
        cd "$TESTDIR"
}

# helper function to determine if a process is alive from its PID
is_process_alive() {
    if kill -0 "$1" 2>/dev/null; then
        return 0  # Process is alive
    else
        return 1  # Process is dead
    fi
}

# using sam kick off server on localhost:3000 which runs the lambda locally
function run_lambda_locally {
        echo "Starting SAM server at localhost:3000" >&2
        sam local start-api >/tmp/urlify-sam-logs 2>&1 &
        SAM_PID=$!
        echo "SAM running as $SAM_PID"
        sleep 5
        if ! is_process_alive "$SAM_PID"; then
                echo "SAM failed to start: $(cat /tmp/urlify-sam-logs)" >&2
                exit 1
        fi
        echo "SAM started successfully"
}

function test_upload_text {
    URL=$(curl -X POST localhost:3000/urlify -F "file=@test_files/test.txt" 2>/dev/null)
    curl -o /tmp/urlify-test-upload-text.txt "$URL" 2>/dev/null
    if ! diff test_files/test.txt /tmp/urlify-test-upload-text.txt >/dev/null; then
        echo "fail: test_upload_text: text contents are different" >&2
    else
        echo "pass: test_upload_text: text contents are the same" >&2
    fi
}

function test_upload_image {
    URL=$(curl -X POST localhost:3000/urlify -F "file=@test_files/rome.jpg" 2>/dev/null)
    curl -o /tmp/urlify-test-rome.jpg "$URL" 2>/dev/null
    if ! diff test_files/rome.jpg /tmp/urlify-test-rome.jpg >/dev/null; then
        echo "fail: test_upload_image: image contents are different" >&2
    else
        echo "pass: test_upload_image: image contents are the same" >&2
    fi
}

echo "Building test image" 2>&1
build_test_image
echo "Running lambda container locally via SAM" 2>&1
run_lambda_locally
echo "Running text upload test" 2>&1
test_upload_text
echo "Running image upload test" 2>&1
test_upload_image
