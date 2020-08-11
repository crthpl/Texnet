package items

type ItemStack struct {
	amnt int8		//can be at most 85 (legally) and at the least 1
	itype int32		//allows for four billion different types of items/blocks
	//nbt string	// unused (for now!)
}

func GiveItem(inv []*ItemStack, ite ItemStack) {
	iltd:=ite.amnt
	for i:=0;i!=52;i++ {
		if inv[i].itype == ite.itype {
			isa:=85-inv[i].amnt
			if iltd>isa {
				inv[i].amnt=85
			} else {
				inv[i].amnt+=iltd
				return
			}
			iltd-=isa
		}
	}
	for i:=0;i!=52;i++ {
		if inv[i].itype == 0 {
			if iltd>85 {
				inv[i].amnt=85
				iltd-=85
			} else {
				inv[i].amnt+=iltd
				iltd = 0
				return
			}
		}
	}
}

func CleanInv(inv []*ItemStack) {
	for i:=0;i!=52;i++ {
		if inv[i].amnt==0 {
			inv[i].itype=0
		}
	}
}