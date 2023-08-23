package goopentelemetry

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/cast"
	"log"
	"os"
	"runtime"
	"strings"
)

func GetEnvOrDefault(key string, defaultValue interface{}) interface{} {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func GetEnv(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Println("Cannot load file .env: ", err)
		panic(err)
	}

	value := GetEnvOrDefault(key, "").(string)
	return value
}

func StringToBool(value string) bool {
	if value == "true" {
		return true
	}

	return false
}

func GetActionName() string {
	c, _, _, _ := runtime.Caller(1)
	f := runtime.FuncForPC(c).Name()
	fs := strings.SplitN(f, ".", 2)
	replacer := strings.NewReplacer("(", "", ")", "", "*", "")
	actionName := replacer.Replace(fs[1])

	return actionName
}

func WriteStringTemplate(stringTemplate string, args ...interface{}) string {
	return fmt.Sprintf(stringTemplate, args...)
}

func AnyToBool(value any) bool {
	return cast.ToBool(value)
}

func GetFunctionName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	details := runtime.FuncForPC(pc)
	if details == nil {
		return ""
	}

	funcName := details.Name()
	lastDot := strings.LastIndex(funcName, ".")
	if lastDot != -1 {
		funcName = funcName[lastDot+1:]
	}
	return funcName
}
