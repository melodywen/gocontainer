package ioc

import (
	"reflect"
)

type Container struct {
	instances map[string]interface{} // 绑定的实例 ， 如果他是单例模式则全部存储到这里面
	bindings  map[string]interface{} // 绑定的策略及其配置
}

func newContainer() *Container {
	instances := make(map[string]interface{})
	obj := &Container{instances: instances}
	return obj
}

//  获取对应的接口名称
// abstract 目前完成了 对象 、 接口 和字符串
func (container *Container) checkAbstract(abstract interface{}) string {
	classInfo := reflect.TypeOf(abstract)
	kind := classInfo.Kind().String()

	var pkgPath string
	var name string

	switch kind {
	case "string":
		return abstract.(string)
	case "struct":

		pkgPath = classInfo.PkgPath()
		name = classInfo.Name()
	case "ptr":
		classInfo = classInfo.Elem()
		pkgPath = classInfo.PkgPath()
		name = classInfo.Name()
	default:
		panic("ioc struct index value must implication string:but current is " + kind)
	}
	return pkgPath + "/" + name
}

// Bind 绑定一个实例
func (container *Container) Bind(abstract interface{}, concrete interface{}, shared bool) {
	// get abstract value to string ,set to component index
	index := container.checkAbstract(abstract)

	container.DropStaleInstances(index)

	container.bindings[index] = map[string]interface{}{
		"shared":   shared,
		"concrete": concrete,
	}
}

// DropStaleInstances 删除一个老的实例
func (container *Container) DropStaleInstances(abstract interface{}) {
	index := container.checkAbstract(abstract)
	if _, ok := container.instances[index]; ok {
		delete(container.instances, index)
	}
}

// Make 对外暴露make 方法
func (container *Container) Make(abstract interface{}, parameters []interface{}) interface{} {
	return container.resolve(abstract, parameters, true)
}

func (container *Container) resolved() {

}

// 生成 一个 对应的 内容
func (container *Container) resolve(abstract interface{}, parameters []interface{}, raiseEvents bool) interface{} {
	index := container.checkAbstract(abstract)
	// 如果没有上下文出现了绑定了实例则直接返回
	if _, ok := container.bindings[index]; ok {
		return container.bindings["concrete"]
	}

	concrete := container.getConcrete(index)

	return concrete
}

func (container *Container) rebound() {

}

// 获取实际的实现方式
func (container *Container) getConcrete(abstract interface{}) (concrete interface{}) {

	index := container.checkAbstract(abstract)
	// TODO : 如果存在上下文的绑定则返回 上下文的内容

	// 如果 设置了绑定的内容则返回绑定的内容
	if _, ok := container.bindings[index]; ok {
		return container.bindings["concrete"]
	}
	concrete = abstract
	return concrete
}

func (container *Container) IsShared(concrete interface{}) (bool) {

	return true
}

// 判断是否可以侯建
func (container *Container) isBuildable(abstract interface{}, concrete interface{}) bool {
	//classInfo := reflect.TypeOf(abstract).Elem()
	////if _, ok := abstract.(classInfo); ok {
	////}
	return true
}

// Build 动态构建一个实例出来：
//	1. 可以是某个类的构造方法
//	2. 也可以是回调函数
// 	3. 或者是一个接口，一个具体的标量值 等等
func (container *Container) Build(concrete interface{}, parameters []interface{}) (object interface{}) {

	// 获取实现类的类型
	concreteType := reflect.TypeOf(concrete)

	switch concreteType.Kind() {
	case reflect.Func:
		// 获取实现类的值
		concreteValue := reflect.ValueOf(concrete)
		// 函数的形参绑定
		var params []reflect.Value
		for _, parameter := range parameters {
			params = append(params, reflect.ValueOf(parameter))
		}
		// 调用函数
		resultList := concreteValue.Call(params)
		// 然后进行克隆反射
		numOut := concreteValue.Type().NumOut()
		response := []interface{}{}
		for m := 0; m < numOut; m++ {
			returnType := concreteValue.Type().Out(m)
			switch returnType.Kind() {
			case reflect.Ptr: //如果是指针类型
				returnNew := reflect.New(returnType.Elem()).Elem() //创建对象 //获取源实际类型(否则为指针类型)
				returnValue := resultList[m]                       //源数据值
				returnValue = returnValue.Elem()                   //源数据实际值（否则为指针）
				returnNew.Set(returnValue)                         //设置数据
				returnNew = returnNew.Addr()                       //创建对象的地址（否则返回值）
				response = append(response, returnNew.Interface()) //返回地址
			default:
				returnNew := reflect.New(returnType).Elem()        //创建对象
				returnValue := resultList[m]                       //源数据值
				returnNew.Set(returnValue)                         //设置数据
				response = append(response, returnNew.Interface()) //返回
			}
		}
		if len(response) == 1 {
			return response[0]
		}
		return response
	default:
		return concrete
	}
}


