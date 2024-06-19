# timestamp: $(date +'%s.%N')
define GELF_TEST_UDP
echo \
{
  "version": "1.0",
  "host": "localhost",
  "short_message": "Short message",
  "full_message": "Full message",
  "timestamp": "1640456345.940518000",
  "level": 1,
  "facility": "Tester",
  "_user_id": 42,
  "_Environment": "test",
  "_AppName": "Tester"
} | gzip -c -f - | nc -w 1 -u localhost 12201
endef
export GELF_TEST_UDP

define GELF_TEST_TCP
echo \
{
  "version": "1.0",
  "host": "localhost",
  "short_message": "Short message",
  "full_message": "Full message",
  "timestamp": "1640456345.940518000",
  "level": 1,
  "facility": "Tester",
  "_user_id": 42,
  "_Environment": "test",
  "_AppName": "Tester"
} | gzip -c -f - | nc -w 1 localhost 12201
endef
export GELF_TEST_TCP

# Podman
PODNAME := unleash

pd-compose:
	@podman-compose -f docker/docker-compose.yaml -p $(PODNAME) up

pd-machine-init:
	@podman machine init --memory=8192 --cpus=2 --disk-size=20

pd-machine-start:
	@podman machine start

pd-machine-rm:
	@podman machine rm

pd-machine-recreate: pd-machine-rm pd-machine-init pd-machine-start

pd-pod-create:
	@podman pod create -n $(PODNAME) --network bridge \
      	-p 4242:4242 -p 5432:5432

pd-pod-rm:
	@podman pod rm -f $(PODNAME)

pd-pod-recreate: pd-pod-rm pd-pod-create

# Install unleash
#  web:
#    image: unleashorg/unleash-server:latest
#    ports:
#      - "4242:4242"
#    environment:
#      # This points Unleash to its backing database (defined in the `db` section below)
#      DATABASE_URL: "postgres://postgres:unleash@db/db"
#      # Disable SSL for database connections. @chriswk: why do we do this?
#      DATABASE_SSL: "false"
#      # Changing log levels:
#      LOG_LEVEL: "warn"
#      # Proxy clients must use one of these keys to connect to the
#      # Proxy. To add more keys, separate them with a comma (`key1,key2`).
#      INIT_FRONTEND_API_TOKENS: "default:development.unleash-insecure-frontend-api-token"
#      # Initialize Unleash with a default set of client API tokens. To
#      # initialize Unleash with multiple tokens, separate them with a
#      # comma (`token1,token2`).
#      INIT_CLIENT_API_TOKENS: "default:development.unleash-insecure-api-token"
#    depends_on:
#      db:
#        condition: service_healthy
 #   command: [ "node", "index.js" ]
 #   healthcheck:
 ##     test: wget --no-verbose --tries=1 --spider http://localhost:4242/health || exit 1
 #     interval: 1s
 #     timeout: 1m
 #     retries: 5
 #    start_period: 15s

pd-postgres:
	@podman run -dit --name unleash --pod=$(PODNAME) \
		-e "DATABASE_URL=postgres://postgres:unleash@db/db" \
	    -e "DATABASE_SSL=false" \
    	-e "LOG_LEVEL=warn" \
	    -e "INIT_FRONTEND_API_TOKENS=default:development.unleash-insecure-frontend-api-token" \
    	-e "INIT_CLIENT_API_TOKENS=default:development.unleash-insecure-api-token" \
		docker.io/unleashorg/unleash-server:latest

# Install Jaeger
#  db:
#    expose:
#      - "5432"
#    image: postgres:15
#    environment:
#      # create a database called `db`
#      POSTGRES_DB: "db"
#      # trust incoming connections blindly (DON'T DO THIS IN PRODUCTION!)
#      POSTGRES_HOST_AUTH_METHOD: "trust"
#    healthcheck:
#      test:
#        [
#          "CMD",
#          "pg_isready",
#          "--username=postgres",
#          "--host=127.0.0.1",
#          "--port=5432"
#        ]
#      interval: 2s
#      timeout: 1m
#      retries: 5
#      start_period: 10s

pd-jaeger:
	@podman run -dit --name postgres --pod=$(PODNAME) \
		-e "POSTGRES_DB=db" \
		-e "POSTGRES_HOST_AUTH_METHOD=trust" \
		docker.io/postgres:15

pd-init: pd-machine-init pd-machine-start pd-pod-create

pd-start: pd-postgres pd-unleash