package api

func (ws *WareSourcing) Append(ws2 WareSourcing) {
	for packtype, locations := range ws2.ByPackType {
		ws.AppendByPackType(packtype, locations...)
	}
	for modName, v := range ws2.ByModule {
		for packtype, locations := range v {
			ws.AppendByModule(modName, packtype, locations...)
		}
	}
	for wareID, locations := range ws2.ByWare {
		ws.AppendByWare(wareID, locations...)
	}
}

func (ws *WareSourcing) AppendByPackType(packtype PackType, locations ...WarehouseLocation) {
	if ws.ByPackType == nil {
		ws.ByPackType = make(map[PackType][]WarehouseLocation)
	}
	if ws.ByPackType[packtype] == nil {
		ws.ByPackType[packtype] = locations
		return
	}
	ws.ByPackType[packtype] = append(ws.ByPackType[packtype], locations...)
}

func (ws *WareSourcing) AppendByModule(modName ModuleName, packtype PackType, locations ...WarehouseLocation) {
	if ws.ByModule == nil {
		ws.ByModule = make(map[ModuleName]map[PackType][]WarehouseLocation)
	}
	if ws.ByModule[modName] == nil {
		ws.ByModule[modName] = make(map[PackType][]WarehouseLocation)
	}
	if ws.ByModule[modName][packtype] == nil {
		ws.ByModule[modName][packtype] = locations
		return
	}
	ws.ByModule[modName][packtype] = append(ws.ByModule[modName][packtype], locations...)
}

func (ws *WareSourcing) AppendByWare(wareID WareID, locations ...WarehouseLocation) {
	if ws.ByWare == nil {
		ws.ByWare = make(map[WareID][]WarehouseLocation)
	}
	if ws.ByWare[wareID] == nil {
		ws.ByWare[wareID] = locations
		return
	}
	ws.ByWare[wareID] = append(ws.ByWare[wareID], locations...)
}
