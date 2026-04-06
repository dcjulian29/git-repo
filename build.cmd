@echo off

echo Build:
goreleaser --snapshot --clean

echo.
echo Lint:
golangci-lint run

echo.
echo Revive:
revive --formatter stylish ./...
