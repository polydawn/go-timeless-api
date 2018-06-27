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

// PivotToWareIDs returns a new and reduced WareSourcing where all data is
// indexed ByWareID for each wareID in the argument set.
// All the ByPackType for a type "tar" will be appended to the ByWareID
// index for all wareIDs of type "tar", and so forth.
// ByModule data is ignored (you should flip that to ByWareID-indexed
// immediately when you load it).
func (ws WareSourcing) PivotToWareIDs(wareIDs map[WareID]struct{}) WareSourcing {
	ws2 := WareSourcing{ByWare: make(map[WareID][]WarehouseLocation, len(wareIDs))}
	for wareID := range wareIDs {
		// Copy over anything already explicitly wareID-indexed.
		ws2.ByWare[wareID] = ws.ByWare[wareID]
		// Append packtype-general info.
		ws2.ByWare[wareID] = append(ws2.ByWare[wareID], ws.ByPackType[wareID.Type]...)
	}
	return ws2
}

// PivotToInputs is a shortcut for calling PivotToWareIDs with the set of
// inputs to a bound Op.
func (ws WareSourcing) PivotToInputs(boundOp BoundOperation) WareSourcing {
	wareIDs := make(map[WareID]struct{}, len(boundOp.InputPins))
	for _, wareID := range boundOp.InputPins {
		wareIDs[wareID] = struct{}{}
	}
	return ws.PivotToWareIDs(wareIDs)
}

// PivotToModuleWare returns WareSourcing where all data is indexed ByWareID
// (like PivotToInputs and PivotToWareIDs), also applying any ByModule-index
// info for the named module.  (This is typically used immediately after
// loading the mirrors info in a module's release catalog, in order to avoid
// needed to carry around the module-oriented info any longer.)
func (ws WareSourcing) PivotToModuleWare(wareID WareID, assumingModName ModuleName) WareSourcing {
	ws2 := WareSourcing{ByWare: make(map[WareID][]WarehouseLocation, 1)}
	// Copy over anything already explicitly wareID-indexed.
	ws2.ByWare[wareID] = ws.ByWare[wareID]
	// Append module info.
	forMod := ws.ByModule[assumingModName]
	if forMod == nil {
		return ws2
	}
	ws2.ByWare[wareID] = append(ws2.ByWare[wareID], forMod[wareID.Type]...)
	return ws2
}
