# Numberblocks (YouTube) — Stremio addon

Static Stremio stream addon: maps Numberblocks episodes (Cinemeta `tt7978538`)
to videos on the official Numberblocks YouTube channel (`ytId` streams).

- `data/episodes.json` — `"season:episode"` → YouTube video id
- `go run ./cmd/gen` — regenerates the addon into `docs/` (served by GitHub Pages)

Install in Stremio: `https://<user>.github.io/numberblocks_stremio/manifest.json`
