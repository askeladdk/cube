package cube

type Type struct {
	Name string
}

func (this *Type) String() string {
	return this.Name
}

var TypeInt32 = Type{"int32"}
var TypeInt64 = Type{"int64"}
