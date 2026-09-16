#!/bin/bash


# find .git \( -name objects -o -name hooks -o -name __pycache__ \) -prune -o -type f -exec chmod a+w {} +

find .git ! -path ".git/objects/*" ! -path ".git/logs/*" ! -path ".git/hooks/*" -type f -exec echo chmod a+x {} +
