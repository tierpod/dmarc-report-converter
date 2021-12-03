package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	imap "github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

func fetchIMAPAttachments(cfg *config) error {
	log.Printf("[INFO] imap: connecting to server %v", cfg.Input.IMAP.Server)

	// connect to server
	c, err := client.DialTLS(cfg.Input.IMAP.Server, nil)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] imap: connected")
	defer func() {
		log.Printf("[DEBUG] imap: logout")
		if err := c.Logout(); err != nil {
			log.Printf("[ERROR] imap: logout error %v", err)
		}
	}()

	if cfg.Input.IMAP.Debug {
		log.Printf("[DEBUG] imap: enable debug")
		c.SetDebug(os.Stdout)
	}

	// login
	err = c.Login(cfg.Input.IMAP.Username, cfg.Input.IMAP.Password)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] imap: logged in")

	// select mailbox
	mbox, err := c.Select(cfg.Input.IMAP.Mailbox, false)
	if err != nil {
		return err
	}

	// get all messages
	if mbox.Messages == 0 {
		return fmt.Errorf("no messages found in mailbox")
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

	messages := make(chan *imap.Message, 1)
	done := make(chan error, 1)

	go func() {
		done <- c.Fetch(seqSet, items, messages)
	}()

	downloadCount := 0
	for msg := range messages {
		if msg == nil {
			return fmt.Errorf("server didn't return message")
		}

		br := msg.GetBody(section)
		if br == nil {
			return fmt.Errorf("server didn't return message body")
		}

		// create a new mail reader
		mr, err := mail.CreateReader(br)
		if err != nil {
			return err
		}

		// process each message's part
		isSuccess := false
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Printf("[ERROR] imap: can't read next part: %v, skip", err)
				break
			}

			switch h := p.Header.(type) {
			case mail.AttachmentHeader:
				// this is an attachment
				filename, err := h.Filename()
				if err != nil {
					log.Printf("[ERROR] imap: %v, skip", err)
					continue
				}
				log.Printf("[INFO] imap: found attachment: %v", filename)

				outFile := filepath.Join(cfg.Input.Dir, filename)
				log.Printf("[INFO] imap: save attachment to: %v", outFile)
				f, err := os.Create(outFile)
				if err != nil {
					log.Printf("[ERROR] imap: %v, skip", err)
					continue
				}

				_, err = io.Copy(f, p.Body)
				if err != nil {
					log.Printf("[ERROR] imap: %v, skip", err)
					continue
				}
				err = f.Close()
				if err != nil {
					log.Printf("[ERROR] imap: %v, skip", err)
					continue
				}
				isSuccess = true
			}
		}

		if isSuccess && cfg.Input.IMAP.Delete {
			log.Printf("[DEBUG] imap: add SeqNum %v to delete set", msg.SeqNum)
			deleteSet.AddNum(msg.SeqNum)
		}
		downloadCount++
	}
	log.Printf("[DEBUG] imap: %v attachments downloaded", downloadCount)

	if err := <-done; err != nil {
		return err
	}

	if cfg.Input.IMAP.Delete {
		log.Printf("[DEBUG] imap: delete emails after fetch")

		// first, mark the messages as deleted
		delItems := imap.FormatFlagsOp(imap.AddFlags, false)
		delFlags := []interface{}{imap.DeletedFlag}

		err := c.Store(deleteSet, delItems, delFlags, nil)
		if err != nil {
			return err
		}

		// then delete it
		if err := c.Expunge(nil); err != nil {
			return err
		}
	}

	return nil
}
