package main

import (
	"context"
	"fmt"
	"kenshop/pkg/log"
	ktrace "kenshop/pkg/trace"
	proto "kenshop/proto/test"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var pwd, _ = os.Getwd()
var logger1 *otelzap.Logger = log.MustNewOtelLogger(log.WithOutputPaths(fmt.Sprintf("%s/e1.log", pwd)))
var logger2 *otelzap.Logger = log.MustNewOtelLogger(log.WithOutputPaths(fmt.Sprintf("%s/e2.log", pwd)))

func testA(ctx context.Context, tracer trace.Tracer) {
	//span:=trace.SpanFromContext(ctx)
	_, cspan := tracer.Start(ctx, "chird pp")
	ktrace.MessageSent.Event(ctx, 222, &proto.ReqMessage{Req: "222"})
	defer cspan.End()
	time.Sleep(2 * time.Second)
}

func Span() {

	tp := ktrace.MustNewTracer(context.TODO(), ktrace.WithName("ken"))
	t, err := tp.NewTraceProvider("192.168.199.128:4318")
	if err != nil {
		panic(err)
	}
	tracer := t.Tracer("ken-tracer", trace.WithInstrumentationAttributes(attribute.Bool("isken", false)))
	ctx, span := tracer.Start(context.Background(), "pp", trace.WithAttributes(attribute.String("span key", "test")))
	testA(ctx, tracer)

	span.End()
	t.Shutdown(context.Background())
}

func Consume1(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, v := range msgs {
		logger1.Info(string(v.Body))
	}
	return consumer.ConsumeSuccess, nil
}

func Consume2(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, v := range msgs {
		logger2.Info(string(v.Body))
	}
	return consumer.ConsumeSuccess, nil
}

func Consumer(gn string, f func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)) {
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"192.168.199.128:9876"}),
		consumer.WithGroupName(gn),
	)
	if err != nil {
		panic(err)
	}

	err = c.Subscribe("line", consumer.MessageSelector{}, f)
	if err != nil {
		panic(err)
	}
	err = c.Start()
	if err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Hour)
	c.Shutdown()
}

func Producer() {
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.199.128:9876"}))
	if err != nil {
		panic(err)
	}
	if err = p.Start(); err != nil {
		panic(err)
	}
	for i := range 10 {
		msg := fmt.Sprintf("this is rq msg %d", i)
		_, err := p.SendSync(context.Background(), primitive.NewMessage("line", []byte(msg)))
		if err != nil {
			panic(err)
		}
	}

	if err = p.Shutdown(); err != nil {
		panic("关闭 producer 失败")
	}

}

func ToSnakeCase(s string) string {
	// 使用正则表达式匹配大写字母并在前面加上下划线
	re := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	snake := re.ReplaceAllString(s, `${1}_${2}`)
	// 将所有字符转换为小写
	return strings.ToLower(snake)
}

// ToExportedCamelCase 转换为 Go 的导出驼峰命名法
func GoExportedCamelCase(s string) string {
	return goCamelCase(s)
}

// ToUnexportedCamelCase 转换为 Go 的不导出驼峰命名法:首字母小写
func GoUnexportedCamelCase(s string) string {
	s = goCamelCase(s)
	return strings.ToLower(s[:1]) + s[1:]
}

func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// 此函数能将str转化为go的驼峰命名格式
func goCamelCase(str string) string {
	var b []byte
	for i := 0; i < len(str); i++ {
		c := str[i]
		switch {
		case c == '.' && i+1 < len(str) && isASCIILower(str[i+1]):
		case c == '.':
			b = append(b, '_')
		case c == '_' && (i == 0 || str[i-1] == '.'):

			b = append(b, 'X')
		case c == '_' && i+1 < len(str) && isASCIILower(str[i+1]):
		case isASCIIDigit(c):
			b = append(b, c)
		default:
			if isASCIILower(c) {
				c -= 'a' - 'A'
			}
			b = append(b, c)
			for ; i+1 < len(str) && isASCIILower(str[i+1]); i++ {
				b = append(b, str[i+1])
			}
		}
	}
	ss := string(b)
	return strings.ReplaceAll(ss, "_", "")
}
func regis() {
	//fmt.Println(ToSnakeCase("ni_Hao_ma_wqwtq_Z"))
	fmt.Println(GoExportedCamelCase("userName"))
	gin.Should
}

func main() {
	//Span()
	// go Producer()
	// go Consumer("ken-consumer1", Consume1)
	// Consumer("ken-consumer2", Consume2)
	regis()
}
