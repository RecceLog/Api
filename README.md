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

# Esempi di richieste
```http request
# INSERIRE STRADA: Panoramica
POST https://reccelogapi.onrender.com/v1/routes/
Content-Type: application/json

{
    "route": {
        "start": {
            "lat": 45.044153,
            "lng": 7.766673
        },
        "finish": {
            "lat": 45.074064,
            "lng": 7.777123
        }
    },
    "note_set": {
        "notes": [
            {
                "position": {
                    "lat": 45.045115,
                    "lng": 7.767779
                },
                "type": "indication",
                "severity": "3",
                "direction":"left",
                "description": "sinistra tre"
            },
            {
                "position": {
                    "lat": 45.045957,
                    "lng": 7.767379
                },
                "type": "indication",
                "severity": "3",
                "direction":"right",
                "description": "destra tre chiude"
            }
        ]
    }
}

###
# INSERIRE STRADA: Strada comunale superga
POST https://reccelogapi.onrender.com/v1/routes/
Content-Type: application/json

{
  "route": {
    "start": {
      "lat": 45.081807,
      "lng": 7.732988
    },
    "finish": {
      "lat": 45.077277,
      "lng": 7.761719
    }
  },
  "note_set": {
    "notes": [
      {
        "position": {
          "lat": 45.045115,
          "lng": 7.767779
        },
        "type": "indication",
        "severity": "3",
        "direction":"left",
        "description": "sinistra tre"
      },
      {
        "position": {
          "lat": 45.045957,
          "lng": 7.767379
        },
        "type": "indication",
        "severity": "3",
        "direction":"right",
        "description": "destra tre chiude"
      }
    ]
  }
}

###
# INSERIRE STRADA: Via cave (rorà)
POST https://reccelogapi.onrender.com/v1/routes/
Content-Type: application/json

{
    "route": {
        "start": {
            "lat": 44.805123,
            "lng": 7.242619
        },
        "finish": {
            "lat": 44.791587,
            "lng": 7.200493
        }
    },
    "note_set": {
        "notes": [
            {
                "position": {
                    "lat": 44.804893,
                    "lng": 7.241788
                },
                "type": "indication",
                "severity": "6",
                "direction":"left",
                "description": "sinistra sei"
            },
            {
                "position": {
                    "lat": 44.804576,
                    "lng": 7.241054
                },
                "type": "indication",
                "severity": "5",
                "direction":"right",
                "description": "destra cinque"
            },
            {
                "position": {
                    "lat": 44.804391,
                    "lng": 7.240395
                },
                "type": "indication",
                "severity": "4",
                "direction":"left",
                "description": "sinistra quattro"
            },
            {
                "position": {
                    "lat": 44.804053,
                    "lng": 7.239881
                },
                "type": "indication",
                "severity": "5",
                "direction":"right",
                "description": "destra cinque"
            },
            {
                "position": {
                    "lat": 44.803604,
                    "lng": 7.238797
                },
                "type": "indication",
                "severity": "3",
                "direction":"right",
                "description": "destra tre"
            },
            {
                "position": {
                    "lat": 44.803744,
                    "lng": 7.23664
                },
                "type": "indication",
                "severity": "6",
                "direction":"left",
                "description": "sinistra sei"
            },
            {
                "position": {
                    "lat": 44.8037284,
                    "lng": 7.234777
                },
                "type": "indication",
                "severity": "6",
                "direction":"left",
                "description": "sinistra sei"
            },
            {
                "position": {
                    "lat": 44.803007,
                    "lng": 7.232098
                },
                "type": "indication",
                "severity": "7",
                "direction":"straight",
                "description": "rettilineo 100"
            }
        ]
    }
}

###
# INSERIRE STRADA: Strada secondaria rorà
POST https://reccelogapi.onrender.com/v1/routes/
Content-Type: application/json

{
    "route": {
        "start": {
            "lat": 44.806422,
            "lng": 7.287241
        },
        "waypoints": [
            {
                "position": {
                    "lat": 44.806649,
                    "lng": 7.272475
                }
            },
            {
                "position": {
                    "lat": 44.804868,
                    "lng": 7.2558
                }
            }
        ],
        "finish": {
            "lat": 44.791587,
            "lng": 7.200493
        }
    },
    "note_set": {
        "notes": [
            {
                "position": {
                    "lat": 44.806329,
                    "lng": 7.282342
                },
                "type": "indication",
                "severity": "4",
                "direction":"right",
                "description": "test"
            },
            {
                "position": {
                    "lat": 44.80648,
                    "lng": 7.281546
                },
                "type": "indication",
                "severity": "5",
                "direction":"left",
                "description": "test"
            }
        ]
    }
}

###
# OTTENERE TUTTE LE STRADE
GET https://reccelogapi.onrender.com/v1/routes

###
# OTTENERE DATI DI STRADA SPECIFICA CON TUTTI I NOTE SET
GET https://reccelogapi.onrender.com/v1/routes/019b85eb-2bbc-7bf0-a08b-d4b1f811750d

###
# OTTENERE TUTTE LE STRADE IN UN RANGE
GET https://reccelogapi.onrender.com/v1/routes/range/10000
Latitude: 45.013902
Longitude: 7.659012

###
# OTTENERE NOTE DI UN SET DI UNA STRADA
GET http://localhost:8080/v1/routes/019aeac6-4dea-71f7-86b6-e05e67ce5167/note-set/019aeac5-4df9-7abc-9d29-a09c9485a5ff

###
# MODIFICA NOTA DI UN SET DI UNA STRADA
PATCH http://localhost:8080/v1/notes/019aeac6-4dea-71f7-86b6-e05e67ce5167/note-set/019aeac6-4df9-7abc-9d29-a09c9485a5ff/note/

###
# AGGIUNGERE SET DI NOTE AD UNA STRADA
POST http://localhost:8080/v1/routes/019aeac6-4dea-71f7-86b6-e05e67ce5167/notes
Content-Type: application/json

{
    "notes": [
        {
            "position": {
                "lat": 44.884646,
                "lng": 7.369042
            },
            "type": "warning",
            "severity": "",
            "direction": "left",
            "description": "autovelox"
        },
        {
            "position": {
                "lat": 44.878532,
                "lng": 7.355322
            },
            "type": "warning",
            "severity": "",
            "direction": "left",
            "description": "altra roba"
        }
    ]
}

###
# CANCELLA SET DI NOTE DI UNA STRADA
DELETE http://localhost:8080/v1/routes/019aeac6-4dea-71f7-86b6-e05e67ce5167/note-set/019aeb3b-ba19-7ebc-8ea1-609050485c48
```