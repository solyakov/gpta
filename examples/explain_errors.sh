#!/bin/bash

ls /non/existent/directory 2>&1 | ../gpta -t "explain this error"