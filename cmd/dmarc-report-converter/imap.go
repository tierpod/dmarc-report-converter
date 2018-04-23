package main

import (
	"io"
	"log"

	imap "github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

func processIMAP(cfg *config) {
	log.Printf("[INFO] imap: connecting to server %v", cfg.Input.IMAP.Server)

	// connect to server
	c, err := client.DialTLS(cfg.Input.IMAP.Server, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[DEBUG] imap: connected")
	defer c.Logout()

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

	log.Printf("[INFO] imap: found messages %v, unseen %v", mbox.Messages, mbox.Unseen)

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

	doneCount := 0
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

				err = readConvert(p.Body, filename, cfg)
				if err != nil {
					log.Printf("[ERROR] imap: %v, skip", err)
					continue
				}
			}
		}

		log.Printf("[DEBUG] imap: add SeqNum %v to delete set", msg.SeqNum)
		deleteSet.AddNum(msg.SeqNum)
		doneCount++
	}

	if err := <-done; err != nil {
		log.Printf("[ERROR] imap: %v", err)
		return
	}

	if cfg.Input.Delete {
		log.Printf("[DEBUG] imap: delete emails after converting")

		delItems := imap.FormatFlagsOp(imap.AddFlags, false)
		delFlags := []interface{}{imap.SeenFlag, imap.DeletedFlag}

		err := c.Store(deleteSet, delItems, delFlags, nil)
		if err != nil {
			log.Printf("[ERROR] imap: %v", err)
			return
		}
	}

	log.Printf("[INFO] imap: done %v items", doneCount)
}
