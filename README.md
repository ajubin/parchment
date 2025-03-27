<p align="center">
  <img src="assets/logo.svg" width="150" alt="Parchment Logo">
</p>


<h1 align="center">Parchment</h1>
<p align="center">🖨️ A minimalist receipt printing server built with Go.</p>

## 🛠️ Installation

```bash
git clone https://github.com/your-username/parchment.git
cd parchment
```

## Development

Run  `make dev` to start the server in watch mode

## Deployment

This project is intended to be deployed on my rasbperry pi, and tunneled via cloudlflared with systemctl.

Before deploying into production, protect your route with basic auth.
To configure the password, run `cp .env.example .env` and modify `API_USER` and `API_PASSWORD`.

To deploy, run `make deploy`

## Feature List

- [x] one script deploy
- [x] POST /print available 
- [x] Possibility to run it without a serial printer
- [x] Prevent content injection on POST /print
- [x] Protect routes with basic auth
- [ ] Parse content as markdown
