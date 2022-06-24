package job

import (
	"context"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"go-taskfile-example/delay"
	"go-taskfile-example/pb"
	"log"
	"math/rand"
	"time"
)

const MaxAccounts = 1000000

func Run(svc pb.AccountServiceServer) {
	ctx := context.Background()

	for i := 0; i < MaxAccounts; i++ {

		years := rand.Intn(90) + 10
		month := rand.Intn(12) + 1
		days := rand.Intn(31)
		birthday := time.Now().AddDate(-years, month, -days)

		account, err := svc.CreateAccount(ctx, &pb.CreateAccountRequest{
			Name:     faker.Name(),
			Email:    faker.Email(),
			Password: faker.Password(),
			Birthday: birthday.Format("2006-01-02"),
			Gender:   faker.Gender(),
		})
		if err != nil {
			log.Println("error on create account: ", err.Error())
		} else {
			fmt.Printf("%+v\n", account)
		}

		delay.Random(10)
	}
}
