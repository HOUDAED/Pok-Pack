# Pok-Pack 🎴

Une application web interactive pour gérer votre collection de cartes Pokémon. Créez un compte, collectez des cartes en ouvrant des packs, consultez vos statistiques personnalisées, et attrapez-les tous !

## 🎯 Fonctionnalités

- **🔐 Authentification Sécurisée** : Création de compte et connexion avec bcrypt pour les mots de passe.
- **📦 Ouverture de Packs** : Obtenez des cartes Pokémon aléatoires.
- **🏆 Collection Personnelle** : Consultez toutes les cartes que vous possédez avec leurs détails.
- **🎨 Statistiques Détaillées** : Visualisez des statistiques sur votre collection. (types, rareté, etc.)
- **🔍 Search/add manual** : Recherchez et ajoutez n'importe quel Pokémon (ID 1-1025) à votre collection.
- **👤 Dashboard** : Vue d'ensemble de votre progression et de vos cartes récentes.
- **💾 Persistance des Données** : Toutes vos données sont sauvegardées localement.

## 🛠️ Technologies Utilisées

### Backend
- **Go 1.24.0** - Langage de programmation.
- **SQLite** avec `modernc.org/sqlite` - Base de données locale.
- **PokeAPI** - Pour récupérer les données des Pokémon en temps réel.
- **bcrypt** (`golang.org/x/crypto`) - Sécurisation des mots de passe.

### Frontend
- **HTML/CSS/JavaScript** - Interface utilisateur responsive.
- **CSS Grid/Flexbox** - Mise en page moderne.

## 📁 Structure du Projet

```
Pok-Pack/
├── cmd/
│   └── server/
│       └── main.go                 # Point d'entrée de l'application
├── internal/pok/
│   ├── app_handlers.go            # Handlers pour les routes principales (pack, collection, etc.)
│   ├── app_render.go              # Rendu des templates
│   ├── auth.go                    # Gestion de l'authentification et des sessions
│   ├── cards_store.go             # Requêtes DB pour les cartes
│   ├── dashboard_handlers.go      # Handlers du tableau de bord
│   ├── db.go                      # Initialisation et configuration de la DB
│   ├── middleware.go              # Middlewares (vérification auth, etc.)
│   ├── models.go                  # Structures de données
│   ├── pokeapi.go                 # Client PokéAPI avec retry
│   ├── query.go                   # Requêtes SQL générales
│   ├── session.go                 # Gestion des sessions utilisateur
│   ├── sql_loader.go              # Chargement des définitions SQL
│   ├── types_stats.go             # Calcul des statistiques par type
│   ├── user_store.go              # Requêtes DB pour les utilisateurs
│   └── user.go                    # Logique utilisateur et authentification
├── web/
│   ├── static/                    # Fichiers CSS, JS et assets statiques
│   │   ├── app.css, app.js        # Styles et scripts pour l'app
│   │   ├── dashboard.css, dashboard.js
│   │   ├── script.js, style.css
│   └── templates/                 # Templates HTML
│       ├── app_layout.html        # Layout principal de l'app
│       ├── app_home.html          # Page d'accueil
│       ├── app_pack.html          # Page d'ouverture de pack
│       ├── app_collection.html    # Page de collection
│       ├── app_add.html           # Page d'ajout manuel
│       ├── app_stats.html         # Page des statistiques
│       ├── connexion.html         # Page de connexion
│       └── inscription.html       # Page d'inscription
├── go.mod                         # Dépendances Go
└── pok.db                         # Base de données SQLite (créée au démarrage)
```

## 🚀 Installation et Utilisation

### Prérequis
- **Go 1.24.0** (ou version supérieure)

### Installation

1. **Cloner le projet** (ou télécharger)
```bash
cd Pok-Pack
```

2. **Démarrer l'application**
```bash
go run ./cmd/server/main.go
```

L'application démarre automatiquement :
- ✅ Initialise la base de données SQLite (`pok.db`)
- ✅ Crée les tables nécessaires
- ✅ Vide les sessions précédentes
- ✅ Lance le serveur web

### Accès à l'application

Ouvrez votre navigateur et accédez à : **http://localhost:8080**

## 📋 Flux Utilisateur

### 1. **Inscription**
- Créez un compte avec un pseudo et un email uniques.
- Obtenez 5 cartes aléatoires de démarrage.

### 2. **Connexion**
- Connectez-vous avec votre pseudo/email et mot de passe.

### 3. **Tableau de Bord**
- Consultez vos cartes récentes.
- Voyez le nombre total de cartes dans votre collection.
- Découvrez les statistiques par type.

### 4. **Ouvrir des Packs**
- Cliquez sur "Ouvrir un pack" pour obtenir 5 cartes aléatoires.
- Chaque ouverture génère de nouvelles cartes.

### 5. **Collection**
- Consultez toutes les cartes que vous possédez.
- Voyez les détails : nom, type, image.

### 6. **Ajouter Manuellement**
- Entrez l'ID d'un Pokémon (1-1025).
- La carte est ajoutée directement à votre collection.

### 7. **Statistiques**
- Analysez la répartition de vos cartes par type.
- Visualisez la composition de votre collection.

## 🔧 Routes Disponibles

| Route | Méthode | Description |
|-------|---------|-------------|
| `/` | GET | Page d'accueil (redirect login/dashboard) |
| `/login` | GET/POST | Page de connexion |
| `/register` | GET/POST | Page d'inscription |
| `/dashboard` | GET | Tableau de bord (auth requise) |
| `/pack` | GET | Page d'ouverture de pack (auth requise) |
| `/pack/open` | POST | Ouvrir un pack (auth requise) |
| `/collection` | GET | Voir sa collection (auth requise) |
| `/card/add` | GET/POST | Ajouter une carte (auth requise) |
| `/stats` | GET | Voir les statistiques (auth requise) |
| `/logout` | GET | Se déconnecter |
| `/static/*` | GET | Fichiers statiques (CSS, JS, images) |

## 🗄️ Base de Données

La base de données SQLite contient les tables suivantes :

- **users** : Informations des utilisateurs. (pseudo, email, mot de passe hashé)
- **sessions** : Gestion des sessions actives.
- **cards** : Cartes possédées par chaque utilisateur.
- **card_types** : Types de cartes. (Feu, Eau, Herbe, etc.)

## 🔒 Sécurité

- ✅ Mots de passe hachés avec bcrypt
- ✅ Validation des sessions utilisateur
- ✅ Authentification requise pour toutes les routes protégées
- ✅ Gestion des erreurs avec messages génériques côté client

## 📡 Intégration PokéAPI

L'application récupère les données des Pokémon depuis **PokéAPI** :
- Nom, ID, types, image officielle
- Retry automatique en cas d'erreur de connexion (jusqu'à 3 tentatives)
- Délai de fallback pour éviter les limitations de taux

## 🐛 Dépannage

### La base de données ne s'initialise pas
```bash
# Supprimez la base de données existante et relancez
rm pok.db
go run ./cmd/server/main.go
```

### Erreur de connexion à PokéAPI
- Vérifiez votre connexion internet
- L'application tentera automatiquement une nouvelle tentative
- Les cartes existantes sont toujours accessibles hors ligne

### Sessions expirées
- Les sessions sont vidées au démarrage de l'application
- Reconnectez-vous si nécessaire

## 📝 Développement

### Ajouter une nouvelle route

1. Créez un handler dans le fichier approprié (`internal/pok/*.go`)
2. Protégez-la avec `RequireAuth()` si nécessaire
3. Enregistrez-la dans `main.go` :
```go
http.HandleFunc("/nouvelle-route", pok.NouveauHandler)
```

### Modifier le schéma de la base de données

1. Éditez les requêtes SQL dans `internal/pok/sql_loader.go`
2. Supprimez `pok.db` pour forcer la recréation des tables
3. Relatance l'application

## 📄 Licence

Ce projet est personnel et à usage libre.

---

**Créé par Edvige et Ryan avec ❤️ en Go**