package loggergo

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
	"strconv"
	"sync"
	"time"
)

const (
	loggerWriterContextKey = "com.ditotechnologies.loggergo.context"
)

type loggerWriter struct {
	endpointDefinitionMutex sync.Mutex
	service                 *timestreamwrite.TimestreamWrite
	projectID               string
}

func (lw *loggerWriter) log(ctx context.Context, eventName string, t time.Time) {
	records := make([]*timestreamwrite.Record, 0)

	record := &timestreamwrite.Record{
		Time:     aws.String(strconv.FormatInt(t.UnixMilli(), 10)),
		TimeUnit: aws.String("MILLISECONDS"),
	}
	records = append(records, record)

	lw.sendRecords(ctx, records)

}

func (lw *loggerWriter) sendRecords(ctx context.Context, records []*timestreamwrite.Record) {
	input := &timestreamwrite.WriteRecordsInput{
		DatabaseName: aws.String("parsley_analytics"),
		TableName:    aws.String(lw.projectID),
		Records:      records,
	}

	output, err := lw.service.WriteRecords(input)
	if err != nil {
		// TODO (dito) should print out the error here
	}
	fmt.Println(output)
}
