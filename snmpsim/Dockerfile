FROM python:3.11-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    gcc \
    python3-dev \
    snmp \
    && rm -rf /var/lib/apt/lists/*

RUN pip install --no-cache-dir snmpsim

WORKDIR /usr/local/snmpsim

COPY data /usr/local/snmpsim/data

EXPOSE 11161/udp

CMD ["snmpsim-command-responder", "--agent-udpv4-endpoint=0.0.0.0:11161", "--process-user=nobody", "--process-group=nogroup"]
