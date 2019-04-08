package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"

	imap "github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

type bufferMessage struct {
	buf      []byte
	filename string
}

func (bm *bufferMessage) reader() io.Reader {
	return bytes.NewReader(bm.buf)
}

func processIMAP(cfg *config) {
	log.Printf("[INFO] imap: connecting to server %v", cfg.Input.IMAP.Server)

	// connect to server
	c, err := client.DialTLS(cfg.Input.IMAP.Server, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[DEBUG] imap: connected")
	defer c.Logout()

	if cfg.ImapDebug {
		log.Printf("[DEBUG] imap: enable debug")
		c.SetDebug(os.Stdout)
	}

	// login
	err = c.Login(cfg.Input.IMAP.Username, cfg.Input.IMAP.Password)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[DEBUG] imap: logged in")

	// select mailbox
	mbox, err := c.Select(cfg.Input.IMAP.Mailbox, false)
	if err != nil {
		log.Fatal(err)
	}

	// get all messages
	if mbox.Messages == 0 {
		log.Printf("[WARN] imap: no message in mailbox")
		return
	}

	from := uint32(1)
	to := mbox.Messages

	log.Printf("[INFO] imap: found %v messages, %v unseen", mbox.Messages, mbox.Unseen)

	// set for all messages
	seqSet := new(imap.SeqSet)
	seqSet.AddRange(from, to)

	// set for delete messages
	deleteSet := new(imap.SeqSet)

	// get the whole message body
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	go func() {
		done <- c.Fetch(seqSet, items, messages)
	}()

	var bufferMessages []bufferMessage
	downloadCount := 0
	for msg := range messages {
		if msg == nil {
			log.Printf("[ERROR] imap: server didn't returned message")
			return
		}

		br := msg.GetBody(section)
		if br == nil {
			log.Printf("[ERROR] imap: server didn't returned message body")
			return
		}

		// create a new mail reader
		mr, err := mail.CreateReader(br)
		if err != nil {
			log.Printf("[ERROR] imap: %v", err)
			return
		}

		// process each message's part
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			switch h := p.Header.(type) {
			case mail.AttachmentHeader:
				// this is an attachment
				filename, _ := h.Filename()
				log.Printf("[INFO] imap: found attachment: %v", filename)
				// download to buffer to prevent long imap connection
				buf, err := ioutil.ReadAll(p.Body)
				if err != nil {
					log.Printf("[ERROR] imap: %v, skip", err)
					continue
				}
				bufferMessages = append(bufferMessages, bufferMessage{buf, filename})
			}
		}

		log.Printf("[DEBUG] imap: add SeqNum %v to delete set", msg.SeqNum)
		deleteSet.AddNum(msg.SeqNum)
		downloadCount++
	}
	log.Printf("[DEBUG] imap: %v attachments downloaded", downloadCount)

	if err := <-done; err != nil {
		log.Printf("[ERROR] imap: %v", err)
		return
	}

	if cfg.Input.Delete {
		log.Printf("[DEBUG] imap: delete emails after converting")

		// first, mark the messages as deleted
		delItems := imap.FormatFlagsOp(imap.AddFlags, false)
		delFlags := []interface{}{imap.DeletedFlag}

		err := c.Store(deleteSet, delItems, delFlags, nil)
		if err != nil {
			log.Printf("[ERROR] imap: %v", err)
			return
		}

		// then delete it
		if err := c.Expunge(nil); err != nil {
			log.Printf("[ERROR] imap: %v", err)
			return
		}
	}

	log.Printf("[DEBUG] imap: logout")
	if err := c.Logout(); err != nil {
		log.Printf("[ERROR] imap: logout error %v", err)
	}

	doneCount := 0
	for _, bm := range bufferMessages {
		err = readConvert(bm.reader(), bm.filename, cfg)
		if err != nil {
			log.Printf("[ERROR] imap: %v, skip", err)
			continue
		}
		doneCount++
	}
	log.Printf("[INFO] imap: %v files converted", doneCount)
}
