package service

import (
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPubSub_SubscribeAndPublish(t *testing.T) {
	ps := NewPubSub()
	postID := "post123"

	ch, unsubscribe := ps.Subscribe(postID)
	defer unsubscribe()

	comment := &models.Comment{ID: "comment1", PostID: postID, Text: "Test comment"}

	ps.Publish(postID, comment)

	select {
	case received := <-ch:
		require.Equal(t, comment, received, "Полученный комментарий не совпадает с отправленным")
	case <-time.After(1 * time.Second):
		require.Fail(t, "Сообщение не пришло подписчику")
	}
}

func TestPubSub_MultipleSubscribers(t *testing.T) {
	ps := NewPubSub()
	postID := "post456"

	ch1, unsubscribe1 := ps.Subscribe(postID)
	defer unsubscribe1()
	ch2, unsubscribe2 := ps.Subscribe(postID)
	defer unsubscribe2()

	comment := &models.Comment{ID: "comment2", PostID: postID, Text: "Another comment"}

	ps.Publish(postID, comment)

	for _, ch := range []<-chan *models.Comment{ch1, ch2} {
		select {
		case received := <-ch:
			require.Equal(t, comment, received, "Один из подписчиков получил неверное сообщение")
		case <-time.After(1 * time.Second):
			require.Fail(t, "Один из подписчиков не получил сообщение")
		}
	}
}

func TestPubSub_Unsubscribe(t *testing.T) {
	ps := NewPubSub()
	postID := "post789"

	ch, unsubscribe := ps.Subscribe(postID)
	unsubscribe()

	comment := &models.Comment{ID: "comment3", PostID: postID, Text: "No one should get this"}

	ps.Publish(postID, comment)

	select {
	case msg, ok := <-ch:
		if ok {
			require.Fail(t, "Отписанный подписчик все равно получил сообщение: %v", msg)
		}
	case <-time.After(500 * time.Millisecond):
	}
}
