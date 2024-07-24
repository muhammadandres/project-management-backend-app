package helper

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

func SetupSES(Recipients []string, Subject string, TextBody string) error {
	sender := "m.andres.novrizal@gmail.com"

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	client := sesv2.NewFromConfig(cfg)

	input := &sesv2.SendEmailInput{
		FromEmailAddress: &sender,
		Destination: &types.Destination{
			ToAddresses: Recipients,
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Text: &types.Content{
						Data: &TextBody,
					},
				},
				Subject: &types.Content{
					Data: &Subject,
				},
			},
		},
	}

	response, err := client.SendEmail(context.TODO(), input)
	if err != nil {
		return err
	}

	log.Println(response.MessageId)

	return nil
}
