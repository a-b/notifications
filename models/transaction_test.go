package models_test

import (
    "github.com/cloudfoundry-incubator/notifications/models"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Transaction", func() {
    var transaction models.TransactionInterface

    BeforeEach(func() {
        TruncateTables()
        transaction = models.Database().Connection().Transaction()
    })

    Describe("Begin/Commit", func() {
        It("commits the transaction to the database", func() {
            err := transaction.Begin()
            if err != nil {
                panic(err)
            }

            repo := models.NewClientsRepo()
            _, err = repo.Create(transaction, models.Client{
                ID:          "my-client",
                Description: "My Client",
            })
            if err != nil {
                panic(err)
            }

            err = transaction.Commit()
            if err != nil {
                panic(err)
            }

            client, err := repo.Find(models.Database().Connection(), "my-client")
            if err != nil {
                panic(err)
            }

            Expect(client.ID).To(Equal("my-client"))
            Expect(client.Description).To(Equal("My Client"))
        })
    })

    Describe("Begin/Rollback", func() {
        It("rolls back the transaction from the database", func() {
            err := transaction.Begin()
            if err != nil {
                panic(err)
            }

            repo := models.NewClientsRepo()
            _, err = repo.Create(transaction, models.Client{
                ID:          "my-client",
                Description: "My Client",
            })
            if err != nil {
                panic(err)
            }

            err = transaction.Rollback()
            if err != nil {
                panic(err)
            }

            _, err = repo.Find(models.Database().Connection(), "my-client")
            Expect(err).To(BeAssignableToTypeOf(models.ErrRecordNotFound{}))
        })
    })
})
