#!/bin/sh

echo "" > $(docker inspect --format='{{.LogPath}}' $1)
