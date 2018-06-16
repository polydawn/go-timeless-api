package api

type SubmoduleSlotRefList []SubmoduleSlotRef

func (s SubmoduleSlotRefList) Len() int      { return len(s) }
func (s SubmoduleSlotRefList) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SubmoduleSlotRefList) Less(i, j int) bool {
	return s[i].SubmoduleRef < s[j].SubmoduleRef || s[i].SlotName < s[j].SlotName
}
