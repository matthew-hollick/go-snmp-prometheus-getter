FROM apache/superset:4.1.1

USER root

# Install additional dependencies
RUN pip install --no-cache-dir \
    elasticsearch-dbapi[opendistro]

# Create SQLite database directory
RUN mkdir -p /app/superset_home && \
    chown -R superset:superset /app/superset_home

# Copy configuration and initialization script
COPY superset_config.py /app/
COPY init.sh /app/
RUN chmod +x /app/init.sh
ENV SUPERSET_CONFIG_PATH /app/superset_config.py
ENV SUPERSET_HOME /app/superset_home

USER superset
