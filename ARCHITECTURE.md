# QuadBoard Architecture

Version: 0.1
Status: Draft

## Vision

QuadBoard est un portail d'applications auto-découvert destiné aux environnements
Podman utilisant des Quadlets.

Le principe est simple :

> Déployez une application, elle apparaît automatiquement.

L'objectif est de supprimer toute configuration manuelle du portail.

Le portail n'est pas un outil d'administration de conteneurs. Il est un point
d'entrée pour les utilisateurs et les administrateurs.

---

# Objectifs

- Un seul binaire.
- Aucune base de données.
- Fonctionne sans configuration.
- Découverte automatique.
- Compatible rootless et rootful.
- Léger.
- Extensible.
- Facile à empaqueter (OCI, Nix, RPM, DEB…).

---

# Non-objectifs

QuadBoard ne cherche pas à remplacer :

- Portainer
- Cockpit
- Grafana
- Traefik Dashboard

Il ne gère pas le cycle de vie des conteneurs.

---

# Philosophie

Le projet privilégie :

- les conventions plutôt que la configuration ;
- les métadonnées déjà présentes dans l'infrastructure ;
- les technologies standards ;
- les dépendances limitées.

---

# Architecture

Le cœur du projet repose sur quatre concepts.

## Provider

Un Provider découvre des applications.

Exemples :

- Quadlet
- Docker Compose
- Kubernetes (futur)
- Nomad (futur)

Interface :

```go
type Provider interface {
    Discover(context.Context) ([]Application, error)
}
```

---

## Application

Une application représente un service visible par un utilisateur.

```go
type Application struct {
    ID          string
    Name        string
    Description string

    URL         string
    Icon        string

    Group       string
    Tags        []string

    Status      Status
}
```

---

## Renderer

Le Renderer transforme les Applications en interface.

Exemples :

- HTML
- JSON
- RSS (éventuellement)

---

## Theme

Un thème définit :

- les templates
- les feuilles de style
- les assets
- les paramètres graphiques

Les thèmes ne doivent jamais modifier la logique métier.

---

# Configuration

Ordre de priorité :

CLI

↓

Variables d'environnement

↓

Configuration TOML

↓

Valeurs par défaut

Le fichier de configuration est optionnel.

---

# Découverte

Au démarrage :

```
Quadlets
        │
        ▼
Lecture des labels
        │
        ▼
Construction des Applications
        │
        ▼
Interface Web
```

La découverte est relancée automatiquement lorsqu'un Quadlet est modifié.

---

# Authentification

QuadBoard ne fournit pas d'authentification.

Elle est déléguée au reverse proxy.

Compatibilité prévue :

- Authelia
- OAuth2 Proxy
- Authentik
- Keycloak
- Reverse proxy HTTP classique

---

# Providers

Version 1

- Quadlet

Version future

- Docker Compose
- Kubernetes
- Nomad

---

# Thèmes

Chaque thème est contenu dans un répertoire.

```
theme/
    modern/
        theme.toml
        templates/
        assets/
```

Les thèmes peuvent modifier :

- couleurs
- polices
- disposition
- composants

Ils ne peuvent pas modifier le backend.

---

# API

Version 1

GET /api/apps

GET /api/status

GET /health

---

# Configuration

Sources possibles :

- fichier TOML
- variables d'environnement
- paramètres CLI

---

# Roadmap

## v0.1

- [ ] serveur HTTP
- [ ] Provider Quadlet
- [ ] page HTML

## v0.2

- [ ] recherche
- [ ] thèmes
- [ ] API

## v0.3

- [ ] Authelia
- [ ] statut des services

## v1.0

- [ ] stabilité
- [ ] documentation
- [ ] image OCI

## Further
- [ ] package Nix

---

# Principes de développement

Chaque fonctionnalité doit :

- être testable ;
- être documentée ;
- être découplée du reste ;
- éviter les dépendances inutiles.

---

# Devise

Deploy.

Discover.

Launch.