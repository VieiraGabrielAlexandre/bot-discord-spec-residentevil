package main

import "github.com/bwmarrin/discordgo"

func NewReCommand(index BioIndex) *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "re",
		Description: "Bio rápida de personagens Resident Evil",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "character",
				Description: "Nome: " + index.List(),
				Required:    true,
			},
		},
	}
}
