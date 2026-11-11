@echo off
title Omnix Backend
cd /d "%~dp0"
go run ./cmd/server
pause
