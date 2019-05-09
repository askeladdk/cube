package cube

type Type struct {
	Name string
}

func (this *Type) String() string {
	return this.Name
}

var (
	TypeAuto  = &Type{"auto"}
	TypeInt32 = &Type{"int32"}
	TypeInt64 = &Type{"int64"}
)
