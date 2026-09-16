#!/bin/bash


find .git \( -name objects -o -name hooks -o -name __pycache__ \) -prune -o -type f -exec chmod a+w {} +
