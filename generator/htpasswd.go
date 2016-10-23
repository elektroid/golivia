package generator


import(
	"github.com/foomo/htpasswd"
	"github.com/gin-gonic/gin"
	"fmt"
)

// unused, meant for static generation (but going with a live server for now)
func GenerateHtaccess(accounts gin.Accounts, dir string) error{
	hashes := htpasswd.HashedPasswords(make(map[string]string))

	for user, password := range accounts {
			err := hashes.SetPassword(user, password, htpasswd.HashBCrypt)
			if err != nil {
				return err
		}
	}

	err := hashes.WriteToFile(fmt.Sprintf("%s/.htpasswd", dir))
	if err!=nil{
		return err
	}

	return nil
}