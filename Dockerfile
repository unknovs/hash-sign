FROM scratch

WORKDIR /
EXPOSE 8080/tcp
COPY publish/ ./

ENV TZ=Europe/Riga

ENTRYPOINT ["/server", "main"]
