package manager

import (
	"chat/adapter"
	adaptercommon "chat/adapter/common"
	"chat/addition/web"
	"chat/admin"
	"chat/auth"
	"chat/channel"
	"chat/globals"
	"chat/manager/conversation"
	"chat/utils"
	"time"

	"database/sql"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

const defaultMessage = "empty response"
const interruptMessage = "interrupted"

func CollectQuota(c *gin.Context, user *auth.User, buffer *utils.Buffer, uncountable bool, err error) {
	db := utils.GetDBFromContext(c)
	quota := buffer.GetQuota()

	if user == nil || quota <= 0 {
		return
	}

	if buffer.IsEmpty() || err != nil {
		return
	}

	if !uncountable {
		user.UseQuota(db, quota)
	}
}

type partialChunk struct {
	Chunk *globals.Chunk
	End   bool
	Hit   bool
	Error error
}

func createStopSignal(conn *Connection) chan bool {
	stopSignal := make(chan bool, 1)
	go func(conn *Connection, stopSignal chan bool) {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer func() {
			ticker.Stop()
			if r := recover(); r != nil && !strings.Contains(fmt.Sprintf("%s", r), "closed channel") {
				stack := debug.Stack()
				globals.Warn(fmt.Sprintf("caught panic from stop signal: %s\n%s", r, stack))
			}
		}()

		for {
			select {
			case <-ticker.C:
				state := conn.PeekStop() != nil // check the stop state
				stopSignal <- state

				if state {
					break
				}
			}
		}
	}(conn, stopSignal)

	return stopSignal
}

func createChatTask(
	conn *Connection, user *auth.User, buffer *utils.Buffer, db *sql.DB, cache *redis.Client,
	model string, instance *conversation.Conversation, segment []globals.Message, plan bool,
) (hit bool, err error) {
	chunkChan := make(chan partialChunk, 24) // the channel to send the chunk data
	interruptSignal := make(chan error, 1)   // the signal to interrupt the chat task routine
	stopSignal := createStopSignal(conn)     // the signal to stop from the client

	defer func() {
		// close all channels
		close(interruptSignal)
		close(stopSignal)
		close(chunkChan)
	}()

	// create a new chat request routine
	go func() {
		defer func() {
			if r := recover(); r != nil && !strings.Contains(fmt.Sprintf("%s", r), "closed channel") {
				stack := debug.Stack()
				globals.Warn(fmt.Sprintf("caught panic from chat request: %s\n%s", r, stack))
			}
		}()

		if globals.IsVideoModel(model) {
			props := adaptercommon.CreateVideoProps(&adaptercommon.VideoProps{
				Model:  model,
				Prompt: segment[len(segment)-1].Content,
			})
			props.User = auth.GetUsernameString(db, user)

			var finalJobJson string
			hit, err := channel.NewVideoRequestWithCache(
				cache, buffer,
				auth.GetGroup(db, user),
				props,
				func(data *globals.Chunk) error {
					if data != nil && data.Content != "" {
						if strings.HasPrefix(data.Content, "{") && strings.Contains(data.Content, "\"id\"") && strings.Contains(data.Content, "\"status\"") {
							finalJobJson = data.Content

							job, err := utils.UnmarshalString[RelayVideoJob](data.Content)
							if err == nil && job.Id != "" && job.Status == "completed" {
								backendUrl := channel.SystemInstance.GetBackend()
								videoUrl := fmt.Sprintf("%s/v1/videos/%s/content", backendUrl, job.Id)
								videoMarkdown := utils.GetVideoMarkdown(videoUrl, "video")

								chunkChan <- partialChunk{Chunk: &globals.Chunk{Content: videoMarkdown}, End: false, Hit: false, Error: nil}
								return nil
							}
						}
					}
					// Send original content for progress updates and other messages
					chunkChan <- partialChunk{Chunk: data, End: false, Hit: false, Error: nil}
					return nil
				},
			)

			if instance != nil && finalJobJson != "" {
				job, err := utils.UnmarshalString[RelayVideoJob](finalJobJson)
				if err != nil {
					globals.Warn(fmt.Sprintf("[video] failed to parse job JSON: %s, finalJobJson: %s", err.Error(), finalJobJson))
				} else if job.Id == "" {
					globals.Warn(fmt.Sprintf("[video] job.Id is empty after parsing, finalJobJson: %s", finalJobJson))
				} else {
					globals.Debug(fmt.Sprintf("[video] saving task_id %s to conversation %d", job.Id, instance.GetId()))
					instance.SetTaskID(job.Id)
					if !instance.SaveConversation(db) {
						globals.Warn(fmt.Sprintf("[video] failed to save conversation with task_id %s", job.Id))
					} else {
						globals.Debug(fmt.Sprintf("[video] successfully saved task_id %s to conversation %d", job.Id, instance.GetId()))
					}
				}
			} else {
				if instance == nil {
					globals.Warn("[video] instance is nil, cannot save task_id")
				} else if finalJobJson == "" {
					globals.Warn("[video] finalJobJson is empty, cannot save task_id")
				}
			}

			chunkChan <- partialChunk{Chunk: nil, End: true, Hit: hit, Error: err}
			return
		}

		hit, err := channel.NewChatRequestWithCache(
			cache, buffer,
			auth.GetGroup(db, user),
			adaptercommon.CreateChatProps(&adaptercommon.ChatProps{
				Model:             model,
				Message:           segment,
				MaxTokens:         instance.GetMaxTokens(),
				Temperature:       instance.GetTemperature(),
				TopP:              instance.GetTopP(),
				TopK:              instance.GetTopK(),
				PresencePenalty:   instance.GetPresencePenalty(),
				FrequencyPenalty:  instance.GetFrequencyPenalty(),
				RepetitionPenalty: instance.GetRepetitionPenalty(),
			}, buffer),

			// the function to handle the chunk data
			func(data *globals.Chunk) error {
				// if interrupt signal is received
				if len(interruptSignal) > 0 {
					return errors.New(interruptMessage)
				}

				// send the chunk data to the channel
				chunkChan <- partialChunk{
					Chunk: data,
					End:   false,
					Hit:   false,
					Error: nil,
				}
				return nil
			},
		)

		// chat request routine is done
		chunkChan <- partialChunk{
			Chunk: nil,
			End:   true,
			Hit:   hit,
			Error: err,
		}
	}()

	for {
		select {
		case data := <-chunkChan:
			if data.Error != nil && data.Error.Error() == interruptMessage {
				// skip the interrupt message
				continue
			}

			hit = data.Hit
			err = data.Error

			if data.End {
				return
			}

			if err := conn.SendClient(globals.ChatSegmentResponse{
				Message: buffer.WriteChunk(data.Chunk),
				Quota:   buffer.GetQuota(),
				End:     false,
				Plan:    plan,
			}); err != nil {
				globals.Warn(fmt.Sprintf("failed to send message to client: %s", err.Error()))
				interruptSignal <- err
				return hit, nil
			}

		case signal := <-stopSignal:
			// if stop signal is received
			if signal {
				globals.Info(fmt.Sprintf("client stopped the chat request (model: %s, client: %s)", model, conn.GetCtx().ClientIP()))
				_ = conn.SendClient(globals.ChatSegmentResponse{
					Quota: buffer.GetQuota(),
					End:   true,
					Plan:  plan,
				})
				interruptSignal <- errors.New("signal")

				return
			}
		}
	}
}

func ChatHandler(conn *Connection, user *auth.User, instance *conversation.Conversation, restart bool) string {
	defer func() {
		if err := recover(); err != nil {
			stack := debug.Stack()
			globals.Warn(fmt.Sprintf("caught panic from chat handler: %s (instance: %s, client: %s)\n%s",
				err, instance.GetModel(), conn.GetCtx().ClientIP(), stack,
			))
		}
	}()

	db := conn.GetDB()
	cache := conn.GetCache()

	model := instance.GetModel()
	segment := adapter.ClearMessages(model, web.ToChatSearchedContext(conn.GetCtx().Request.Context(), instance, restart))

	check, plan := auth.CanEnableModelWithSubscription(db, cache, user, model, segment)
	conn.Send(globals.ChatSegmentResponse{
		Conversation: instance.GetId(),
	})

	if check != nil {
		message := check.Error()
		conn.Send(globals.ChatSegmentResponse{
			Message: message,
			Quota:   0,
			End:     true,
		})
		return message
	}

	buffer := utils.NewBuffer(model, segment, channel.ChargeInstance.GetCharge(model))
	hit, err := createChatTask(conn, user, buffer, db, cache, model, instance, segment, plan)

	admin.AnalyseRequest(model, buffer, err)
	if adapter.IsAvailableError(err) {
		globals.Warn(fmt.Sprintf("%s (model: %s, client: %s)", err, model, conn.GetCtx().ClientIP()))

		auth.RevertSubscriptionUsage(db, cache, user, model)
		conn.Send(globals.ChatSegmentResponse{
			Message: err.Error(),
			End:     true,
		})
		return err.Error()
	}

	if !hit {
		CollectQuota(conn.GetCtx(), user, buffer, plan, err)
	}

	if buffer.IsEmpty() {
		conn.Send(globals.ChatSegmentResponse{
			Message: defaultMessage,
			End:     true,
		})
		return defaultMessage
	}

	conn.Send(globals.ChatSegmentResponse{
		End:   true,
		Quota: buffer.GetQuota(),
		Plan:  plan,
	})

	return buffer.ReadWithDefault(defaultMessage)
}
