<p align="center">
  <img src="assets/parchment.webp" width="150" alt="Parchment Logo">
</p>


<h1 align="center">Parchment</h1>
<p align="center">🖨️ A minimalist receipt printing server built with Go.</p>

## 🛠️ Installation

```bash
git clone https://github.com/your-username/parchment.git
cd parchment
go build -o parchment main.go
```

## Development

If you have no thermal printer, run `export MODE=dev` before running your code.

## Deployment

This project is intended to be deployed on my rasbperry pi, and tunneled via cloudlflared with systemctl
To deploy, run `./deploy.sh`

## Feature List

- [x] one script deploy
- [x] POST /print available 
- [x] Possibility to run it without a serial printer
- [x] Prevent content injection on POST /print
- [ ] Protect routes with a token
- [ ] Parse content as markdown
