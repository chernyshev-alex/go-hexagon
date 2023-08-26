# Database
MYSQL_USER ?= user
MYSQL_PASSWORD ?= password
MYSQL_ADDRESS ?= 127.0.0.1:3306
MYSQL_DATABASE ?= article

export PATH   := $(PWD)/bin:$(PATH)
export SHELL  := bash
export OSTYPE := $(shell uname -s)

up: dev-env dev-air           
down: docker-stop              
destroy: docker-teardown clean  


go-generate: $(MOCKERY) 
	go generate ./...
