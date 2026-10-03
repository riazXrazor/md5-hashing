# MD5 Implementation in Go

A lightweight, pure Go implementation of the **MD5 (Message-Digest Algorithm)** as specified in [RFC 1321](https://tools.ietf.org/html/rfc1321).

> [!NOTE]
> **Disclaimer / Learning Project:**  
> This was coded purely for learning and understanding how the algorithm works under the hood. It was implemented while referencing Wikipedia and YouTube video explanations, so the code is unoptimized and messy. It is not intended for production use.

## Features

- Computes standard 128-bit MD5 hashes.
- Supports input from both **files** and **standard input (`stdin`)**.
- Built-in debug logging via environment variables.
- Zero external dependencies (standard library only).

---

## Requirements

- [Go](https://golang.org/dl/) 1.22 or newer

---

## How to Run

### 1. Run Directly with `go run`

**Hash a file:**
```bash
go run main.go path/to/file.txt
```

**Hash from standard input (stdin):**
```bash
echo -n "hello" | go run main.go
```

---

### 2. Build and Run Binary

**Build the executable:**
```bash
go build -o md5 main.go
```

**Run using the compiled binary:**
```bash
# Using a file
./md5 test_files/hello.txt

# Using standard input
echo -n "hello" | ./md5
```

---

## Debug Mode

You can enable verbose step-by-step debug logging by setting `DEBUG=1`:

```bash
DEBUG=1 go run main.go test_files/hello.txt
```

---

## Running Tests

Run the included unit test suite:

```bash
go test -v ./...
```
