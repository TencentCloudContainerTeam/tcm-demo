FROM mongo
RUN mkdir -p /app/data/
COPY users_data.json /app/data/
COPY script.sh /docker-entrypoint-initdb.d/
RUN chmod +x /docker-entrypoint-initdb.d/script.sh
