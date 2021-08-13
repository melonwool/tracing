package gormplugin

type gormOption string

func (g gormOption) String() string {
	return string(g)
}
func (g gormOption) Before() string {
	return string(g) + "_before"
}
func (g gormOption) After() string {
	return string(g) + "_after"
}

const (
	_gormCreate gormOption = "gorm:create"
	_gormUpdate gormOption = "gorm:update"
	_gormQuery  gormOption = "gorm:query"
	_gormDelete gormOption = "gorm:delete"
	_gormRow    gormOption = "gorm:row"
	_gormRaw    gormOption = "gorm:raw"
)

type operationStage string

func (op operationStage) Name() string {
	return string(op)
}

func (op operationStage) Before() string {
	return string(op) + "_before"
}
func (op operationStage) After() string {
	return string(op) + "_after"
}

const (
	_stageBeforeCreate operationStage = "opentracing:before_create"
	_stageAfterCreate  operationStage = "opentracing:after_create"
	_stageBeforeUpdate operationStage = "opentracing:before_update"
	_stageAfterUpdate  operationStage = "opentracing:after_update"
	_stageBeforeQuery  operationStage = "opentracing:before_query"
	_stageAfterQuery   operationStage = "opentracing:after_query"
	_stageBeforeDelete operationStage = "opentracing:before_delete"
	_stageAfterDelete  operationStage = "opentracing:after_delete"
	_stageBeforeRow    operationStage = "opentracing:before_row"
	_stageAfterRow     operationStage = "opentracing:after_row"
	_stageBeforeRaw    operationStage = "opentracing:before_raw"
	_stageAfterRaw     operationStage = "opentracing:after_raw"
)
