package config

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

func ConnectS3() (*session.Session, error) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile:           "hinagane-user",
		SharedConfigState: session.SharedConfigEnable,
	}))

	return sess, nil
}
