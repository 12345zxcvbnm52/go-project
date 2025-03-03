package errors

//自定义的Error Code总共7位,第一位若为1则表示rpc服务错误码,为2则表示http服务错误码,
//第2-3位共同表示服务号,第4-7位共同表示一个错误码的唯一标识,
//如1110001中第一个1表示为rpc服务错误码,之后11表示第一个服务(这里是用户服务),
//再之后0001唯一标识一个错误码
//无论是http服务还是rpc服务,服务标识中的10统一分配给公有的错误,如110001服务器内部错误将在各服务中共享使用,
//注解由:符号分割为三部分,第一部分表示该错误码对应的grpc错误码的名称,第二部分为http错误码,第三部分为对外不敏感的信息,

const (
	//OK:200:OK
	CodeSuccess = 1100000
	//Internal:500:服务器内部错误
	CodeInternalError = 1100001 + iota
	// PermissionDenied:403:缺少访问的权限
	CodePermissionDenied
	// Canceled:499:客户端关闭请求或连接超时
	CodeCanceled
	// Internal:500:Rockmq消息队列操作失败
	CodeBadRockmq
	// Internal:500:商品服务连接失败
	CodeBadGoodsClient
	// Internal:500:库存服务连接失败
	CodeBadInventoryClient
	// Internal:500:用户服务连接失败
	CodeBadUserClient
	// Internal:500:订单服务连接失败
	CodeBadOrderClient
)

const (
	//Unauthenticated:401:Token已失效
	CodeExpired = 2100001 + iota
	//Unauthenticated:401:Authorization header为空
	CodeMissingHeader
	//Unauthenticated:401:authorization header参数
	CodeInvalidAuthHeader
	//InvalidArgument:400:传入参数有误,认证失败
	CodeValidationInvalid
	//Unauthenticated:401:Token有误或者已失效
	CodeTokenInvalid
)

const (
	// NotFound:404:用户不存在
	CodeUserNotFound = 1110001 + iota
	// AlreadyExists:400:用户已存在
	CodeUserAlreadyExist
	// InvalidArgument:400:错误的账户或密码
	CodeBadAuth
)

const (
	//NotFound:404:所选分类不存在
	CodeCategoryNotFound = 1120001 + iota
	//NotFound:404:所选滑窗不存在
	CodeBannerNotFound
	//NotFound:404:对应商品不存在
	CodeGoodsNotFound
	//NotFound:404:对应品牌不存在
	CodeBrandNotFound
	//NotFound:404:该品牌未创建任何商品类型
	CodeCategoryBrandNotFound

	//AlreadyExists:400:欲创建的商品类型已存在
	CodeCategoryAlreadyExist
	//AlreadyExists:400:欲创建的滑窗已存在
	CodeBannerAlreadyExist
	//AlreadyExists:400:欲创建的商品已存在
	CodeGoodsAlreadyExist
	//AlreadyExists:400:欲创建的品牌已存在
	CodeBrandAlreadyExist
	//AlreadyExists:400:该品牌欲添加的商品类型已存在
	CodeCategoryBrandAlreadyExist

	//FailedPrecondition:400:当前类型存在子类型,无法修改或删除
	CodeCategoryRefered
)

const (
	//NotFound:404:对应库存信息不存在
	CodeInventoryNotFound = 1130001 + iota
	//AlreadyExists:400:欲创建的商品库存已存在
	CodeInventoryAlreadyExist
	//ResourceExhausted:400:指定商品缺少库存
	CodeLackInventory
)

const (
	//NotFound:404:对应的订单不存在
	CodeOrderNotFound = 1140001 + iota
	//AlreadyExists:400:欲创建的订单已存在
	CodeOrderAlreadyExist
	//NotFound:404:购物车内不存在商品
	CodeCartNoItems
	//InvalidArgument:400:购物车内没有选中的商品
	CodeCartNoSelected

	//NotFound:404:未找到对应的订单内商品信息
	CodeOrderGoodsNotFound
	//AlreadyExists:400:欲创建的订单内商品信息已存在
	CodeOrderGoodsAlreadyExist

	//Aborted:400:订单创建失败
	CodeOrderFailedCreate
)
