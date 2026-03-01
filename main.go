package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN não definido no ambiente")
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("erro criando sessão: %v", err)
	}

	bios, err := loadBios(getBiosFilePath())
	if err != nil {
		log.Fatalf("erro carregando bios: %v", err)
	}
	index := NewBioIndex(bios)

	// Para slash commands, o intent de Guilds é suficiente.
	dg.Identify.Intents = discordgo.IntentsGuilds

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("logado como %s#%s", r.User.Username, r.User.Discriminator)
	})

	dg.AddHandler(makeInteractionHandler(index))

	err = dg.Open()
	if err != nil {
		log.Fatalf("erro abrindo conexão: %v", err)
	}
	defer dg.Close()

	appID := dg.State.User.ID
	_, err = dg.ApplicationCommandCreate(appID, "", NewReCommand(index))
	if err != nil {
		log.Fatalf("erro registrando comando: %v", err)
	}
	log.Println("comando /re registrado (global).")

	log.Println("rodando. CTRL+C para sair.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop
}
