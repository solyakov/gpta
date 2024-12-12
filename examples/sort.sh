#!/bin/bash

cat <<EOF | ../gpta -vit "sort this list, capitalizing the first letter of each word"
cherry
apple
banana
EOF