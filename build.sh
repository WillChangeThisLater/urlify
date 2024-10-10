#!/bin/bash
go build -o urlify cmd/cli/main.go

cat <<EOF
urlify is built! now you'll need to link it into your system path
you should be able to do that with something like

\`\`\`bash
ln -s $(pwd)/urlify /usr/local/bin/urlify # if you hit a permission error, you may need to prefix this with 'sudo'
\`\`\`
EOF
