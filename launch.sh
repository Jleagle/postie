#!/usr/bin/env bash

if [ "${ENV}" == "local" ]
then

    dep ensure
    realize start

else

    postie

fi
