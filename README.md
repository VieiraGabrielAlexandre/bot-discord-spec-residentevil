# Resident Evil Bio Bot

Bot de Discord para retornar bios rápidas de personagens de Resident Evil via slash command.

## Requisitos

- Go instalado
- Token do Discord

## Configuração

- `DISCORD_TOKEN`: token do bot
- `BIOS_FILE` (opcional): caminho para o arquivo JSON de bios (default: `bios.json`)

## Como rodar

```bash
go run .
```

## Formato do `bios.json`

```json
{
  "leon": {
    "name": "Leon S. Kennedy",
    "role": "Agente / Sobrevivente",
    "bio": "Texto curto da bio.",
    "traits": "Traços principais."
  }
}
```

## Slash command

- `/re character:<nome>`
