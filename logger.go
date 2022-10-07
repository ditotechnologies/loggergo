package loggergo

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
	"time"
)

//var writer = &loggerWriter{}

func WarnErrorWithContext(ctx context.Context, e error) {
	LogWithContext(ctx, fmt.Sprintf("err_%s", e.Error()))
}

func LogWithContext(ctx context.Context, eventName string) {
	// TODO (dito) make this fault tolerant by trying to write to the file system
	t := time.Now()
	go func() {
		loggerWriterInterface := ctx.Value(loggerWriterContextKey)
		if loggerWriterInterface == nil {
			return
		}
		writer, ok := loggerWriterInterface.(*loggerWriter)
		if !ok {
			panic("Could not load the writer. Please reach out to Dito Technologies LLC technical support")
		}
		writer.log(ctx, eventName, t)
	}()
}

func SetupWithContext(ctx context.Context, projectID string, accessKeyID string, accessKeySecret string) context.Context {
	s, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(accessKeyID, accessKeySecret, ""),
	})
	if err != nil {
		panic("could not create loggergo values session, this is unrecoverable")
	}
	service := timestreamwrite.New(s)
	writer := &loggerWriter{
		service:   service,
		projectID: projectID,
	}
	return context.WithValue(ctx, loggerWriterContextKey, writer)
}
