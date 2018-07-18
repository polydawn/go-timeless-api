package api

type WarehouseLocation string

type WareSourcing struct {
	ByPackType map[PackType][]WarehouseLocation                `refmt:",omitempty"`
	ByModule   map[ModuleName]map[PackType][]WarehouseLocation `refmt:",omitempty"`
	ByWare     map[WareID][]WarehouseLocation                  `refmt:",omitempty"`
}
