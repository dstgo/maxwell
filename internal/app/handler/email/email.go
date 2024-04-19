package email

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/dstgo/maxwell/internal/app/conf"
	"github.com/dstgo/maxwell/internal/app/data/mq"
	"github.com/ginx-contribs/ginx/pkg/resp/statuserr"
	"github.com/ginx-contribs/str2bytes"
	"github.com/matcornic/hermes/v2"
	"github.com/redis/go-redis/v9"
	"github.com/wneessen/go-mail"
	"golang.org/x/net/context"
)

// Message represents an email message
type Message struct {
	ContentType mail.ContentType `mapstructure:"contentType"`
	From        string           `mapstructure:"from"`
	To          []string         `mapstructure:"to"`
	CC          []string         `mapstructure:"cc"`
	Bcc         []string         `mapstructure:"bcc"`
	Subject     string           `mapstructure:"subject"`
	Body        string           `mapstructure:"body"`
}

// BuildMailMsg build *mail.Msg from Message
func BuildMailMsg(msg Message) (*mail.Msg, error) {
	mailMsg := mail.NewMsg()
	if err := mailMsg.From(msg.From); err != nil {
		return nil, err
	}
	if err := mailMsg.To(msg.To...); err != nil {
		return nil, err
	}
	if err := mailMsg.Cc(msg.CC...); err != nil {
		return nil, err
	}
	if err := mailMsg.Bcc(msg.Bcc...); err != nil {
		return nil, err
	}
	mailMsg.Subject(msg.Subject)
	mailMsg.SetBodyString(msg.ContentType, msg.Body)
	return mailMsg, nil
}

func NewEmailHandler(cfg conf.EmailConf, client *mail.Client, queue mq.Queue) (*Handler, error) {
	h := &Handler{Cfg: cfg}

	h.pub = &Publisher{queue: queue, topic: cfg.MQ.Topic, maxLen: cfg.MQ.MaxLen}
	// run the consumer
	for _, consumer := range cfg.MQ.Consumers {
		c := &Consumer{
			topic:     cfg.MQ.Topic,
			group:     cfg.MQ.Group,
			consumer:  consumer,
			batchSize: cfg.MQ.BatchSize,
			queue:     queue,
			email:     client,
		}
		if err := c.Consume(context.Background()); err != nil {
			return nil, err
		}
		h.cons = append(h.cons, c)
	}

	h.product = hermes.Hermes{
		Product: hermes.Product{
			Name:      "dstgo",
			Copyright: "Copyright © dstgo",
		},
	}

	return h, nil
}

// Handler is email handler
type Handler struct {
	Cfg     conf.EmailConf
	pub     *Publisher
	cons    []*Consumer
	product hermes.Hermes
}

func (h *Handler) SendHermesEmail(ctx context.Context, subject string, to []string, email hermes.Email) error {
	html, err := h.product.GenerateHTML(email)
	if err != nil {
		return err
	}

	msg := Message{
		ContentType: mail.TypeTextHTML,
		From:        h.Cfg.Username,
		To:          to,
		Subject:     subject,
		Body:        html,
	}

	if err := h.pub.Publish(ctx, msg); err != nil {
		return err
	}
	return nil
}

func (h *Handler) SendEmail(ctx context.Context, msg Message) error {
	msg.From = h.Cfg.Username
	return h.pub.Publish(ctx, msg)
}

// Publisher store email message into queue and then return right now
type Publisher struct {
	queue mq.Queue

	topic  string
	maxLen int64
}

// Publish publishes email message
func (s *Publisher) Publish(ctx context.Context, msg Message) error {
	marshal, err := sonic.Marshal(msg)
	if err != nil {
		return statuserr.InternalError(err)
	}

	_, err = s.queue.Publish(ctx, s.topic, map[string]any{"mail": marshal}, s.maxLen)
	if err != nil {
		return statuserr.InternalError(err)
	}
	return err
}

// Consumer is responsible for reading messages from queues and then send these emails.
type Consumer struct {
	topic     string
	group     string
	consumer  string
	batchSize int64

	queue mq.Queue
	email *mail.Client
}

func (c *Consumer) Consume(ctx context.Context) error {
	return c.queue.Consume(ctx, c.topic, c.group, c.consumer, c.batchSize, c.consume)
}

func (c *Consumer) consume(ctx context.Context, client *redis.Client, topic, group string, id string, value any) error {
	val, ok := value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("mismatched value type from mq, expected map[string]any, but got %T", value)
	}

	var mailMsg string
	if val["mail"] != nil {
		mailMsg = val["mail"].(string)
	}

	var msg Message
	err := sonic.Unmarshal(str2bytes.Str2Bytes(mailMsg), &msg)
	if err != nil {
		return err
	}

	buildMail, err := BuildMailMsg(msg)
	if err != nil {
		return err
	}

	return c.email.DialAndSendWithContext(ctx, buildMail)
}
