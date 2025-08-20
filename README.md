# 📦 Go-Archiver 🐹

<p align="center">
  <img src="https://raw.githubusercontent.com/Diuxx/resources/7161a6770faf33ecd64923f20692360af29b0f13/go-archiver-logo.svg" alt="Go-Archiver logo" width="580">
</p>

Small project that allows you to archive a specified path to a specified destination.  
It can be used as a **scheduled task** (Windows Task Scheduler, cron, Synology Task Scheduler).  
Fast, portable, written in Go.

## Table of Contents
- [About](#-about)
- [What's New](#-whats-new)
- [Features](#-features)
- [Installation](#-installation)
- [Usage](#-usage)
- [How to Build](#-how-to-build)
- [Examples](#-examples)
- [Contacts](#-contacts)

## 🚀 About
Go-Archiver is a **lightweight CLI tool** to automate folder backups.  
It takes an `origin` path, compresses it (zip), and places the archive in a `destination` folder.  
Typical use cases:
- 🗂️ Regular backups of work folders  
- 🎬 Auto-archiving NAS media libraries  
- 💾 Keeping lightweight snapshots of projects  

## ✨ What's New
- **v0.2.0**  
  - Added support for scheduled execution logs  
  - Improved error handling (disk not mounted detection)  
  - Configurable log prefix via CLI flag  

- **v0.1.0**  
  - First working release: basic archive to destination  


## 🔧 Features
- Compresses any folder into a **.zip archive**  
- Archives named with **date/time suffix** for easy tracking  
- Supports both **local drives** and **network drives/NAS**  
- Automatic **disk mount check** (prevents errors if NAS is offline)  
- Optional **cleanup**: remove source files after archive  
- Logging with customizable prefix  

---

## 📥 Installation
1. [Download latest release](https://github.com/your-username/go-archiver/releases)  
2. Place binary in a folder of your choice (e.g. `C:\tools\go-archiver`)  
3. Add folder to your **PATH** (optional, for global access)

---

## ▶️ Usage

```bash
# Basic usage
go-archiver --origin "C:\Users\me\Documents" --dest "D:\Backups"

# With logging
go-archiver --origin "Z:\films" --dest "D:\Archives" --log-prefix="films"

# Verbose mode
go-archiver --origin "/mnt/nas/data" --dest "/backups" -v
```

## How to Build

# Clone the repo
git clone https://github.com/your-username/go-archiver.git
cd go-archiver

# Build (Windows)
go build -o go-archiver.exe .

# Build (Linux)
GOOS=linux GOARCH=amd64 go build -o go-archiver .

## Examples

# Daily backup via Windows Task Scheduler :
go-archiver.exe --origin "C:\Projects" --dest "E:\DailyBackups"

## Contacts

Author: [Diuxx](https://github.com/Diuxx)
Issues & Suggestions: GitHub Issues
