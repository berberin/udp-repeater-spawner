#! /bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

kill $(ps aux | grep "[r]epeater .*/config/$1" | awk '{print $2}')
rm $DIR/config/$1.json
rm $DIR/log/$1.log