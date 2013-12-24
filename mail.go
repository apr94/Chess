            package main


            import (

              "crypto/tls"

              "fmt"

              "log"

              "net"

              "net/mail"

              "net/smtp"

            )


            func sendMail(rec string, subj string, body string) {

              // the basics

              from := mail.Address{"Go Chess Server", "Email@gmail.com"}

              to := mail.Address{"Player 2", rec}


              // setup the remote smtpserver & auth info

              smtpserver := "smtp.gmail.com:25"

              auth := smtp.PlainAuth("", "Email@gmail.com", "password", "smtp.gmail.com")


              // setup a map for the headers

              header := make(map[string]string)

              header["From"] = from.String()

              header["To"] = to.String()

              header["Subject"] = subj


              // setup the message

              message := ""

              for k, v := range header {

                message += fmt.Sprintf("%s: %s\r\n", k, v)

              }

              message += "\r\n" + body


              // create the smtp connection

              c, err := smtp.Dial(smtpserver)

              if err != nil {

                log.Panic(err)

              }


              // set some TLS options, so we can make sure a non-verified cert won't stop us sending

              host, _, _ := net.SplitHostPort(smtpserver)

              tlc := &tls.Config{

                InsecureSkipVerify: true,

                ServerName:         host,

              }

              if err = c.StartTLS(tlc); err != nil {

                log.Panic(err)

              }


              // auth stuff

              if err = c.Auth(auth); err != nil {

                log.Panic(err)

              }


              // To && From

              if err = c.Mail(from.Address); err != nil {

                log.Panic(err)

              }

              if err = c.Rcpt(to.Address); err != nil {

                log.Panic(err)

              }


              // Data

              w, err := c.Data()

              if err != nil {

                log.Panic(err)

              }

              _, err = w.Write([]byte(message))

              if err != nil {

                log.Panic(err)

              }

              err = w.Close()

              if err != nil {

                log.Panic(err)

              }

              c.Quit()

        }


