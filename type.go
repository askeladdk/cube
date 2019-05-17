package cube

type Type struct {
	Name string
}

func (this *Type) String() string {
	return this.Name
}

var (
	TypeUntyped64 = &Type{"u64"}
)
