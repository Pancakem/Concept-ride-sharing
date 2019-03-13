#!/bin/bash

mkdir keys

cd keys

ssh-keygen -t rsa -N "" -f key
