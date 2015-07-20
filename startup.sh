#!/bin/bash
set -e
NSQD_HOST="http://172.17.42.1:4151"
MONGODB_URL="mongodb://172.17.42.1/mydb"
case $1 in 
	production)
		NSQD_HOST="http://172.17.42.1:4151"
		MONGODB_URL="mongodb://172.17.42.1/mydb"
		;;
	testing)
		NSQD_HOST="http://172.17.42.1:4151"
		MONGODB_URL="mongodb://172.17.42.1/mydb"
		;;
esac
export NSQD_HOST=$NSQD_HOST
export MONGODB_URL=$MONGODB_URL
echo "NSQD_HOST set to:" $NSQD_HOST
echo "MONGODB_URL set to:" $MONGODB_URL

$GOBIN/auth
