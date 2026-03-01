package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type CharBio struct {
	Name   string
	Role   string
	Bio    string
	Traits string
}

var bios = map[string]CharBio{
	"leon": {
		Name:   "Leon S. Kennedy",
		Role:   "Agente / Sobrevivente",
		Bio:    "Entrou em serviço no pior dia possível e acabou virando referência em operações de bio-terror.",
		Traits: "Calmo sob pressão, protetor, improvisa bem.",
	},
	"jill": {
		Name:   "Jill Valentine",
		Role:   "Operadora tática",
		Bio:    "Veterana de incidentes biológicos, conhecida por raciocínio rápido e resiliência absurda.",
		Traits: "Tática, resiliente, foco em missão.",
	},
	"claire": {
		Name:   "Claire Redfield",
		Role:   "Civil / Sobrevivente",
		Bio:    "Teimosamente corajosa, faz o impossível para proteger inocentes e encontrar respostas.",
		Traits: "Empática, corajosa, não abandona ninguém.",
	},
	"chris": {
		Name:   "Chris Redfield",
		Role:   "Líder operacional",
		Bio:    "Obstinado e disciplinado, vive para impedir novas tragédias — mesmo que custe tudo.",
		Traits: "Determinado, forte, liderança direta.",
	},
	"ada": {
		Name:   "Ada Wong",
		Role:   "Agente independente",
		Bio:    "Misteriosa e sempre com um plano B (e C). Ajuda quando convém — e some quando quer.",
		Traits: "Enigmática, estratégica, imprevisível.",
	},
	"wesker": {
		Name:   "Albert Wesker",
		Role:   "Antagonista / Estrategista",
		Bio:    "Frio e calculista, manipula pessoas e eventos como peças em um tabuleiro.",
		Traits: "Calculista, ambicioso, controlador.",
	},
	"hunk": {
		Name:   "HUNK",
		Role:   "Operador",
		Bio:    "O tipo que entra, pega o alvo e sai — e quase sempre é o único que volta.",
		Traits: "Silencioso, preciso, missão acima de tudo.",
	},
	"ethan": {
		Name:   "Ethan Winters",
		Role:   "Sobrevivente",
		Bio:    "Um cara comum jogado no horror. Persistência é a arma principal.",
		Traits: "Persistente, protetor, improvável.",
	},
	"nemesis": {
		Name:   "Nemesis",
		Role:   "Arma biológica",
		Bio:    "Perseguição implacável focada em alvo. Não negocia. Não cansa.",
		Traits: "Implacável, brutal, objetivo único.",
	},
}

func main() {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN não definido no ambiente")
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("erro criando sessão: %v", err)
	}

	// Para slash commands, o intent de Guilds é suficiente.
	dg.Identify.Intents = discordgo.IntentsGuilds

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("logado como %s#%s", r.User.Username, r.User.Discriminator)
	})

	dg.AddHandler(onInteractionCreate)

	err = dg.Open()
	if err != nil {
		log.Fatalf("erro abrindo conexão: %v", err)
	}
	defer dg.Close()

	// Registra comando global: /re character:<texto>
	// Observação: comandos globais podem demorar um pouco pra aparecer. Em servidor específico é mais rápido.
	cmd := &discordgo.ApplicationCommand{
		Name:        "re",
		Description: "Bio rápida de personagens (genérico por enquanto)",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "character",
				Description: "Nome: leon, jill, claire, chris, ada, wesker, hunk, ethan, nemesis",
				Required:    true,
			},
		},
	}

	appID := dg.State.User.ID
	_, err = dg.ApplicationCommandCreate(appID, "", cmd)
	if err != nil {
		log.Fatalf("erro registrando comando: %v", err)
	}
	log.Println("comando /re registrado (global).")

	log.Println("rodando. CTRL+C para sair.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop
}

func onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	query := strings.ToLower(strings.TrimSpace(opts[0].StringValue()))
	b, ok := bios[query]
	if !ok {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Não achei esse personagem. Tenta: leon, jill, claire, chris, ada, wesker, hunk, ethan, nemesis",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       b.Name,
		Description: b.Bio,
		Color:       0x8B0000, // vermelho escuro estilo Resident Evil
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Tipo",
				Value:  b.Role,
				Inline: true,
			},
			{
				Name:   "Traços",
				Value:  b.Traits,
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Umbrella Archives • Experimental Bot",
		},
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
