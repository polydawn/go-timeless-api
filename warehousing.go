package api

type WarehouseLocation string

// WareSourcing contains suggestions on WarehouseLocations which may be
// able to provide Wares.
//
// This information may be indexed in several different ways: most specifically
// (and inflexibly, and verbosely) by specific WareID; or by module name;
// or by pack type in general.  (Non-content-addressible WarehouseLocations
// only semantically make sense when indexed by specific WareID; since the
// other forms of indexing will recommend the WarehouseLocation for more than
// one specific WareID, it stands to reason that they the WarehouseLocation
// ought to specify a system which can store more than one Ware!)
//
// WareSourcing is meant to be reasonable to provide to *more than one*
// Operation (each of which may also have more than one input, of course) --
// the various mechanisms of indexing allow such generalized suggestions.
type WareSourcing struct {
	ByPackType map[PackType][]WarehouseLocation                `refmt:",omitempty"`
	ByModule   map[ModuleName]map[PackType][]WarehouseLocation `refmt:",omitempty"`
	ByWare     map[WareID][]WarehouseLocation                  `refmt:",omitempty"`
}

// WareStaging contains instructions on where to store wares that are
// output from an Operation.
//
// WareStaging only takes a single warehouse location per packtype.
// It is intended that if you want to replicate the ware storage to
// multiple locations, you should do this later, *not* while saving
// the output from the Operation.  An Operation may fail if the
// WarehouseLocation provided by the WareStaging info is not writable.
//
// It is semantically unreasonable to provide a non-content-addressable
// WarehouseLocation in WareStaging info: WareStaging info is meant to
// be reasonable to provide to *more than one* Operation (each of which
// may also have more than one output, of course) -- therefore it is only
// sensible to provide a WarehouseLocation which is capable of storing
// more than one Ware!  (You may still run Repeatr with non-CA
// WarehouseLocation configurations for specific outputs; it's only the
// higher level pipelining tools which become opinionated about this.)
type WareStaging struct {
	ByPackType map[PackType]WarehouseLocation
}
