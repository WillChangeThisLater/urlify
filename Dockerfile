FROM golang:latest as build

WORKDIR /lambda

# Copy dependencies list
COPY go.mod go.sum ./

# Build with lambda.norpc tag
# afaik this reduces the image size
COPY cmd/lambda/main.go .
RUN go build -tags lambda.norpc -o main main.go

FROM public.ecr.aws/lambda/provided:al2023
COPY --from=build /lambda/main ./main
ENTRYPOINT [ "./main" ]
