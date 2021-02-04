package channeldata

import "html/template"

// MailData holds info for sending an email
type MailData struct {
	ToName       string
	ToAddress    string
	FromName     string
	FromAddress  string
	AdditionalTo []string
	Subject      string
	Content      template.HTML
	Template     string
	CC           []string
	UseHermes    bool
	Attachments  []string
	StringMap    map[string]string
	IntMap       map[string]int
	FloatMap     map[string]float32
	RowSets      map[string]interface{}
}

// MailJob is the unit of work to be performed when sending an email to chan
type MailJob struct {
	MailMessage MailData
}
