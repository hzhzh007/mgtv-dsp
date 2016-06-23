package logic

const (
	UidImeiType     = "imei"
	UidImeiSha1Type = "imei_sha1"
	UidImeiMd5Type  = "imei_md5"
	UidAnidType     = "android_id"
	UidMacSha1Type  = "mac_sha1"
	UidMacMd5Type   = "mac_md5"
	UidOpenudidType = "open_udid"
	UidIdfaType     = "idfa"
	UidCookieType   = "cookie_id"
)

type Uid struct {
	Type  string
	Value string
}
type Uids []Uid

func (uids *Uids) AddId(idType, value string) bool {
	if value != "" {
		*uids = append(*uids, Uid{Type: idType, Value: value})
		return true
	}
	return false
}
