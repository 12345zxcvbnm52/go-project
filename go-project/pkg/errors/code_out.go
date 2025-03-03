package errors

import (
	"google.golang.org/grpc/codes"
)

var CodeSuccess Coder

var CodeInternalError Coder

var CodePermissionDenied Coder

var CodeCanceled Coder

var CodeBadRockmq Coder

var CodeBadGoodsClient Coder

var CodeBadInventoryClient Coder

var CodeBadUserClient Coder

var CodeBadOrderClient Coder

func init() {

	CodeSuccess = mustNewCoder(1100000, 200, codes.OK, "OK")

	CodeInternalError = mustNewCoder(1100001, 500, codes.Internal, "服务器内部错误")

	CodePermissionDenied = mustNewCoder(1100002, 403, codes.PermissionDenied, "缺少访问的权限")

	CodeCanceled = mustNewCoder(1100003, 499, codes.Canceled, "客户端关闭请求或连接超时")

	CodeBadRockmq = mustNewCoder(1100004, 500, codes.Internal, "Rockmq消息队列操作失败")

	CodeBadGoodsClient = mustNewCoder(1100005, 500, codes.Internal, "商品服务连接失败")

	CodeBadInventoryClient = mustNewCoder(1100006, 500, codes.Internal, "库存服务连接失败")

	CodeBadUserClient = mustNewCoder(1100007, 500, codes.Internal, "用户服务连接失败")

	CodeBadOrderClient = mustNewCoder(1100008, 500, codes.Internal, "订单服务连接失败")

}

var CodeExpired Coder

var CodeMissingHeader Coder

var CodeInvalidAuthHeader Coder

var CodeValidationInvalid Coder

var CodeTokenInvalid Coder

func init() {

	CodeExpired = mustNewCoder(2100001, 401, codes.Unauthenticated, "Token已失效")

	CodeMissingHeader = mustNewCoder(2100002, 401, codes.Unauthenticated, "Authorization header为空")

	CodeInvalidAuthHeader = mustNewCoder(2100003, 401, codes.Unauthenticated, "authorization header参数")

	CodeValidationInvalid = mustNewCoder(2100004, 400, codes.InvalidArgument, "传入参数有误,认证失败")

	CodeTokenInvalid = mustNewCoder(2100005, 401, codes.Unauthenticated, "Token有误或者已失效")

}

var CodeUserNotFound Coder

var CodeUserAlreadyExist Coder

var CodeBadAuth Coder

func init() {

	CodeUserNotFound = mustNewCoder(1110001, 404, codes.NotFound, "用户不存在")

	CodeUserAlreadyExist = mustNewCoder(1110002, 400, codes.AlreadyExists, "用户已存在")

	CodeBadAuth = mustNewCoder(1110003, 400, codes.InvalidArgument, "错误的账户或密码")

}

var CodeCategoryNotFound Coder

var CodeBannerNotFound Coder

var CodeGoodsNotFound Coder

var CodeBrandNotFound Coder

var CodeCategoryBrandNotFound Coder

var CodeCategoryAlreadyExist Coder

var CodeBannerAlreadyExist Coder

var CodeGoodsAlreadyExist Coder

var CodeBrandAlreadyExist Coder

var CodeCategoryBrandAlreadyExist Coder

var CodeCategoryRefered Coder

func init() {

	CodeCategoryNotFound = mustNewCoder(1120001, 404, codes.NotFound, "所选分类不存在")

	CodeBannerNotFound = mustNewCoder(1120002, 404, codes.NotFound, "所选滑窗不存在")

	CodeGoodsNotFound = mustNewCoder(1120003, 404, codes.NotFound, "对应商品不存在")

	CodeBrandNotFound = mustNewCoder(1120004, 404, codes.NotFound, "对应品牌不存在")

	CodeCategoryBrandNotFound = mustNewCoder(1120005, 404, codes.NotFound, "该品牌未创建任何商品类型")

	CodeCategoryAlreadyExist = mustNewCoder(1120006, 400, codes.AlreadyExists, "欲创建的商品类型已存在")

	CodeBannerAlreadyExist = mustNewCoder(1120007, 400, codes.AlreadyExists, "欲创建的滑窗已存在")

	CodeGoodsAlreadyExist = mustNewCoder(1120008, 400, codes.AlreadyExists, "欲创建的商品已存在")

	CodeBrandAlreadyExist = mustNewCoder(1120009, 400, codes.AlreadyExists, "欲创建的品牌已存在")

	CodeCategoryBrandAlreadyExist = mustNewCoder(1120010, 400, codes.AlreadyExists, "该品牌欲添加的商品类型已存在")

	CodeCategoryRefered = mustNewCoder(1120011, 400, codes.FailedPrecondition, "当前类型存在子类型,无法修改或删除")

}

var CodeInventoryNotFound Coder

var CodeInventoryAlreadyExist Coder

var CodeLackInventory Coder

func init() {

	CodeInventoryNotFound = mustNewCoder(1130001, 404, codes.NotFound, "对应库存信息不存在")

	CodeInventoryAlreadyExist = mustNewCoder(1130002, 400, codes.AlreadyExists, "欲创建的商品库存已存在")

	CodeLackInventory = mustNewCoder(1130003, 400, codes.ResourceExhausted, "指定商品缺少库存")

}

var CodeOrderNotFound Coder

var CodeOrderAlreadyExist Coder

var CodeCartNoItems Coder

var CodeCartNoSelected Coder

var CodeOrderGoodsNotFound Coder

var CodeOrderGoodsAlreadyExist Coder

var CodeOrderFailedCreate Coder

func init() {

	CodeOrderNotFound = mustNewCoder(1140001, 404, codes.NotFound, "对应的订单不存在")

	CodeOrderAlreadyExist = mustNewCoder(1140002, 400, codes.AlreadyExists, "欲创建的订单已存在")

	CodeCartNoItems = mustNewCoder(1140003, 404, codes.NotFound, "购物车内不存在商品")

	CodeCartNoSelected = mustNewCoder(1140004, 400, codes.InvalidArgument, "购物车内没有选中的商品")

	CodeOrderGoodsNotFound = mustNewCoder(1140005, 404, codes.NotFound, "未找到对应的订单内商品信息")

	CodeOrderGoodsAlreadyExist = mustNewCoder(1140006, 400, codes.AlreadyExists, "欲创建的订单内商品信息已存在")

	CodeOrderFailedCreate = mustNewCoder(1140007, 400, codes.Aborted, "订单创建失败")

}
