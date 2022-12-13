#!/bin/bash

CGO_ENABLED=0 go build
mv term-util.nvim bin/
