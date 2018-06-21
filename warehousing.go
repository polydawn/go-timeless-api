package api

type WarehouseLocation string

type WareSourcing struct {
	ByKind   map[PackType][]WarehouseLocation
	ByModule map[ModuleName]map[PackType][]WarehouseLocation
	ByWare   map[WareID][]WarehouseLocation
}
