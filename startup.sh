#!/bin/bash
set -e
NSQD_HOST="http://172.17.42.1:4151"
REDIS_HOST="172.17.42.1:7100"
MONGODB_URL="mongodb://172.17.42.1/mydb"
ENABLE_SNAPPY=true
case $1 in 
	production)
		NSQD_HOST="http://172.17.42.1:4151"
		REDIS_HOST="172.17.42.1:7100"
		MONGODB_URL="mongodb://172.17.42.1/mydb"
		;;
	testing)
		NSQD_HOST="http://172.17.42.1:4151"
		REDIS_HOST="172.17.42.1:7100"
		MONGODB_URL="mongodb://172.17.42.1/mydb"
		;;
esac
export NSQD_HOST=$NSQD_HOST
export REDIS_HOST=$REDIS_HOST
export MONGODB_URL=$MONGODB_URL
export ENABLE_SNAPPY=$ENABLE_SNAPPY
echo "NSQD_HOST set to:" $NSQD_HOST
echo "REDIS_HOST set to:" $REDIS_HOST
echo "MONGODB_URL set to:" $MONGODB_URL
echo "ENABLE_SNAPPY:" $ENABLE_SNAPPY

# signal handle
trap 'kill -TERM $PID' TERM
$GOBIN/auth &
PID=$!
wait $PID
wait $PID
EXIT_STATUS=$?
