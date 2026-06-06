package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

const (
	smtpServer   = "smtp.mail.ru"
	smtpPort     = 465
	smtpUser     = "tenders_info@internet.ru"
	smtpPassword = "DVEgz21kg0zDLUE9cc1q"
	fromEmail    = "tenders_info@internet.ru"
	webHost      = "http://localhost:8080" //Uri web app(domain)
)

func SendResetEmail(email, token string) error {
	resetLink := fmt.Sprintf("%s/main/password/reset/form?token=%s", webHost, token)
	subject := "Восстановление пароля"
	body := fmt.Sprintf("Для восстановления пароля перейдите по ссылке ниже, \nесли вы не запрашивали восстановление пароля, \nПРОИГНОРИРУЙТЕ данное письмо:\n %s \nС Уважением Администрация портала коммерческих закупок!", resetLink)

	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpServer) 

	// Создаем TLS-конфигурацию
	tlsConfig := &tls.Config{
		ServerName:         smtpServer, 
		InsecureSkipVerify: false,      
	}

	// Подключаемся к SMTP-серверу с TLS
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", smtpServer, smtpPort), tlsConfig) 
	if err != nil {
		return fmt.Errorf("Ошибка при подключении к SMTP-серверу: %v", err)
	}
	defer conn.Close()

	// Создаем SMTP-клиент
	c, err := smtp.NewClient(conn, smtpServer) 
	if err != nil {
		return fmt.Errorf("Ошибка при создании SMTP-клиента: %v", err)
	}
	defer c.Quit()

	if err = c.Auth(auth); err != nil {
		return fmt.Errorf("Ошибка аутентификации: %v", err)
	}

	if err = c.Mail(fromEmail); err != nil {
		return fmt.Errorf("Ошибка при установке отправителя: %v", err)
	}
	if err = c.Rcpt(email); err != nil {
		return fmt.Errorf("Ошибка при установке получателя: %v", err)
	}

	// Отправляем данные
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("Ошибка при получении потока данных: %v", err)
	}
	defer w.Close()

	msg := []byte("To: " + email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("Ошибка при записи данных: %v", err)
	}

	return nil
}
