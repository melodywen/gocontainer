package ioc

import (
	"fmt"
	"reflect"
	"testing"
)

type Person struct {
	username string
	age      int
}

func NewPerson(username string, age int) *Person {
	return &Person{username: username, age: age}
}

type SuperPerson struct {
	Person
	username string
}

func (person *Person) say() string {
	return "i am " + person.username
}

func (person *SuperPerson) say() string {
	return "i am " + person.username
}

func (person *SuperPerson) work() string {
	return "i am working"
}

type HuMen interface {
	say() string
}

type Chinese interface {
	HuMen
	work() string
}

func add(a int, b int) int {
	return a + b
}
func add2(a int, b int) (int, int) {
	return a + b, a - b
}
func addToPerson(a int, b int) Person {
	return *NewPerson("melody", a+b)
}

func addToPerson2(a int, b int) *Person {
	return NewPerson("melody", a+b)
}

func addToPerson3(a int, b int) (*Person, int) {
	return NewPerson("melody", a+b), a - b
}

func Test_container_Bind(t *testing.T) {
	type fields struct {
		instances map[string]interface{}
	}
	type args struct {
		abstract interface{}
		concrete interface{}
		shared   bool
	}

	p1 := &Person{username: "melody", age: 18}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name: "测试bind方法：是一个string 和 标量 类型",
			fields: fields{instances: map[string]interface{}{
				"username": map[string]interface{}{
					"concrete": "melody", "shared": true,
				},
			}},
			args: args{
				abstract: "username",
				concrete: "melody",
				shared:   true,
			},
		},
		{
			name: "测试bind方法：测试是接口和接口类型",
			fields: fields{instances: map[string]interface{}{
				(func() string {
					return reflect.TypeOf(p1).Elem().PkgPath() + "/" + reflect.TypeOf(p1).Elem().Name()
				})(): map[string]interface{}{
					"concrete": &Person{username: "melody", age: 18}, "shared": true,
				},
			}},
			args: args{
				abstract: p1,
				concrete: &Person{username: "melody", age: 18},
				shared:   true,
			},
		},
		{
			name: "测试bind方法：测试是接口和接口类型",
			fields: fields{instances: map[string]interface{}{
				(func() string {
					return reflect.TypeOf(p1).Elem().PkgPath() + "/" + reflect.TypeOf(p1).Elem().Name()
				})(): map[string]interface{}{
					"concrete": &Person{username: "melody", age: 18}, "shared": true,
				},
			}},
			args: args{
				abstract: *p1,
				concrete: &Person{username: "melody", age: 18},
				shared:   true,
			},
		},
	}
	for _, tt := range tests {
		fmt.Println("进入到循环")
		t.Run(tt.name, func(t *testing.T) {
			mock := newContainer()
			mock.Bind(tt.args.abstract, tt.args.concrete, tt.args.shared)
			if !reflect.DeepEqual(mock.instances, tt.fields.instances) {
				t.Errorf("newContainer() = %v, want %v", mock, tt.fields.instances)
			} else {
				fmt.Println("测试通过", mock)
			}
		})
	}
}

func Test_newContainer(t *testing.T) {
	tests := []struct {
		name string
		want *Container
	}{
		// TODO: Add test cases.
		{
			name: "测试",
			want: newContainer(),
		},
	}

	for _, tt := range tests {
		fmt.Println("进入到循环")
		t.Run(tt.name, func(t *testing.T) {
			got := newContainer()
			fmt.Println("伪造到的数据：", got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newContainer() = %v, want %v", got, tt.want)
			} else {
				fmt.Println("测试通过")
			}
		})
	}
}

func Test_container_DropStaleInstances(t *testing.T) {
	type fields struct {
		instances map[string]interface{}
	}
	type args struct {
		abstract interface{}
		concrete interface{}
		shared   bool
	}
	p1 := &Person{username: "melody", age: 18}

	var tests = []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "测试一下",
			fields: fields{instances: map[string]interface{}{}},
			args: args{
				abstract: p1,
				concrete: &Person{username: "melody", age: 18},
				shared:   true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newContainer()
			got.Bind(tt.args.abstract, tt.args.concrete, tt.args.shared)
			fmt.Println("伪造到的数据：", got)
			got.DropStaleInstances(tt.args.abstract)
			fmt.Println("伪造到的数据：", got)
			if !reflect.DeepEqual(got.instances, tt.fields.instances) {
				t.Errorf("newContainer() = %v, want %v", got, tt.fields)
			} else {
				fmt.Println("测试通过")
			}
		})
	}
}

func Test_container_isBuildable(t *testing.T) {
	type fields struct {
		instances map[string]interface{}
		bindings  map[string]interface{}
	}
	var humen *HuMen
	fmt.Println(reflect.TypeOf(humen).Elem())

	type args struct {
		abstract interface{}
		concrete interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
		//{
		//	name: "判定接口与接口是否为继承关系",
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			con := &Container{
				instances: tt.fields.instances,
				bindings:  tt.fields.bindings,
			}
			if got := con.isBuildable(tt.args.abstract, tt.args.concrete); got != tt.want {
				t.Errorf("isBuildable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainer_Build(t *testing.T) {
	type fields struct {
		instances map[string]interface{}
		bindings  map[string]interface{}
	}
	type args struct {
		concrete   interface{}
		parameters []interface{}
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantObject interface{}
	}{
		// TODO: Add test cases.
		{
			name:   "如果是一个函数,返回一个值",
			fields: fields{},
			args: args{
				concrete:   add,
				parameters: []interface{}{4, 5},
			},
			wantObject: 9,
		},
		{
			name:   "如果是一个函数,返回二个值",
			fields: fields{},
			args: args{
				concrete:   add2,
				parameters: []interface{}{4, 5},
			},
			wantObject: []interface{}{9, -1},
		},
		{
			name:   "如果是一个函数,返回是一个结构体",
			fields: fields{},
			args: args{
				concrete:   addToPerson,
				parameters: []interface{}{4, 5},
			},
			wantObject: *NewPerson("melody", 9),
		},
		{
			name:   "如果是一个函数,返回是一个结构体指针",
			fields: fields{},
			args: args{
				concrete:   addToPerson2,
				parameters: []interface{}{4, 5},
			},
			wantObject: NewPerson("melody", 9),
		},
		{
			name:   "如果是一个函数,返回是一个结构体指针 和int ",
			fields: fields{},
			args: args{
				concrete:   addToPerson3,
				parameters: []interface{}{4, 5},
			},
			wantObject: []interface{}{NewPerson("melody", 9), -1},
		},
		{
			name:   "如果是一个标量  ",
			fields: fields{},
			args: args{
				concrete:   2,
				parameters: []interface{}{},
			},
			wantObject: 2,
		},
		{
			name:   "如果是一个切片  ",
			fields: fields{},
			args: args{
				concrete:   []Person{{username: "melody"}},
				parameters: []interface{}{},
			},
			wantObject: []Person{{username: "melody"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := newContainer()
			//gotObject := container.Build(tt.args.concrete, tt.args.parameters)
			//fmt.Println(gotObject, index, tt.wantObject, gotObject == 9)
			if gotObject := container.Build(tt.args.concrete, tt.args.parameters); !reflect.DeepEqual(gotObject, tt.wantObject) {
				t.Errorf("Build() = %v, want %v", gotObject, tt.wantObject)
			} else {
				fmt.Println("验证成功", gotObject)
			}
		})
	}
}
