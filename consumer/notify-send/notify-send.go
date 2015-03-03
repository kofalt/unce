package notifysend


import (
	"os/exec"
	. "kofalt.com/unce/def"
)


type NotifySendConsumer struct {

}

func New() *NotifySendConsumer {
	return &NotifySendConsumer{}
}

func (ns *NotifySendConsumer) Consume(e *Event) {
	exec.Command("notify-send", "--urgency", "normal", "--expire-time", "10", e.Summary, e.Message).Start()
}
