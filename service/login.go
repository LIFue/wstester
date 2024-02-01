package service

type Login struct {
	Uid      string
	Password string
	LoginUri string
	Version  string
	Platform string
}

//func NewLogin(uid, password, loginUri string) Login {
//	return Login{
//		Uid:      uid,
//		Password: password,
//		LoginUri: loginUri,
//		Version:  model.VERSION,
//		Platform: model.PLATFORM_WEB,
//	}
//}
//
//func (l Login) Login(server model.Server) (string, error) {
//	var urlParam = model.NewLoginUrlParam(l.Version, l.Platform)
//	var payload = model.NewReqAuthMethodPayload(l.Uid)
//
//	var authMethod = new(model.LoginPrepareResp)
//	loginUrl := l.genLoginUrl(server)
//	log.Infof("loginUrl: %s", loginUrl)
//	if err := new(httputil.HttpUtil).SendRequest(
//		httputil.NewRequest(http.MethodPost, loginUrl, payload), authMethod); err != nil {
//		log.Errorf("get auth method err: %s", err.Error())
//		return "", errors.Wrap(err, "request auth method error")
//	}
//
//	encodePassword := utils.EncodeMD5(l.Password)
//	loginPassword := utils.EncodeMD5(fmt.Sprintf("%s:%s:%s:%s", l.Uid, authMethod.Realm, encodePassword, authMethod.Nonce))
//
//	loginParam := model.NewLoginParam(l.Uid, loginPassword)
//	urlParam = model.NewLoginUrlParam(l.Version, l.Platform)
//	var loginResp model.RespLogin
//	if err := new(httputil.HttpUtil).SendRequest(
//		httputil.NewRequest(http.MethodPost, loginUrl, loginParam), &loginResp); err != nil {
//		log.Errorf("login err: %s", err.Error())
//		return "", errors.Wrap(err, "request auth method error")
//	}
//
//	if len(loginResp.Ws) == 0 {
//		return "", errors.New("login resp error, ws is nil")
//	}
//
//	return loginResp.Ws[0], nil
//}
//
//func (l Login) genLoginUrl(server model.Server) string {
//	return fmt.Sprintf("%s://%s:%s%s", server.Protocol, server.Ip, server.Port, l.LoginUri)
//}
