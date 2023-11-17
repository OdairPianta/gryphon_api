package notifications

import (
	"fmt"
	"os"

	"github.com/OdairPianta/gryphon_api/models"
)

func ForgotPassword(user models.User) error {
	var appName = os.Getenv("APP_NAME")
	var passwordRecoveyUrl = os.Getenv("APP_URL") + "/password_recovery/" + fmt.Sprint(user.ResetPasswordToken)
	var message string = "Você esta recebendo este e-mail porque solicitou a recuperação de senha da sua conta.\n\n" +
		"Por favor, clique no link abaixo ou cole o link em seu navegador para completar o processo:\n\n" +
		passwordRecoveyUrl + "\n\n" +
		"Se você não solicitou a recuperação de senha, ignore este e-mail e sua senha permanecerá inalterada.\n" +
		"\n\n" +
		"Obrigado,\n" +
		"Equipe " + appName

	return ToMail(user.Email, "Recuperação de senha", message)

}
