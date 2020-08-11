package items

type ItemStack struct {
	Amnt int8		//can be at most 85 (legally) and at the least 1
	Itype int32		//allows for four billion different types of items/blocks
	//nbt string	// unused (for now!)
}

func GiveItem(inv *[52]ItemStack, ite ItemStack) {
	iltd:=ite.Amnt
	for i:=0;i!=52;i++ {
		if inv[i].Itype == ite.Itype {
			isa:=85-inv[i].Amnt
			if iltd>isa {
				inv[i].Amnt=85
			} else {
				inv[i].Amnt+=iltd
				return
			}
			iltd-=isa
		}
	}
	for i:=0;i!=52;i++ {
		if inv[i].Itype == 0 {
			if iltd>85 {
				inv[i].Amnt=85
				iltd-=85
			} else {
				inv[i].Amnt+=iltd
				iltd = 0
				return
			}
		}
	}
}

func CleanInv(inv [52]ItemStack) ([52]ItemStack) {
	for i:=0;i!=52;i++ {
		if inv[i].Amnt==0 {
			inv[i].Itype=0
		}
	}
	return inv
}