package gormplugin

import "gorm.io/gorm"

const (
	_createOp = "gorm:create"
	_updateOp = "gorm:update"
	_queryOp  = "gorm:query"
	_deleteOp = "gorm:delete"
	_rowOp    = "gorm:row"
	_rawOp    = "gorm:raw"
)

func (op OpentracingPlugin) beforeCreate(db *gorm.DB) {
	op.injectBefore(db, _createOp)
}

func (op OpentracingPlugin) beforeUpdate(db *gorm.DB) {
	op.injectBefore(db, _updateOp)
}

func (op OpentracingPlugin) beforeQuery(db *gorm.DB) {
	op.injectBefore(db, _queryOp)
}

func (op OpentracingPlugin) beforeDelete(db *gorm.DB) {
	op.injectBefore(db, _deleteOp)
}

func (op OpentracingPlugin) beforeRow(db *gorm.DB) {
	op.injectBefore(db, _rowOp)
}

func (op OpentracingPlugin) beforeRaw(db *gorm.DB) {
	op.injectBefore(db, _rawOp)
}

func (op OpentracingPlugin) after(db *gorm.DB) {
	op.extractAfter(db)
}
