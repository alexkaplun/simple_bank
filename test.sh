#!/usr/bin/env bash
go test -coverprofile=test_out/cover.out ./... && go tool cover -html=test_out/cover.out -o test_out/cover.html
