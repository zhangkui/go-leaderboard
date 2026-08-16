# 官方 Go 镜像，自带完整工具链
FROM golang:1.22

WORKDIR /app

COPY . .

RUN go build ./...

CMD ["bash"]
