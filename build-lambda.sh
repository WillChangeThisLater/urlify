#!/bin/bash

set -euo pipefail

# deployment script. this is keyed to my personal AWS creds and ECR repo

REGION="us-east-2"
LOCAL_TAG="urlify:latest"
REPO=${LOCAL_TAG%%:*}
VERSION=${LOCAL_TAG#*:}

# function for extracting AWS account number from
# whatever system you're running this on
function getAWSAccountNum() {
        aws sts get-caller-identity | jq -r -e '.Account'
}

function getImageDigest() {
    aws ecr list-images --repository-name "$REPO" --region "$REGION" | jq -r ".imageIds.[] | select(.imageTag == \"$VERSION\") | .imageDigest"
}

# function for extracting the ECR URI for a repository
# from its name
#
# for this script we'll only be using 'urlify', but might
# as well make it generic
function getECRUri() {

        # set -u will actually fail if no argument is passed
        # but double check here in any case
        if [ -z "$1" ]; then
                echo "Usage: getEcrUri <repoName>" >&2
                exit 1
        fi
        aws ecr describe-repositories | jq -r -e ".repositories.[] | select(.repositoryName == \"$1\") | .repositoryUri"
}

# build the image
echo "Building docker image" >&2
docker build --platform linux/amd64 -t "$LOCAL_TAG" . >/dev/null 2>&1

# figure out the account number
echo "Getting AWS Account number" >&2
ACCOUNTNUM="$(getAWSAccountNum)"

# make sure the account number is set
#
# because "set -euxo pipefail" is set at the top of the script
# and 'jq -e' is set in the getAWSAccountNum() function, the
# script SHOULD stop if either the 'aws sts get-caller-identity'
# command returns a non-zero exit code or for some reason .Account
# doesn't show up in the JSON
#
# btw, if you forget what "set -euxo pipefail" means,
# go remind yourself: http://redsymbol.net/articles/unofficial-bash-strict-mode/
# it's important.
#
# anyway, the behavior described above should guarantee that, if the function
# above succeeds, ACCOUNTNUM is set to something. but i am paranoid and
# double check that here
if [ -z ${ACCOUNTNUM} ]; then
        echo "Account number not found"
fi

# I assume that the ECR repo exists already
echo "Getting ECR URI" >&2
ECR_URI="$(getECRUri $REPO)"

echo "Logging in to ECR" >&2
aws ecr get-login-password --region "$REGION" | docker login --username AWS --password-stdin "$ACCOUNTNUM".dkr.ecr."$REGION".amazonaws.com >/dev/null 2>&1

REMOTE_TAG="$ECR_URI:$VERSION"
echo "Tagging image as $REMOTE_TAG" >&2
docker tag "$LOCAL_TAG" "$REMOTE_TAG" >/dev/null 2>&1

echo "Pushing image $REMOTE_TAG" >&2
docker push "$REMOTE_TAG" >/dev/null 2>&1

imageDigest=$(getImageDigest)
echo "Image digest: $imageDigest"
