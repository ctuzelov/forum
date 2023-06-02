FROM golang:alpine AS builder
LABEL authors = "aabdrakh & nazerkegazizova & ctuzelov"
LABEL project-name = "forum"
WORKDIR /app
COPY . .
RUN apk add build-base && go build -o forum cmd/main.go
FROM alpine
WORKDIR /app 
COPY --from=builder /app .
CMD ["./forum"]
EXPOSE 8080