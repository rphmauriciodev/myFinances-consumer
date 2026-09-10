package main

import (
	"context"
	"log/slog"

	"github.com/rphmauriciodev/myFinances-consumer/internal/platform"
	"github.com/rphmauriciodev/myFinances-consumer/internal/processing"
	"github.com/rphmauriciodev/myFinances-consumer/internal/queue"
	"github.com/spf13/viper"
)

func main() {

	slog.Info("Iniciando consumidor...")

	v := viper.New()

	v.AutomaticEnv()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := platform.NewPostgresPool()
	if err != nil {
		slog.Error("Erro ao criar pool de conexões com o Postgres", "erro", err)
		return
	}

	repo := platform.NewPostgresRepository(pool)

	processor := processing.NewProcessor(repo)
	if err != nil {
		slog.Error("Erro ao inicializar o processador", "erro", err)
		return
	}

	queue, err := queue.NewQueue(ctx, v.GetString("SQS_QUEUE_URL"))
	if err != nil {
		slog.Error("Erro ao criar a fila SQS", "erro", err)
		return
	}

	slog.Info("Consumidor iniciado com sucesso!")

	for {
		messages, err := queue.ReceiveMessages(ctx, 10)
		if err != nil {
			slog.Error("Erro ao receber mensagens da fila SQS", "erro", err)
			continue
		}

		for _, message := range messages {
			if message.Body == nil || message.ReceiptHandle == nil {
				slog.Error("Mensagem inválida recebida da fila SQS", "mensagem", message)
				continue
			}
			slog.Info("Mensagem recebida da fila SQS", "mensagem", message.Body)

			transaction, err := processing.ParseTransaction(*message.Body)
			if err != nil {
				slog.Error("Erro ao processar a mensagem da fila SQS", "erro", err)
				continue
			}

			if err := processor.ProcessTransactions(ctx, transaction); err != nil {
				slog.Error("Erro ao processar a transação", "erro", err)
				continue
			}

			if err := queue.DeleteMessage(ctx, *message.ReceiptHandle); err != nil {
				slog.Error("Erro ao deletar a mensagem da fila SQS", "erro", err)
				continue
			}

			slog.Info("Mensagem processada e deletada da fila SQS com sucesso")
		}
	}
}
