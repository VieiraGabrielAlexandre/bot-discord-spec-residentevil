package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func makeInteractionHandler(index BioIndex) func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}
		if i.ApplicationCommandData().Name != "re" {
			return
		}

		opts := i.ApplicationCommandData().Options
		if len(opts) == 0 {
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{Content: "Use: /re character:<nome>"},
			})
			return
		}

		query := opts[0].StringValue()
		bio, ok := index.Find(query)
		if !ok {
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Não achei esse personagem. Tenta: %s", index.List()),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

		embed := buildBioEmbed(bio)
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
	}
}

func buildBioEmbed(bio CharBio) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       bio.Name,
		Description: bio.Bio,
		Color:       0x8B0000,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Tipo",
				Value:  bio.Role,
				Inline: true,
			},
			{
				Name:   "Traços",
				Value:  bio.Traits,
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Umbrella Archives • Experimental Bot",
		},
	}
}
