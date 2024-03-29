FROM golang:latest AS builder

COPY . /github.com/Mobo140/projects/tg-bot-notes_links/
WORKDIR /github.com/Mobo140/projects/tg-bot-notes_links/

RUN go mod download
RUN go build -o /bin/ cmd/bot/main.go


FROM golang:latest

WORKDIR /root/

COPY --from=0 /github.com/Mobo140/projects/tg-bot-notes_links/bin/bot .
COPY --from=0 /github.com/Mobo140/projects/tg-bot-notes_links/configs configs/

EXPOSE 80 

CMD ["./bot"]
