@echo off
docker build -t calculator .
docker run -p 8888:8080 calculator
pause