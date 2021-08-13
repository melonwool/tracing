package gormplugin

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tracerLog "github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

const (
	callBackBeforeName = "opentracing:before"
	callBackAfterName  = "opentracing:after"
	opentracingSpanKey = "opentracing:span"
)

type OpentracingPlugin struct{}

func (op *OpentracingPlugin) Name() string {
	return "opentracingPlugin"
}

func (op *OpentracingPlugin) Initialize(db *gorm.DB) (err error) {
	// 开始前
	db.Callback().Create().Before(_gormCreate.Before()).Register(callBackBeforeName, op.beforeCreate)
	db.Callback().Query().Before(_gormQuery.Before()).Register(callBackBeforeName, op.beforeQuery)
	db.Callback().Delete().Before(_gormDelete.Before()).Register(callBackBeforeName, op.beforeDelete)
	db.Callback().Update().Before(_gormUpdate.Before()).Register(callBackBeforeName, op.beforeUpdate)
	db.Callback().Row().Before(_gormRow.Before()).Register(callBackBeforeName, op.beforeRow)
	db.Callback().Raw().Before(_gormRaw.Before()).Register(callBackBeforeName, op.beforeRaw)

	// 结束后
	db.Callback().Create().After(_gormCreate.After()).Register(callBackAfterName, op.after)
	db.Callback().Query().After(_gormQuery.After()).Register(callBackAfterName, op.after)
	db.Callback().Delete().After(_gormDelete.After()).Register(callBackAfterName, op.after)
	db.Callback().Update().After(_gormUpdate.After()).Register(callBackAfterName, op.after)
	db.Callback().Row().After(_gormRow.After()).Register(callBackAfterName, op.after)
	db.Callback().Raw().After(_gormRaw.After()).Register(callBackAfterName, op.after)
	return
}

func (op OpentracingPlugin) injectBefore(db *gorm.DB, opName string) {
	span, _ := opentracing.StartSpanFromContext(db.Statement.Context, opName)
	ext.DBType.Set(span, "mysql")
	db.InstanceSet(opentracingSpanKey, span)
}

func (op OpentracingPlugin) extractAfter(db *gorm.DB) {
	v, ok := db.InstanceGet(opentracingSpanKey)
	if !ok || v == nil {
		return
	}
	span, ok := v.(opentracing.Span)
	if !ok || span == nil {
		return
	}
	span.LogFields(tracerLog.String("sql", db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)))
	defer span.Finish()
}
