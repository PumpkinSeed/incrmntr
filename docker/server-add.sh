#!/bin/sh

couchbase-cli server-add -c cb1:8091 \ 
--server-add=cb2:8091 --server-add-username=Administrator \ 
--server-add-password=password --services=data