package emailhelper

import (
	"amartha/loan-service/helpers/envhelper"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

func CreateClientFromSession() (*ses.SES, error) {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(envhelper.GetEnvVar("AWS_DEPLOYMENT_REGION")),
	})
	if err != nil {
		err = fmt.Errorf("error creating SES client: %s", err.Error())
		return nil, err
	}

	return ses.New(session), nil
}

const DefaultEmailSender = "diskon.hunter.official@gmail.com"

type SendEmailArgs struct {
	Sender        string
	RecipientList []*string //slice of aws.String
	CcList        []*string //slice of aws.String
	Subject       string
	HtmlBody      string
	TextBody      string
}

func SendEmail(sesClient *ses.SES, args SendEmailArgs) (err error) {
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: args.CcList,
			ToAddresses: args.RecipientList,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(args.HtmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(args.TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(args.Subject),
			},
		},
		Source: aws.String(args.Sender),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	_, err = sesClient.SendEmail(input)
	if err != nil {
		err = fmt.Errorf("error sending email %s: %v", args.Subject, err)
		return err
	}

	return nil
}
