package main

import (
	"fmt"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
)

var riscErr = errors.New("risky operation failed")

// サンプル関数：エラーを返す可能性がある
func riskyOperation() error {
	// 何らかの理由でエラーを生成する
	return errors.Wrap(stack1(), "riskyOperation failed")
}

func stack1() error {
	return errors.Wrap(stack2(), "stack1")
}

func stack2() error {
	return riscErr
}

type MyError struct{}

func (e *MyError) Error() string {
	return "myError"
}

func riskyOperation2() error {
	return &MyError{}
}

func riskyOperation3() error {
	err := errors.New("error happen!")
	err = errors.WithMessage(err, "message")
	err = errors.WithMessage(err, "message2")
	err = errors.WithHint(err, "hint")
	err = errors.WithHint(err, "hint2")
	err = errors.WithDetail(err, "detail")
	err = errors.WithDetail(err, "detail2")
	err = errors.WithStack(err)
	return err
}

func main() {
	// riskyOperation を実行し、エラーをチェック
	err := riskyOperation()
	if err != nil {
		// エラーがある場合、それをラップして追加情報を提供
		wrappedErr := errors.Wrap(err, "main で riskyOperation が失敗")

		// ラップされたエラーを検査
		if errors.Is(wrappedErr, riscErr) {
			fmt.Println(wrappedErr)
		}

		logger, _ := zap.NewProduction()
		defer logger.Sync()
		loggerDev ,_ := zap.NewDevelopment()
		defer loggerDev.Sync()

		logger.Error(
			fmt.Errorf("%+v", wrappedErr).Error(),
		)

		fmt.Println("")
		loggerDev.Error(
			fmt.Errorf("%+v", wrappedErr).Error(),
		)

		// fmt.Printf("%+v\n", wrappedErr)
	} else {
		// エラーがない場合の処理
		fmt.Println("riskyOperation was successful")
	}

	fmt.Print("\n\n\n")

	var target *MyError
	err = riskyOperation2()
	if err != nil {
		if errors.As(err, &target) {
			fmt.Println("same my error")
		}
	}

	fmt.Print("\n\n\n")

	err = riskyOperation3()
	if err != nil {
		fmt.Println(errors.GetAllHints(err))
		fmt.Println(errors.GetAllDetails(err))
		fmt.Println(errors.FlattenHints(err))
		fmt.Println(errors.FlattenDetails(err))
	}
}
