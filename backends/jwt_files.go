package backends

import (
	"github.com/iegomez/mosquitto-go-auth/backends/files"
	"github.com/iegomez/mosquitto-go-auth/hashing"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type filesJWTChecker struct {
	checker *files.Checker
	options tokenOptions
}

func NewFilesJWTChecker(authOpts map[string]string, logLevel log.Level, hasher hashing.HashComparer, options tokenOptions) (jwtChecker, error) {
	log.SetLevel(logLevel)

	/*	We could ask for a file listing available users with no password, but that gives very little value
		versus just assuming users in the ACL file are valid ones, while general rules apply to any user.
		Thus, padswords file makes no sense for JWT, we only need to check ACLs.
	*/
	aclPath, ok := authOpts["jwt_acl_path"]
	if !ok || aclPath == "" {
		return nil, errors.New("missing acl file path")
	}

	var checker, err = files.NewChecker(authOpts["backends"], "", aclPath, logLevel, hasher)
	if err != nil {
		return nil, err
	}

	return &filesJWTChecker{
		checker: checker,
		options: options,
	}, nil
}

func (o *filesJWTChecker) GetUser(username, token string) (bool, error) {
	tokenUsername, err := getUsernameForToken(o.options, token, o.options.skipUserExpiration)

	if err != nil {
		log.Debugf("jwt local get user error: %s", err)
		return false, err
	}

	if tokenUsername != username {
		log.Printf("jwt local get user error: username does not match token")
		return false, err
	}

	return true, nil
}

func (o *filesJWTChecker) GetSuperuser(username string) (bool, error) {
	return false, nil
}

func (o *filesJWTChecker) CheckAcl(username, topic, clientid string, acc int32) (bool, error) {
	return o.checker.CheckAcl(username, topic, clientid, acc)
}

func (o *filesJWTChecker) Halt() {
	// NO-OP
}
