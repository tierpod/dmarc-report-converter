package main

import (
	"io"
	"io/ioutil"
	"log"
	"path/filepath"

	imap "github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message"
)

func fetchIMAPAttachments(cfg *config) error {

	log.Println("[INFO] Connecting to server...")

	c, err := client.DialTLS(cfg.Input.IMAP.Server, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[INFO] Connected to: %s ", cfg.Input.IMAP.Server)

	defer c.Logout()

	if err := c.Login(cfg.Input.IMAP.Username, cfg.Input.IMAP.Password); err != nil {
		log.Fatal(err)
	}
	log.Printf("[INFO] Logged in as: %s", cfg.Input.IMAP.Username)

	mbox, err := c.Select(cfg.Input.IMAP.Mailbox, false)
	if err != nil {
		log.Fatal(err)
	}

	seqset := new(imap.SeqSet)
	seqset.AddRange(1, mbox.Messages)

	// set for delete messages
	deleteSet := new(imap.SeqSet)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchRFC822}, messages)
	}()

	countMessages := 0
	countProcessedMessages := 0

	for msg := range messages {
		downloadSuccess := false
		countMessages += 1
		for _, r := range msg.Body {

			entity, err := message.Read(r)
			if err != nil {

				log.Fatal(err)
			}

			multiPartReader := entity.MultipartReader()

			for e, err := multiPartReader.NextPart(); err != io.EOF; e, err = multiPartReader.NextPart() {

				kind, _, cErr := e.Header.ContentType()
				if cErr != nil {
					log.Fatal(cErr)
				}

				if kind != "application/gzip" && kind != "application/zip" && kind != "application/octet-stream" {
					continue
				}

				_, v, _ := e.Header.ContentDisposition()

				c, rErr := ioutil.ReadAll(e.Body)
				if rErr != nil {
					log.Fatal(rErr)
				}

				outFile := filepath.Join(cfg.Input.Dir, v["filename"])

				log.Printf("[INFO] * Extracting file %s", outFile)

				if fErr := ioutil.WriteFile(outFile, c, 0777); err != nil {
					log.Fatal(fErr)
				}
				downloadSuccess = true
				countProcessedMessages += 1
			}
		}
		if downloadSuccess && cfg.Input.IMAP.Delete {
			log.Printf("[DEBUG] imap: add SeqNum %v to delete set", msg.SeqNum)
			deleteSet.AddNum(msg.SeqNum)
		}
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	if countProcessedMessages > 0 && cfg.Input.IMAP.Delete {
		log.Printf("[INFO] imap: delete emails after fetch")
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

	log.Printf("[INFO] Total messages: %d, Processed messages: %d", countMessages, countProcessedMessages)
	return nil
}
