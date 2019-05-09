package cube

type Type struct {
	Name string
}

func (this *Type) String() string {
	return this.Name
}

var (
	TypeAuto  = &Type{"auto"}
	TypeInt64 = &Type{"int64"}
)
