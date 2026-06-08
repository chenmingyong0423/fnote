package web

import (
	"time"

	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"
	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/domain"
	"github.com/chenmingyong0423/gkit"
	"github.com/gin-gonic/gin"
)

func (h *WebsiteConfigHandler) AdminGetConfigCompletion(ctx *gin.Context) (*apiwrap.ResponseBody[ConfigCompletionVO], error) {
	health, err := h.serv.GetConfigHealth(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]ConfigCompletionItemVO, 0, len(health.Items))
	done := 0
	for _, item := range health.Items {
		if item.Configured {
			done++
		}
		items = append(items, ConfigCompletionItemVO{
			Key:           item.Key,
			Label:         item.Label,
			Configured:    item.Configured,
			MissingFields: item.MissingFields,
			Href:          item.Href,
			Description:   item.Description,
		})
	}
	return apiwrap.SuccessResponseWithData(ConfigCompletionVO{
		Completed: done == len(items),
		Total:     len(items),
		Done:      done,
		Items:     items,
	}), nil
}

func (h *WebsiteConfigHandler) AdminGetConfigHealth(ctx *gin.Context) (*apiwrap.ResponseBody[ConfigHealthVO], error) {
	health, err := h.serv.GetConfigHealth(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(toConfigHealthVO(health)), nil
}

func (h *WebsiteConfigHandler) AdminUpdateConfigCheckState(ctx *gin.Context, req UpdateConfigCheckStateReq) (*apiwrap.ResponseBody[any], error) {
	var snoozedUntil *time.Time
	if req.SnoozedUntil != nil {
		snoozedTime := time.Unix(*req.SnoozedUntil, 0).Local()
		snoozedUntil = &snoozedTime
	} else if req.SnoozeDays > 0 {
		snoozedTime := time.Now().Local().AddDate(0, 0, req.SnoozeDays)
		snoozedUntil = &snoozedTime
	}
	err := h.serv.UpdateConfigCheckState(ctx, domain.ConfigCheckState{
		CheckKey:      ctx.Param("key"),
		Status:        req.Status,
		IgnoredReason: req.IgnoredReason,
		SnoozedUntil:  snoozedUntil,
	})
	return apiwrap.SuccessResponse(), err
}

func toConfigHealthVO(health *domain.ConfigHealth) ConfigHealthVO {
	items := make([]ConfigHealthItemVO, 0, len(health.Items))
	for _, item := range health.Items {
		items = append(items, toConfigHealthItemVO(item))
	}
	return ConfigHealthVO{
		Score:       health.Score,
		Required:    toConfigHealthGroupVO(health.Required),
		Recommended: toConfigHealthGroupVO(health.Recommended),
		Optional:    toConfigHealthGroupVO(health.Optional),
		Items:       items,
	}
}

func toConfigHealthGroupVO(group domain.ConfigHealthGroup) ConfigHealthGroupVO {
	return ConfigHealthGroupVO{
		Done:      group.Done,
		Total:     group.Total,
		Completed: group.Completed,
	}
}

func toConfigHealthItemVO(item domain.ConfigHealthItem) ConfigHealthItemVO {
	var snoozedUntil *int64
	if item.SnoozedUntil != nil {
		snoozedUntil = gkit.ToPtr(item.SnoozedUntil.Unix())
	}
	return ConfigHealthItemVO{
		Key:           item.Key,
		Label:         item.Label,
		Level:         item.Level,
		Status:        item.Status,
		Configured:    item.Configured,
		MissingFields: item.MissingFields,
		Description:   item.Description,
		Href:          item.Href,
		SnoozedUntil:  snoozedUntil,
		IgnoredReason: item.IgnoredReason,
	}
}
