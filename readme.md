# Byn Payment System Application

This is a simple payment system application implemented in Go. It provides functionalities to manage accounts, perform transactions, and track balances.

## Requirements

- Go programming language 1.18 or higher

## Usage

1. Run the application with usage scenarios:

    ```bash
    go run .
    ```

2. Run utilities tests (iban generation in parallel threads, incremental generator, paddleft function):
    ```bash
    go test ./src/util
    ```
3. Run test that checks that payment service is singleton
   ```bash
   go test ./src/byn
   ```

## Demo

You can run this code on [replit.com](https://replit.com/@nikitabort22092/WebByn#main.go). Just open the link, click green button "Fork & Run" is the top right corner of the screen. Wait about 40sec and you will see output of some examples

If you have any suggestions, bug reports, or feature requests, please contact me by [email](mailto:nikitabort22092000@gmail.com) or phone (if you know it)