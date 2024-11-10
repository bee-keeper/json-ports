FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    sqlite3 \
    libsqlite3-dev \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

RUN addgroup --system appgroup && adduser --system --ingroup appgroup --disabled-password --gecos "" appuser

WORKDIR /app

COPY ./ports /usr/local/bin/ports
COPY ./data/ports.json /app/data/

RUN chmod +x /usr/local/bin/ports

RUN chown -R appuser:appgroup /app
USER appuser

ENTRYPOINT ["/usr/local/bin/ports", "./data/ports.json"]

CMD []
