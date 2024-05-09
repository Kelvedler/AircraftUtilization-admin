PARENT_PATH=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )
docker run -v $PARENT_PATH/migration:/migration --network host migrate/migrate -path=/migration/ -database "postgres://local_user:local_pass@127.0.0.1:5433/local_db?sslmode=disable" up 1
