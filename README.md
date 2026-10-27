# ConcOrrient

ConcOrrient is a small Go web app for keeping shared ideas alive across different schedules. A thread can be open to the whole workspace or directed to one teammate. Participants can add thoughts, update a thread’s status, and see activity notifications.

## Run it

```powershell
go run .
```

Then open [http://localhost:8080](http://localhost:8080). Go 1.23 or newer is required. The first run creates `data.json` beside the app and seeds a sample workspace. Edit the initial workspace in `internal/workspace/seed.go` before first launch if you want different starting content.

`CONCORRIENT_ADDR` changes the listen address, and `CONCORRIENT_DATA` changes the JSON data file location.

## Technology

- Go 1.23+ with the standard-library `net/http` server; there is no Go web framework.
- Plain HTML, CSS, and JavaScript for the web/PWA client; there is no frontend framework.
- Android native client written in Java with an Android `WebView`, built using Gradle and Android Gradle Plugin 8.13.1.
- JSON file persistence; there is no database configured yet.
- Build tools: Go toolchain, Android SDK/API 36 and Build Tools 36.0.0, Android Studio's JDK 21, and Gradle 8.13.

## Project layout

- `main.go` starts the server, loads embedded web assets, and wires dependencies.
- `internal/httpapi/` owns routes, request decoding, response formatting, and HTTP handlers.
- `internal/workspace/` owns collaboration models, business rules, notifications, and JSON persistence.
- `web/` contains the browser/PWA interface and embedded static assets.
- `android/` contains the native Android client and its Gradle build.

## Install on Android

Build the native Android app from the `android` directory with Android Studio, or run `./gradlew assembleDebug` (Windows: `gradlew.bat assembleDebug`). The installable APK is written to `android/app/build/outputs/apk/debug/app-debug.apk`.

The phone app connects to the Go workspace running on your computer. Keep both devices on the same Wi-Fi network, start the Go server, and enter the computer's local network address (for example, `http://192.168.1.20:8080`) when the app opens. The app's **Server** button lets you change that address later. Windows Firewall may need to allow incoming connections on port 8080. This first native build is an installable Android client, not a standalone server; collaboration data remains on the computer running Go.

## First version

- One shared workspace with a few sample threads and teammates
- Direct threads notify the invited teammate and the other participant when a thought is added
- Open threads notify the rest of the workspace
- Thread status filters, labels, activity history, and notifications
- JSON file persistence with automatic refresh for other open browser tabs

The member picker is a lightweight prototype identity switcher so the collaboration flow can be tried without configuring accounts. Add authentication and a database before exposing this server to a wider network.
