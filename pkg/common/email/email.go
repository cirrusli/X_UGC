package email

import (
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	"math/rand"
	"net/smtp"
	"time"
)

func SendCode(userEmail string) (string, error) {
	//生成六位数验证码
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06v", r.Int31n(1000000))
	em := email.NewEmail()
	// 设置 sender 发送方 的邮箱 ， 此处可以填写自己的邮箱
	em.From = "2963703604@qq.com"

	// 设置 receiver 接收方 的邮箱  此处也可以填写自己的邮箱， 就是自己发邮件给自己
	em.To = []string{userEmail}

	// 设置主题
	em.Subject = "验证码到啦"

	//发送的内容
	html := fmt.Sprintf(`<div>
        <div>
            尊敬的用户，您好！
        </div>
        <div style="padding: 8px 40px 8px 50px;">
            <p>本次的验证码为%s,为保证账号安全，验证码有效期为1分钟。请确认为本人操作，切勿向他人泄露，感谢您的理解与使用。</p>
        </div>
        <div>
            <p>此邮件为系统发送，请勿回复。</p>
        </div>
    </div>`, code)
	em.HTML = []byte(html)
	//设置服务器相关的配置
	err := em.SendWithTLS("smtp.qq.com:465",
		smtp.PlainAuth("", "2963703604@qq.com", "hidden", "smtp.qq.com"),
		&tls.Config{
			InsecureSkipVerify: true,
			ServerName:         "smtp.qq.com",
		})
	if err != nil {
		return "", err
	}
	return code, nil
}
