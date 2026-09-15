package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/rphmauriciodev/myFinances-consumer/internal/dto"
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

			select {
			case <-time.After(5 * time.Second):
			case <-ctx.Done():
				return
			}
			continue
		}

		if len(messages) == 0 {
			slog.Info("Nenhuma mensagem recebida da fila SQS, aguardando por novas mensagens...")
		}

		for _, message := range messages {
			if message.Body == nil || message.ReceiptHandle == nil {
				slog.Error("Mensagem inválida recebida da fila SQS", "mensagem", message)
				continue
			}
			slog.Info("Mensagem recebida da fila SQS", "mensagem", message.Body)

			transactionDTO, err := dto.ParseTransaction(*message.Body)
			if err != nil {
				slog.Error("Erro ao parsear a mensagem para DTO Transaction", "mensagem", message.Body)
				continue
			}

			if err := processor.ProcessTransactions(ctx, *transactionDTO); err != nil {
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
