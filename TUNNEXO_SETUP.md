# Tunnexo Agent — Complete Setup Guide

macOS, Linux aur Windows par local application ko public URL se share karne ki step-by-step guide.

- Website: https://tunnexo.live
- Downloads: https://github.com/SharmaG-Dev/Portune-Agent/releases
- Public URL format: `https://<subdomain>.tunnexo.live`
- Command name: `tunnexo`

## Current setup ek nazar mein

| Cheez | User ko kya karna hai? |
| --- | --- |
| Installation | Apne OS/CPU ka archive extract karke binary PATH mein install karein. |
| Backend URL | Kuch set nahi karna: `https://tunnexo.live/tunnel` built in hai. |
| Agent name | Kuch set nahi karna: local IP aur `tunnexo` prefix se automatically banta hai. |
| Authentication | Abhi valid server-issued token environment ya optional `.env` se dena hai. |
| Local application | Apni app start karke uska URL `--target` mein dein. |
| Public URL | Terminal mein server ka assigned URL copy karein. |

**Is build ka working flow:** Download → Install → Token set karein → Run. Automatic token registration abhi implement nahi hui hai; token ke bina command setup error dega. `.env` banana mandatory nahi hai.

### Installation ke baad quick start

Local app port 4000 par running ho aur token maintainer se mil gaya ho, tab agent wali terminal mein:

**macOS / Linux:**

```sh
export TUNNEXO_AGENT_TOKEN='REPLACE_WITH_YOUR_REAL_AGENT_TOKEN'
tunnexo agent --target http://localhost:4000
```

**Windows PowerShell:**

```powershell
$env:TUNNEXO_AGENT_TOKEN = 'REPLACE_WITH_YOUR_REAL_AGENT_TOKEN'
tunnexo agent --target http://localhost:4000
```

Placeholder ki jagah actual token aur `4000` ki jagah app ka port use karein. Token wali command shell history mein save ho sakti hai; editor se optional `.env` banane ka alternative section 6 mein hai.

### Guide mein kya milega?

- Sections 1–2: agent ka kaam aur requirements.
- Sections 3–5: macOS, Linux aur Windows installation.
- Section 6: token setup aur optional configuration.
- Sections 7–9: local app, tunnel start, stop aur restart.
- Sections 10–13: command reference, daily use, updates aur troubleshooting.
- Sections 14–15: maintainer release notes aur license.


## 1. Tunnexo kya karta hai?

Maan lijiye aapki website aapke computer par `http://localhost:4000` par chal rahi hai. Tunnexo agent local application ko tunnel server se connect karta hai. Server ek public URL assign karta hai, jaise `https://demo123.tunnexo.live`.

Koi is public URL ko kholta hai to request tunnel server se aapke agent aur phir aapki local application tak jaati hai. Application ka response isi connection se visitor tak wapas jaata hai.

`demo123` sirf example hai. Actual URL terminal mein server ke response se milega. Subdomain server assign karta hai; is CLI mein custom subdomain choose karne ka flag nahi hai. Restart ke baad same URL milna guaranteed nahi hai.

## 2. Setup se pehle kya chahiye?

1. Internet/network connection wala computer.
2. Apne operating system aur processor ka Tunnexo archive.
3. Filhaal maintainer/server administrator se mila valid agent token (automatic registration pending).
4. Chalti hui local application, jaise `http://localhost:4000`.
5. macOS/Linux par Terminal; Windows par PowerShell.

Downloaded executable chalane ke liye Go install karna zaroori nahi hai. Go sirf source code build karne ke liye chahiye.

| System / processor | Archive |
| --- | --- |
| macOS Apple Silicon (M-series), ARM64 | `tunnexo_darwin_arm64.zip` |
| macOS Intel, AMD64 | `tunnexo_darwin_amd64.zip` |
| Linux x86-64, AMD64 | `tunnexo_linux_amd64.tar.gz` |
| Linux ARM64 / aarch64 | `tunnexo_linux_arm64.tar.gz` |
| Windows x64, AMD64 | `tunnexo_windows_amd64.zip` |

Windows ARM64 aur 32-bit systems ke liye native package is release script mein nahi hai.

**Guide ka order:** Apne OS ka section 3, 4 ya 5 follow karein. Phir sabhi users ke liye section 6 se aage ke steps follow karein.

## 3. macOS setup

### 3.1 Processor check karein

```sh
uname -m
```

- `arm64`: Apple Silicon archive download karein.
- `x86_64`: Intel archive download karein. Apple Silicon par translated terminal bhi `x86_64` dikha sakta hai; Apple menu → About This Mac mein chip confirm karein.

### 3.2 Download aur extract karein

[Releases page](https://github.com/SharmaG-Dev/Portune-Agent/releases) se correct ZIP download karke `Downloads` folder mein rakhein. Maintainer ne ZIP directly bheja hai to woh bhi isi folder mein rakh sakte hain.

**Apple Silicon:**

```sh
mkdir -p "$HOME/Downloads/tunnexo"
unzip "$HOME/Downloads/tunnexo_darwin_arm64.zip" -d "$HOME/Downloads/tunnexo"
cd "$HOME/Downloads/tunnexo"
```

**Intel — upar ke Apple Silicon commands ki jagah yeh chalayein:**

```sh
mkdir -p "$HOME/Downloads/tunnexo"
unzip "$HOME/Downloads/tunnexo_darwin_amd64.zip" -d "$HOME/Downloads/tunnexo"
cd "$HOME/Downloads/tunnexo"
```

| Command | Kya karta hai? |
| --- | --- |
| `mkdir -p` | Extraction folder banata hai; existing folder ho to use karta hai. |
| `unzip ... -d ...` | Downloaded ZIP ko specified folder mein extract karta hai. |
| `cd ...` | Terminal ko extracted folder mein le jaata hai. |

Agar browser ne ZIP automatically extract kar diya ho, extracted folder mein terminal kholkar next step karein. Agar file ka naam/path alag hai to command mein wahi actual path use karein.

### 3.3 Executable check aur install karein

```sh
chmod +x ./tunnexo
./tunnexo --help
mkdir -p "$HOME/.local/bin"
cp ./tunnexo "$HOME/.local/bin/tunnexo"
export PATH="$HOME/.local/bin:$PATH"
tunnexo version
```

| Command | Kya karta hai? |
| --- | --- |
| `chmod +x` | Binary ko run karne ki permission deta hai. |
| `./tunnexo --help` | Current folder ki binary chala kar commands dikhata hai. |
| `cp ...` | Binary ko aapke user ke install folder mein copy karta hai. |
| `export PATH=...` | Current terminal ko batata hai ki `tunnexo` binary kahan milegi. |
| `tunnexo version` | Installed binary ka version print karta hai. |

Future terminals mein bhi command available rakhne ke liye, default zsh mein:

```sh
nano "$HOME/.zshrc"
```

File mein yeh line ek baar add karein:

```sh
export PATH="$HOME/.local/bin:$PATH"
```

Nano mein `Ctrl+O`, Enter se save aur `Ctrl+X` se exit karein. Phir:

```sh
source "$HOME/.zshrc"
command -v tunnexo
```

`source` shell configuration reload karta hai. `command -v` selected executable ka path dikhata hai. Agar aap Bash use karte hain, apni Bash startup file mein PATH line add karein.

macOS agar downloaded binary ko block kare, source verify karne ke baad System Settings → Privacy & Security mein us application ke liye approval option check karein.

**Ab section 6 par jaayein.**

## 4. Linux setup

### 4.1 Architecture check karein

```sh
uname -m
```

- `x86_64`: AMD64 archive.
- `aarch64` / `arm64`: ARM64 archive.

### 4.2 Download aur extract karein

[Releases page](https://github.com/SharmaG-Dev/Portune-Agent/releases) se archive download karke `Downloads` folder mein rakhein. Neeche ke paths ko apne actual download folder ke hisaab se badal sakte hain.

**AMD64:**

```sh
mkdir -p "$HOME/Downloads/tunnexo"
tar -xzf "$HOME/Downloads/tunnexo_linux_amd64.tar.gz" -C "$HOME/Downloads/tunnexo"
cd "$HOME/Downloads/tunnexo"
```

**ARM64 — AMD64 commands ki jagah:**

```sh
mkdir -p "$HOME/Downloads/tunnexo"
tar -xzf "$HOME/Downloads/tunnexo_linux_arm64.tar.gz" -C "$HOME/Downloads/tunnexo"
cd "$HOME/Downloads/tunnexo"
```

`tar -xzf` compressed archive extract karta hai; `-C` destination folder select karta hai.

### 4.3 Install karein

```sh
chmod +x ./tunnexo
./tunnexo --help
mkdir -p "$HOME/.local/bin"
cp ./tunnexo "$HOME/.local/bin/tunnexo"
export PATH="$HOME/.local/bin:$PATH"
tunnexo version
```

Yeh executable permission set karta hai, help verify karta hai, binary install karta hai aur current terminal ka PATH update karta hai. User folder mein install hota hai, isliye `sudo` ki zaroorat nahi hai.

Default Bash ke future interactive terminals ke liye:

```sh
nano "$HOME/.bashrc"
```

Yeh line ek baar add karein:

```sh
export PATH="$HOME/.local/bin:$PATH"
```

Save karke exit karein, phir:

```sh
source "$HOME/.bashrc"
command -v tunnexo
```

Zsh use karte hain to `.bashrc` ki jagah `.zshrc` use karein.

**Ab section 6 par jaayein.**

## 5. Windows setup (PowerShell)

### 5.1 Download karein

Windows Settings → System → About → System type mein x64 system confirm karein. [Releases page](https://github.com/SharmaG-Dev/Portune-Agent/releases) se `tunnexo_windows_amd64.zip` download karke `Downloads` mein rakhein.

PowerShell kholein. Neeche ke commands PowerShell ke liye hain, Command Prompt ke liye nahi.

### 5.2 Extract aur install karein

```powershell
$tunnexoZip = Join-Path $HOME 'Downloads\tunnexo_windows_amd64.zip'
$tunnexoExtract = Join-Path $HOME 'Downloads\tunnexo'
Expand-Archive -LiteralPath $tunnexoZip -DestinationPath $tunnexoExtract -Force
Set-Location $tunnexoExtract
.\tunnexo.exe --help

$tunnexoInstall = Join-Path $env:LOCALAPPDATA 'Tunnexo'
New-Item -ItemType Directory -Path $tunnexoInstall -Force | Out-Null
Copy-Item '.\tunnexo.exe' -Destination $tunnexoInstall -Force
```

| Command | Kya karta hai? |
| --- | --- |
| `Join-Path` | Aapke user folders ke saath correct Windows path banata hai. |
| `Expand-Archive` | ZIP extract karta hai; `-Force` matching extracted files replace kar sakta hai. |
| `Set-Location` | Working folder change karta hai, `cd` ki tarah. |
| `.\tunnexo.exe --help` | Current folder ki executable ki help dikhata hai. |
| `New-Item` | User ka installation folder banata hai. |
| `Copy-Item` | Binary install karta hai; `-Force` previous binary replace karta hai. |

### 5.3 PATH set karein

Isi PowerShell window mein:

```powershell
$tunnexoUserPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if (($tunnexoUserPath -split ';') -notcontains $tunnexoInstall) {
    $tunnexoNewPath = (@($tunnexoUserPath, $tunnexoInstall) | Where-Object { $_ }) -join ';'
    [Environment]::SetEnvironmentVariable('Path', $tunnexoNewPath, 'User')
}
$env:Path = "$tunnexoInstall;$env:Path"
Get-Command tunnexo
tunnexo version
```

Yeh existing user PATH preserve karke install folder add karta hai. `$env:Path` current PowerShell window ko turant update karta hai. `Get-Command` batata hai kaunsi binary run hogi. Future sessions ke liye terminal application close karke dobara open karein; zaroorat pade to sign out/in karein.

User folder installation ke liye administrator PowerShell ki zaroorat nahi hai.

**Ab section 6 follow karein.**

## 6. Authentication aur automatic defaults — sabhi platforms

Server URL aur agent name prefix user ko set nahi karne hain:

| Setting | Automatic default |
| --- | --- |
| Backend endpoint | `https://tunnexo.live/tunnel` |
| Agent name prefix | `tunnexo` |
| Public URL | Server-assigned `https://<subdomain>.tunnexo.live` |

Agent name local IP aur prefix se automatically banta hai. Prefix public subdomain select nahi karta.

**Token abhi required hai.** Neeche ke do methods mein se ek choose karein. Agent token automatically generate ya permanently store nahi karta.

### Method A: current terminal mein token set karein — `.env` ki zaroorat nahi

Maintainer se valid server-issued token lein. Agent chalane wali terminal mein:

macOS/Linux:

```sh
export TUNNEXO_AGENT_TOKEN='REPLACE_WITH_YOUR_REAL_AGENT_TOKEN'
```

Windows PowerShell:

```powershell
$env:TUNNEXO_AGENT_TOKEN = 'REPLACE_WITH_YOUR_REAL_AGENT_TOKEN'
```

Yeh current terminal ka token set karta hai; doosri terminal automatically inherit nahi karti. Isi terminal mein section 8 ka agent command chalayein. Real token wali command shell history mein store ho sakti hai; file method neeche available hai.

### Method B: optional `.env` file mein token save karein

Yeh alternative hai, mandatory setup nahi. Agent `.env` ko current working folder mein padhta hai, executable ke installation folder mein automatically nahi.

macOS/Linux:

```sh
mkdir -p "$HOME/tunnexo-work"
cd "$HOME/tunnexo-work"
(umask 077; touch .env)
nano .env
```

Windows PowerShell:

```powershell
New-Item -ItemType Directory -Path "$HOME\tunnexo-work" -Force | Out-Null
Set-Location "$HOME\tunnexo-work"
notepad .env
```

Editor mein sirf token add karein:

```dotenv
TUNNEXO_AGENT_TOKEN=REPLACE_WITH_YOUR_REAL_AGENT_TOKEN
```

Nano mein `Ctrl+O`, Enter, `Ctrl+X` se save/exit karein. Windows Notepad mein filename `.env` rakhein, `.env.txt` nahi; zaroorat ho to Save as type **All files** choose karein.

macOS/Linux par:

```sh
chmod 600 .env
```

Yeh config ko owner-only read/write permissions deta hai. Token public repository, screenshot ya shared guide mein na bhejein.

### Advanced users: optional server/prefix overrides

Normal public-server setup mein yeh step skip karein. Custom backend ya prefix use karna ho tabhi optional `.env` mein set karein:

```dotenv
TUNNEXO_SERVER_URL=https://tunnexo.live/tunnel
TUNNEXO_AGENT_NAME_PREFIX=tunnexo
```

Prefix mein sirf letters, digits aur hyphens allowed hain. Configuration priority har setting ke liye:

1. `TUNNEXO_*` environment value.
2. `.env` mein `TUNNEXO_*` value.
3. Legacy `PORTUNE_*` environment value.
4. `.env` mein legacy `PORTUNE_*` value.
5. Missing/empty server aur prefix ke liye built-in defaults.

Explicit empty new-name value legacy fallback ko stop karta hai. Empty token par clear authentication setup error aayega. Purana `PORTUNE_SERVER_URL` abhi override kar sakta hai; old Render endpoint hata dein ya Tunnexo endpoint se replace karein.

## 7. Local application start aur check karein

**Terminal 1** mein apni application uske normal start command se chalayein. Tunnexo aapki application khud start nahi karta.

Browser mein `http://localhost:4000` kholkar check karein. Agar app port `3000`, `5173` ya kisi aur port par chalti hai, wahi address use karein.

Optional command-line check:

**macOS / Linux:**

```sh
curl -i http://localhost:4000
```

**Windows PowerShell:**

```powershell
Test-NetConnection -ComputerName localhost -Port 4000
Invoke-WebRequest -Uri 'http://localhost:4000'
```

`curl` / `Invoke-WebRequest` local server ka HTTP response check karte hain. `Test-NetConnection` sirf port reachable hai ya nahi batata hai. App ki `/` route 404 de sakti hai; aise mein app ki known working route check karein. Connection refused aaye to app start karein ya correct port use karein.

### Optional: sirf demo ke liye local web server

Agar Python installed hai aur koi app nahi hai, **Terminal 1** mein ek alag demo folder se server chala sakte hain.

macOS/Linux:

```sh
mkdir -p "$HOME/tunnexo-demo"
cd "$HOME/tunnexo-demo"
printf '%s\n' '<h1>Hello from Tunnexo</h1>' > index.html
python3 -m http.server 4000 --bind 127.0.0.1
```

Windows PowerShell (Python launcher installed ho):

```powershell
New-Item -ItemType Directory -Path "$HOME\tunnexo-demo" -Force | Out-Null
Set-Location "$HOME\tunnexo-demo"
Set-Content -Path index.html -Value '<h1>Hello from Tunnexo</h1>'
py -m http.server 4000 --bind 127.0.0.1
```

Yeh `index.html` banata/replaces karta hai aur demo folder ko port 4000 par serve karta hai. Is demo ko secrets ya personal files wale folder se na chalayein. Python installed nahi hai to apni existing application use karein; Python agent ki requirement nahi hai.

## 8. Tunnel start karein

**Terminal 2** mein token set karein (section 6), phir command chalayein. Environment-based token ke saath kisi bhi folder se run kar sakte hain. Optional `.env` use karte hain to pehle uske folder mein jaayein.

macOS/Linux:

```sh
tunnexo agent --target http://localhost:4000
```

Windows PowerShell:

```powershell
tunnexo agent --target http://localhost:4000
```

| Part | Meaning |
| --- | --- |
| `tunnexo` | Installed executable chalata hai. |
| `agent` | Tunnel agent start karta hai. |
| `--target` | Local application ka address specify karta hai; required hai. |
| `http://localhost:4000` | Requests is address par forward hongi. Apna actual port use karein. |

Command configuration load karega, server se authenticate karega aur tunnel request karega. Success par terminal mein connection details aur **Public URL** dikhega. Wahi exact URL copy karke share karein.

Example URL: `https://demo123.tunnexo.live` — yeh reserved ya guaranteed working URL nahi hai.

Local app aur agent dono running rakhein. Agent terminal close karne, internet disconnect hone ya computer sleep hone par tunnel unavailable ho sakta hai. Public URL share karne se local app doosron ke liye reachable hoti hai; app ka required authentication enabled rakhein.

### Common target examples

Inmein se apni app ka **ek** command use karein:

```sh
tunnexo agent --target http://localhost:3000
tunnexo agent --target http://localhost:4000
tunnexo agent --target http://localhost:5173
tunnexo agent --target http://127.0.0.1:8080
```

Target mein `http://` ya `https://` required hai. Bare `localhost:4000`, embedded username/password, query string aur fragment allowed nahi hain. Browser se aane wali request query strings alag cheez hain; yeh restriction configured `--target` value par hai.

### PATH install skip kiya ho to

Extracted folder mein terminal kholein. Section 6 ke Method A se isi terminal mein token set karein, ya Method B se isi folder mein `.env` banayein. Phir:

macOS/Linux:

```sh
./tunnexo agent --target http://localhost:4000
```

Windows PowerShell:

```powershell
.\tunnexo.exe agent --target http://localhost:4000
```

`./` aur `.\` current folder ki executable select karte hain.

## 9. Stop aur restart

Agent wale **Terminal 2** mein `Ctrl+C` press karein. Agent disconnect hoga; **Terminal 1** ki local app alag se running reh sakti hai. Local app bhi stop karni ho to uske terminal mein `Ctrl+C` karein.

Restart ke liye token set wali terminal mein wahi command dobara chalayein. Nayi terminal mein Method A ka token phir set karna hoga. Method B use karte hain to `.env` wale folder mein jaayein:

```sh
tunnexo agent --target http://localhost:4000
```

Restart ke baad terminal se fresh public URL check karein.

## 10. Command reference

| Command | Result |
| --- | --- |
| `tunnexo` | Root help dikhata hai; tunnel start nahi karta. |
| `tunnexo --help` | Available commands aur root flags. |
| `tunnexo help agent` | Agent command ki help. |
| `tunnexo agent --help` | Agent usage aur `--target` flag. |
| `tunnexo version` | Binary mein embedded version. |
| `tunnexo agent --target http://localhost:4000` | Local port 4000 ke liye agent start. |
| `tunnexo completion --help` | Shell completion generation ki help. |
| `Ctrl+C` | Running agent ko stop karta hai. |

Is build mein `login`, `config`, `--subdomain` aur `--version` commands/flags implement nahi hain. Version ke liye `tunnexo version` aur token ke liye `.env`/environment use karein.

## 11. Daily use — short checklist

1. Local app start karein aur uska port check karein.
2. Dusra terminal kholein.
3. Us terminal mein token set karein, ya optional `.env` ke folder mein jaayein.
4. `tunnexo agent --target http://localhost:4000` run karein; actual port substitute karein.
5. Terminal ka public URL share karein.
6. Kaam complete hone par agent mein `Ctrl+C` karein.

Installation aur PATH setup roz repeat karne ki zaroorat nahi hai.

## 12. Update kaise karein?

1. Running agent ko `Ctrl+C` se stop karein.
2. Apne platform ka naya archive Releases page se download aur extract karein.
3. Extracted folder se installed executable replace karein.

macOS/Linux:

```sh
chmod +x ./tunnexo
cp ./tunnexo "$HOME/.local/bin/tunnexo"
tunnexo version
```

Windows PowerShell:

```powershell
Copy-Item '.\tunnexo.exe' -Destination "$env:LOCALAPPDATA\Tunnexo\tunnexo.exe" -Force
tunnexo version
```

Yeh commands extracted folder mein run karein. Optional configuration `tunnexo-work/.env` mein rakhi hai to binary replace karne se config change nahi hogi. Phir token set wali terminal se agent restart karein; optional `.env` use karte hain to uske folder se run karein.

Auto-update command available nahi hai. GitHub release tag change karne se existing downloaded binary ka version nahi badalta; nayi binary install karni hoti hai.

## 13. Troubleshooting

| Problem | Kya check / fix karein? |
| --- | --- |
| `command not found` / command not recognized | PATH step complete karein; new terminal kholein. Extracted folder mein `./tunnexo` ya `.\tunnexo.exe` try karein. |
| `permission denied` on macOS/Linux | Executable par `chmod +x ./tunnexo` chalayein. |
| `exec format error` / incompatible app | OS aur CPU ke liye correct archive download karein. |
| Download link 404 | Releases page check karein; archive abhi publish nahi hua ho to maintainer se file lein. |
| `automatic agent registration is not configured yet` | Section 6 se valid token supply karein. Automatic token registration is build mein available nahi hai. Optional file use karte hain to folder aur `.env` filename check karein. |
| File mein correct value ke baad bhi old config use ho rahi hai | Same-name environment variable `.env` ko override karta hai. Current session ka stale variable clear karein. |
| Cannot detect a non-loopback local IP | Active network connection check karein; agent ko usable local IP chahiye. |
| Connection / authentication failure | Endpoint, token, internet aur server availability verify karein. |
| Could not resolve host / DNS error | Domain spelling aur DNS/network check karein; production domain readiness maintainer se confirm karein. |
| Public URL par 502 / local connection refused | Local app running ho aur `--target` ka scheme/host/port sahi ho. |
| App public hostname reject karti hai | Local dev server ki allowed-host/origin configuration mein assigned hostname allow karein. |
| Browser HTTPS certificate error | Tunnel server ke wildcard certificate aur routing ke liye maintainer se contact karein. |
| Tunnel sleep ke baad unavailable | Computer awake aur network connected rakhein; agent restart karke returned URL check karein. |
| Agent prefix / target validation error | Prefix mein letters/digits/hyphens; target mein valid HTTP(S) address use karein. |

Stale environment setting clear karne ka example — jab built-in public backend ya optional `.env` wali server value use karni ho:

macOS/Linux:

```sh
unset TUNNEXO_SERVER_URL
```

Windows PowerShell:

```powershell
Remove-Item Env:TUNNEXO_SERVER_URL -ErrorAction SilentlyContinue
```

Yeh current terminal ka override remove karta hai; `.env` ko delete nahi karta. Doosri stale setting ke liye uska exact variable name use karein. Legacy `PORTUNE_*` fallback bhi check karein. Built-in default use karne ke liye environment aur optional `.env` dono se custom server overrides remove karein.

Support ko OS, CPU, `tunnexo version`, command aur redacted error output bhejein. Token ya full private `.env` share na karein.

## 14. Maintainer notes: release aur production readiness

End-user installation ke liye is section ke commands run nahi karne hain.

**Current verification limits:** Production HTTPS/tunnel connection verify nahi hua hai. GitHub download links tab active honge jab matching archives Release mein publish honge. Unpublished release ke users ko archive directly provide karein.

End users ko source build karne ki zaroorat nahi hai. Maintainer ke source checkout se:

```sh
go test ./...
go vet ./...
bash scripts/build-release.sh
```

- Pehle do commands tests aur static checks chalate hain.
- Release script 5 platform packages `dist/` mein banata hai; Go 1.27.1+, Bash, `zip` aur `tar` chahiye.
- Script README, yeh guide, example configuration aur MIT license packages mein include karta hai.
- Release script nayi GitHub release/tag create ya upload nahi karta.
- Version abhi `cmd/version.go` mein set hota hai. Version update karke binaries rebuild karein, matching `vX.Y.Z` tag ki GitHub release banayein aur 5 `dist/tunnexo_*` archives upload karein.
- Repository URL abhi `SharmaG-Dev/Portune-Agent` hai; product name Tunnexo hai.

Production URLs ke liye server ka public domain `tunnexo.live`, wildcard DNS `*.tunnexo.live`, matching HTTPS certificate aur tunnel routing configure honi chahiye. Agent server ka returned URL display karta hai; domain ko locally rewrite nahi karta. Automatic token registration ke liye backend API contract aur implementation abhi required hai. Public endpoint aur real token ke saath end-to-end verification bhi pending hai.

## 15. License

Tunnexo Agent MIT License ke under distributed hai. Archive/source ke [LICENSE](LICENSE) file mein complete terms hain. MIT text aur copyright notice ko redistribution mein preserve karein.

---

**Share karne ke liye:** Yeh guide aur recipient ke OS/CPU ka matching archive bhejein. Automatic registration available hone tak valid token separately provide karein; public server URL built in hai. Guide mein placeholder token hi rehne dein.
