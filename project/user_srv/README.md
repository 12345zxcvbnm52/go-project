# Message 
基于go实现的电商微服务

# 关于Model.Logic里的冗余错误检查
主要是提醒自己日后扩展时可以考虑给增删改添加退避策略等属性

# 关于Model.Logic编写要求
所有的对一行数据的删改查应当只能基于unique键和主键,如FindOneById,FindOneByMobile,
涉及多行数据的查则推荐返回数据,RowEffected,error,同时只允许unique键和主键查询,
其余的复杂逻辑则考虑在handler中编写(暂时未考虑多行数据的删改查)

# Model.Logic title
这样写就可以把所有数据全读到model结构体中,然后根据rpc的要求返回需要的数据,
而且未来可以给每一个logic方法添加一个可有的DB参数用于omit等行为

# 关于错误处理
建议未来的错误处理主要在srv层(统一错误),而web根据得到的error进行翻译 x
建议用户端的错误统一转化为status.Error 按codes分好,web端只需直接打印即可 v

# 关于传递到微服务的数据
建议所有有关参数都在传递到微服务之前处理好,

# 关于logic函数的gorm查询
用Find与切片查询多条数据时,要保证查询条件非空,不然就是全表查询,查询单条数据则考虑(?),