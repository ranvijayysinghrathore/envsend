# 🎯 START HERE - Simple Setup Instructions

## What You Need to Do

### 1️⃣ Install Go and Docker

**Install Go:**
1. Go to: https://go.dev/dl/
2. Download: `go1.21.x.windows-amd64.msi`
3. Run installer → Click "Next" through everything
4. **Close and reopen your terminal**

**Install Docker Desktop:**
1. Go to: https://www.docker.com/products/docker-desktop/
2. Download Docker Desktop for Windows
3. Run installer
4. **Restart your computer**
5. **Start Docker Desktop** from Start menu (wait for it to fully start)

---

### 2️⃣ Check Installation

Open PowerShell or Command Prompt in the envsend folder and run:

```cmd
.\check-prerequisites.bat
```

This will tell you if everything is installed correctly.

---

### 3️⃣ Run Setup

Once Go and Docker are installed:

```cmd
.\setup-local.bat
```

This will:
- Download dependencies
- Start databases
- Build the CLI and servers

**Takes about 5 minutes**

---

### 4️⃣ Start the Server

```cmd
bin\api-server.exe
```

Leave this running!

---

### 5️⃣ Test It (in a new terminal)

```cmd
cd "d:\Programs\Next js {projects}\envsend"

REM Create a test file
echo API_KEY=test123 > test.env

REM Send it
bin\envsend.exe send test.env --server http://localhost:8080

REM You'll get a URL - copy it and retrieve:
bin\envsend.exe receive "YOUR_URL_HERE" > retrieved.env

REM Check it worked
type retrieved.env
```

---

## ❓ Troubleshooting

**"go is not recognized"**
→ Restart your terminal after installing Go

**"docker is not recognized"**
→ Make sure Docker Desktop is running (check system tray)

**Scripts don't run**
→ Use `.\script-name.bat` instead of just `script-name.bat`

---

## 📁 Files You Need

- `check-prerequisites.bat` - Check if Go/Docker are installed
- `setup-local.bat` - Automated setup
- `SETUP_GUIDE.md` - Detailed instructions

---

**That's it! Once setup is done, you have a secure secret-sharing system running locally! 🎉**
