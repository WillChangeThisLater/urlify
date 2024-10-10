## How do I use this?
```bash
./build.sh
sudo ln -s $(pwd)/urlify /usr/local/bin/urlify
urlify build.sh
```

## Why would I use this?
I mainly use this with [my LLM CLI tool](https://github.com/WillChangeThisLater/go-llm). I use it like this:

```bash
echo "When was this building built, and by who?" | lm --imageURLs "$(urlify rome.jpg)"
```

## Lambda
In it's current form, `urlify` is extremely limited: the tool requires AWS creds to upload files to a bucket.
I am working on improving this via a lambda function. The idea is as follows:

  1) Build container for the code and upload to ECR via `./build-lambda.sh`
  2) Create lambda which uses the containerized code to process requests
  3) Throw that lambda behind API Gateway. Route traffic from a domain I control to API Gateway

Infrastructure wise this requires: ECR, Lambda, API Gateway, Route53, S3. Though the use case is quite simplistic
Ultimately I want terraform code to IaC this
