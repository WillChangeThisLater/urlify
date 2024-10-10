#!/bin/bash

cat <<EOF
I have the following repo structure

$(tree)

The core of my logic is in pkg/urlify

\`\`\`pkg/urlify/urlify.go
$(cat pkg/urlify/urlify.go)
\`\`\`

I would like to import the \`Urlify\`
function into the CLI app living in cmd/cli/main.go

How can I do this?
EOF
