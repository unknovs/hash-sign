FROM scratch

WORKDIR /
EXPOSE 80/tcp
COPY publish/ ./

ENV TZ=Europe/Riga

ENTRYPOINT ["/server", "main"]
