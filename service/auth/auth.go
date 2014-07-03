package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/emicklei/go-restful"
	log "github.com/featen/ags/utils/log"
	"github.com/gorilla/sessions"
)

const (
	AdminId = "1000"
)

var sysMagicNumber []byte
var CookieStore *sessions.CookieStore

func SetSysMagicNumber(m []byte) {
	sysMagicNumber = m
	CookieStore = sessions.NewCookieStore(sysMagicNumber)
}

func AuthFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	b, userid := authHandler(req.Request, resp.ResponseWriter)
	if !b {
		log.Debug("unauthorized request %s %s", req.Request.Method, req.Request.URL)
		resp.WriteErrorString(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	req.SetAttribute("agsuserid", userid)
	chain.ProcessFilter(req, resp)
}

func AddCookie(req *http.Request, resp http.ResponseWriter) {
	s, _ := CookieStore.Get(req, "ags-session")
	s.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 0,
	}
	t := time.Now().String()
	s.Values["id"] = AdminId
	s.Values["time"] = t
	s.Values["magic"] = genMagic(AdminId, t)
	s.Save(req, resp)
}

func DelCookie(req *restful.Request, resp *restful.Response) {
	s, _ := CookieStore.Get(req.Request, "ags-session")

	s.Values["id"] = ""
	s.Values["time"] = ""
	s.Values["magic"] = ""
	s.Save(req.Request, resp.ResponseWriter)
}

func authHandler(r *http.Request, w http.ResponseWriter) (bool, string) {
	s, err := CookieStore.Get(r, "ags-session")
	if err != nil {
		log.Debug("Cannot get session: %v", err)
		return false, ""
	}

	if s.Values["id"] == nil || s.Values["id"].(string) != AdminId || s.Values["time"] == nil || s.Values["magic"] == nil {
		return false, ""
	}

	b := check(s.Values["id"].(string), s.Values["time"].(string), s.Values["magic"].(string))
	if b == true {
		return true, s.Values["id"].(string)
	} else {
		return false, ""
	}
}

func check(id string, n string, magic string) bool {
	if m := genMagic(id, n); m == magic {
		return true
	} else {
		return false
	}
}

func genMagic(id string, n string) string {
	h := md5.New()
	io.WriteString(h, "ags-")
	io.WriteString(h, id)
	io.WriteString(h, "-"+n)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Encode(p string) string {
	key := []byte(sysMagicNumber)
	plaintext := []byte(p)

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Debug("%v", err)
		return ""
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Debug("%v", err)
		return ""
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	return fmt.Sprintf("%x\n", ciphertext)
}

func Decode(c string) string {
	key := []byte(sysMagicNumber)
	ciphertext, _ := hex.DecodeString(c)

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Debug("%v", err)
		return ""
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		log.Debug("ciphertext too short")
		return ""
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s\n", ciphertext)
}
