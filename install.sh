#!/bin/bash
cd `dirname $0`
set -e 

export GOBIN=${HOME}/.local/bin
go install 
cd dummywindow
go install 

