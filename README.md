# 📁 Go Static File Server

A simple static file server written in Go, inspired by Python’s built-in http.server.

## 📌 About

The purpose of this repository is to implement a lightweight static file server in Go, following the same idea as Python’s `http.server`, developed for the Web Development I course.

It allows you to:

* Serve static files directly from your machine
* Browse directories through a clean HTML interface
* View files in the browser with correct MIME types

## 🚀 Features

* 📂 Directory listing with navigation (`..` support)
* 📄 File serving with automatic MIME type detection
* 🎨 Simple HTML templates for rendering pages
* 🛡️ Protection against directory traversal attacks
* 🧩 SVG icons for files and folders
* ⚡ Minimal and dependency-free (standard library only)

## 🛠️ How It Works

This project uses Go’s standard library, especially:

* `net/http` → to create the web server
* `os` and `filepath` → to interact with the file system
* `mime` → to detect content types
* `strings` → to dynamically build HTML responses

### Core Logic

* The server starts on port **8080**
* The root directory is automatically set to the executable’s location
* Every request:

    1. Is sanitized using `path.Clean` (to prevent unsafe paths)
    2. Is validated to ensure it stays inside the root directory
    3. Is handled differently depending on the request:

        * **Directory** → generates an HTML listing
        * **File** → serves the file content

### Directory Listing

* Built dynamically using HTML templates
* Each file/folder is rendered as a clickable link
* Includes:

    * Folder navigation (`..`)
    * SVG icons for better visualization

## ▶️ Running the Project

### 1. Clone the repository

```bash
git clone https://github.com/dgd-yz/static-server-go.git
cd static-server-go
```

### 2. Run the server

```bash
go run main.go
```

### 3. Open in browser

```
http://localhost:8080
```

## 📷 Example

When accessing a directory, you'll see a listing like:

```
📁 folder/
📄 file.txt
```

## 🔒 Security Notes

* The server prevents access outside the root directory
* Any attempt to escape the base path returns a **403 Forbidden**

## 📚 Learning Goals

This project was built to practice:

* HTTP servers in Go
* File system manipulation
* Dynamic HTML generation
* Basic web security concepts

## 📄 License

This project is open-source and available under the MIT License.