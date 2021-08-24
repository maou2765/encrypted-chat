#! /bin/bash
# docker run -d --name mysql-container --env="MYSQL_ROOT_PASSWORD=my-secret-pw" --publish 6603:3306 mysql
RUNNING_CONTAINER=$(docker container ls|grep 'mysql:latest')
echo $RUNNING_CONTAINER
echo "$RUNNING_CONTAINER"
if [ -z "$RUNNING_CONTAINER" ]
then
  CONTAINER=$(docker container ls -a|grep 'mysql:latest')
  if [ -z "$CONTAINER" ]
  then
    docker run -d --name mysql-container --env="MYSQL_ROOT_PASSWORD=my-secret-pw" --publish 6603:3306 mysql
  else
    CONTAINER=${CONTAINER:0:12}
    docker container start $CONTAINER
  fi
fi
npx nodemon --watch './**/*.go' --signal SIGTERM --exec 'go' run main.go
