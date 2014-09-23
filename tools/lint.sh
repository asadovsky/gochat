#!/bin/bash

set -e
set -u

find . -name '*.js' -print0 | xargs -0 gjslint --nojsdoc --nobeep
