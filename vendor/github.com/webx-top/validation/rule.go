package validation

type Rule struct {
	ZipCode   string
	Telephone string
	Mobile    string
	Base64    string
	IPv4      string
	Email     string
	AlphaDash string
}

func (r *Rule) GetPhone() string {
	mobile := "(" + r.Mobile[1:len(r.Mobile)-1] + ")"
	tel := "(" + r.Telephone[1:len(r.Telephone)-1] + ")"
	return "^" + mobile + "|" + tel + "$"
}

var DefaultRule = &Rule{
	ZipCode:   "^[1-9]\\d{5}$",
	Telephone: "^(0\\d{2,3}(\\-)?)?\\d{7,8}$",
	Mobile:    "^((\\+86)|(86))?(1(([35][0-9])|[8][0-9]|[7][01356789]|[4][579]|[6][2567]))\\d{8}$",
	Base64:    "^(?:[A-Za-z0-99+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$",
	IPv4:      "^((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)$",
	Email:     "[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?",
	AlphaDash: "[^\\w-]",
}
