# canopy
The device layer over the forest. Canopy tracks field devices — status, battery, last-seen — and serves them over a small, honest REST API.

## Getting Started

### 1. Clone and configure

```bash
git clone https://github.com/anoop-dryad/canopy.git
cd bridgehead

# copy env config
cp .envrc.example .envrc

# edit .envrc with your local values
# then allow direnv to load it
direnv allow
```

### 2. Set up the database

```bash
# start postgres (Homebrew)
brew services start postgresql@18

# create DB and app user (uses .envrc values)
bash app/scripts/setup_db.sh
```

### 3. Run migrations

```bash
make migrate
```

### 4. Start the server

```bash
make run
```

Server starts on `http://localhost:8080`

### 5. Open Swagger UI

```
http://localhost:8080/swagger/v1/index.html
```
