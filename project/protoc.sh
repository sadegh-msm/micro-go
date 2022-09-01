#!/bin/sh

service=$1

if [ $service = "auth" ] ; then
  cd ./../authentication-service/auth
  . ./auth_proto.sh
  cd ./../../project

elif [ $service = "log" ]; then
  cd ./../logger-service/logs
  . ./logs.sh
  cd ./../../project

elif [ $service = "shortner" ]; then
  cd ./../urlshortner-service/shortner
  . ./shortner.sh
  cd ./../../project
fi
