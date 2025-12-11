# Api

## Dipendenze
 - `goose` per le migrazioni del database
 - `docker`, se si ha docker e non si vuole eseguire il codice sull'host allora le seguenti dipendenze non servono
 - `github.com/joho/godotenv` per caricare variabili d'ambiente
 - `github.com/gin-gonic/gin` web framework
 - `github.com/golang-jwt/jwt/v5`
 - `github.com/google/uuid`
 - `github.com/jackc/pgx/v5` postgres driver

## Struttura progetto
- `./Auth` contiene i file necessari per avviare [keycloak](https://www.keycloak.org/), identity provider per automatizzare autenticazione e autorizzazione degli utenti. La configurazione di keycloak è specificata in `./Auth/docker-compose.yml`. Le cartelle `imports` e `themes` contengono file di configurazione del server che vengono importati quando viene creato il server.
- `./cmd` contiene solamente il file `main.go`, punto di partenza per l'api.
- `./internal` contiene tutti i file dello sviluppo effettivo dell'api. I file al di fuori delle sotto cartelle sono utilities.
  - `./internal/domain` contiene le strutture utili per rappresentare le entità del progetto.
  - `./internal/infrastructure` contiene file riguardanti connessioni/servizi connessi al database.
    - `./internal/infrastructure/adapters/postgresql/migrations` contiene le migrazioni per il database (postgres). Le migrazioni sono gestite dal package `goose`.
    - `./internal/infrastructure/database` contiene file utili per la connessione dell'api al database.
    - `./internal/infrastructure/services` contiene file per inserire, ottenere, modificare o eliminare dati sul database dall'api.
  - `./internal/presentation` contiene file utili per la configurazione degli endpoint dell'api.

## Avviare il progetto il locale
Per avviare il server keycloak, creare un file `.env` all'interno di `./Auth` con le seguenti variabili d'ambiente:
  - `DB`: nome del db utilizzato da keycloak
  - `DB_USER`: nome dell'utente postgres
  - `DB_USER_PASSWORD`: password dell'utente postgres
  - `KC_HOST`: indirizzo del server keycloak
  - `KC_HOST_PORT`: porta di ascolto del server keycloak
  - `KC_ADMIN`: nome dell'user temporaneo admin di keycloak
  - `KC_ADMIN_PASSWORD`: password dell'user temporaneo admin di keycloak
  - `KC_THEME`: schermata di login custom (in questo inserire `custom-login-theme`)

Per avviare il server `go`, con un relativo database locale, invece, creare un file `.env` all'interno della cartella "base" con le seguenti variabili d'ambiente:
  - `GOOSE_DBSTRING`: stringa per connessione al database
  - `GOOSE_DRIVER`: database driver
  - `GOOSE_MIGRATION_DIR`: cartella contenente le migrazioni dei database (`./internal/infrastructure/adapters/postgresql/migrations`)
  - `JWT_PUBLIC_KEY`: chiave pubblica del server keycloak per validare i `jwt`
  - `DEV_DB_HOST`: usare il nome del servizio definito in `./docker-compose-db.yml` (postgis)
  - `DEV_DB_PORT`: la porta sul quale si vuole esporre il database
  - `DEV_DB`: nome del database locale
  - `DEV_DB_USER`: username per il database locale
  - `DEV_DB_PASSWORD`: password per il database locale

Creare una rete docker per far comunicare i container dei database e dell'api:
```sh
# Check if network exists
docker network inspect reccelog-shared-network

# If return value is an empty list, create
docker network create reccelog-shared-network
```

Crea l'immagine docker per il server `go`:
```sh
docker compose -f docker-compose-api.yml build
```

Ora avviare i container con:
```sh
# Database container
docker compose -f docker-compose-db.yml up -d
# After starting database, you should run migrations
goose up

# Auth container
docker compose -f Auth/docker-compose.yml up -d

# Api container
docker compose -f docker-compose-api.yml up -d
# Or if you want to build and run
docker compose -f docker-compose-api.yml up -d --build
```