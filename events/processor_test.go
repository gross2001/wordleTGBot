package events

import (
	"testing"
	"wordle/clients/telegram"
	"wordle/vocab/sqlite"
)

func TestProcessor_checkWordIsOk(t *testing.T) {
	type fields struct {
		vocab      sqlite.TatSQLVocab
		client     telegram.Client
		currentDay DayInfo
		users      map[int64]UsersAnswer
		dialogs    Dialogs
	}
	type args struct {
		word string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Processor{
				vocab:      tt.fields.vocab,
				client:     tt.fields.client,
				currentDay: tt.fields.currentDay,
				users:      tt.fields.users,
				dialogs:    tt.fields.dialogs,
			}
			got, got1 := p.checkWordIsOk(tt.args.word)
			if got != tt.want {
				t.Errorf("Processor.checkWordIsOk() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Processor.checkWordIsOk() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
